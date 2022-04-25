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
	usdt := fmt.Sprintf("%s:   %.2f", bold.Render("USDT"), m.usdt)
	btcusd := fmt.Sprintf("%s: %d", bold.Render("BTCUSD"), m.btcusd)

	futures := ""
	for contract, size := range m.futures {
		futures += fmt.Sprintf("%s: %d\n", bold.Render(contract), size)
	}

	return margin.Render(fmt.Sprintf("%s\n%s\n%s", usdt, btcusd, futures))
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
