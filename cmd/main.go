package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	dns "github.com/qrave1/amodns/internal/dns"
	"github.com/qrave1/amodns/internal/interfaces/tui"
)

func main() {
	dnsChanger := dns.NewDNSchanger()
	styler := tui.NewDefaultStyler()
	m, err := tui.NewModel(
		dnsChanger,
		[]tui.Choice{tui.StageChoice, tui.ProdChoice},
		styler,
	)
	if err != nil {
		fmt.Printf("error creating model: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("error run amodns: %v\n", err)
		os.Exit(1)
	}
}
