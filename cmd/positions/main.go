package main

import (
	"fmt"
	"os"
	"time"

	"github.com/stevenwilkin/carry/binance"
	"github.com/stevenwilkin/carry/bybit"
	"github.com/stevenwilkin/carry/deribit"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	_ "github.com/joho/godotenv/autoload"
)

var (
	b      *binance.Binance
	by     *bybit.Bybit
	d      *deribit.Deribit
	margin = lipgloss.NewStyle().Margin(1, 2, 0, 2)
	bold   = lipgloss.NewStyle().Bold(true)
)

type usdtMsg float64
type btcusdMsg int
type futuresMsg map[string]int

type model struct {
	usdt    float64
	btcusd  int
	futures map[string]int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	case usdtMsg:
		m.usdt = float64(msg)
	case btcusdMsg:
		m.btcusd = int(msg)
	case futuresMsg:
		m.futures = map[string]int(msg)
	}

	return m, nil
}

func (m model) View() string {
	var output string

	w := len("USDT:")
	if len(m.futures) > 0 {
		w = len("BTC-31MAR23:")
	} else if m.btcusd != 0 {
		w = len("BTCUSD:")
	}
	width := lipgloss.NewStyle().Width(w)

	if m.usdt != 0 {
		output += fmt.Sprintf("%s %.2f\n", width.Render(bold.Render("USDT")+":"), m.usdt)
	}

	if m.btcusd != 0 {
		output += fmt.Sprintf("%s %d\n", width.Render(bold.Render("BTCUSD")+":"), m.btcusd)
	}

	for contract, size := range m.futures {
		entry := fmt.Sprintf("%s %d\n", width.Render(bold.Render(contract)+":"), size)
		output += entry
	}

	return margin.Render(output)
}

func main() {
	p := tea.NewProgram(model{}, tea.WithAltScreen())

	b = &binance.Binance{
		ApiKey:    os.Getenv("BINANCE_API_KEY"),
		ApiSecret: os.Getenv("BINANCE_API_SECRET")}

	by = &bybit.Bybit{
		ApiKey:    os.Getenv("BYBIT_API_KEY"),
		ApiSecret: os.Getenv("BYBIT_API_SECRET")}

	d = &deribit.Deribit{
		ApiId:     os.Getenv("DERIBIT_API_ID"),
		ApiSecret: os.Getenv("DERIBIT_API_SECRET")}

	go func() {
		t := time.NewTicker(1 * time.Second)

		for {
			usdt, err := b.GetBalance()
			if err != nil {
				panic(err)
			}

			p.Send(usdtMsg(usdt))
			<-t.C
		}
	}()

	go func() {
		t := time.NewTicker(1 * time.Second)

		for {
			btcusd := by.GetSize()

			p.Send(btcusdMsg(btcusd))
			<-t.C
		}
	}()

	go func() {
		t := time.NewTicker(1 * time.Second)

		for {
			futures := d.GetPositions()

			p.Send(futuresMsg(futures))
			<-t.C
		}
	}()

	if err := p.Start(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
