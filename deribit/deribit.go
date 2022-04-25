package deribit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"
)

type Deribit struct {
	ApiId        string
	ApiSecret    string
	Test         bool
	_accessToken string
	expiresIn    time.Time
}

func (d *Deribit) accessToken() (string, error) {
	if d._accessToken != "" && d.expiresIn.After(time.Now()) {
		return d._accessToken, nil
	}

	log.WithField("venue", "deribit").Debug("Fetching access token")

	v := url.Values{}
	v.Set("client_id", d.ApiId)
	v.Set("client_secret", d.ApiSecret)
	v.Set("grant_type", "client_credentials")

	u := url.URL{
		Scheme:   "https",
		Host:     d.hostname(),
		Path:     "/api/v2/public/auth",
		RawQuery: v.Encode()}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var response authResponse
	json.Unmarshal(body, &response)
	d._accessToken = response.Result.AccessToken
	expirySecs := time.Second * time.Duration(response.Result.ExpiresIn-10)
	d.expiresIn = time.Now().Add(expirySecs)

	return d._accessToken, nil
}

func (d *Deribit) hostname() string {
	if d.Test {
		return "test.deribit.com"
	} else {
		return "www.deribit.com"
	}
}

func (d *Deribit) get(path string, params url.Values, result interface{}) error {
	u := url.URL{
		Scheme:   "https",
		Host:     d.hostname(),
		Path:     path,
		RawQuery: params.Encode()}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	accessToken, err := d.accessToken()
	if err != nil {
		log.Error(err.Error())
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Warn(err.Error())
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	json.Unmarshal(body, result)

	return nil
}

func (d *Deribit) GetPositions() map[string]int {
	result := map[string]int{}
	var response positionsResponse

	err := d.get("/api/v2/private/get_positions",
		url.Values{"currency": {"BTC"}, "kind": {"future"}}, &response)

	if err != nil {
		return result
	}

	for _, position := range response.Result {
		result[position.InstrumentName] = int(position.Size)
	}

	return result
}
