package tui

import (
	"fmt"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/qrave1/amodns/internal/dns"
)

// Определяем возможные действия
type Choice string

const (
	StageChoice Choice = "stage"
	ProdChoice  Choice = "prod"
)

// Модель программы
type Model struct {
	dns            *dns.DNSchanger
	activeConnName string
	activeConnDNS  string

	choices  []Choice // Варианты выбора DNS
	cursor   int      // Позиция курсора
	selected Choice   // Выбранный DNS

	styler *Styler
}

func NewModel(dns *dns.DNSchanger, choices []Choice, styler *Styler) (*Model, error) {
	name, err := dns.GetActiveConnectionName()
	if err != nil {
		return nil, fmt.Errorf("error get active connection name: %w", err)
	}

	connDNS, err := dns.GetCurrentDNS()
	if err != nil {
		return nil, fmt.Errorf("error get current DNS: %w", err)
	}
	return &Model{
		dns:            dns,
		choices:        choices,
		activeConnName: name,
		activeConnDNS:  connDNS,
		styler:         styler,
	}, nil
}

// Инициализация модели
func (m Model) Init() tea.Cmd {
	return nil
}

// Обновление модели на основе входящих сообщений
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case "enter":
			dnsAddr := m.dns.MapEnvToAddr(string(m.choices[m.cursor]))
			conn, err := m.dns.GetActiveConnectionName()
			if err != nil {
				log.Printf("error getting active connection name: %v", err)
				return m, tea.Quit
			}
			err = m.dns.SetDNS(conn, dnsAddr)
			if err != nil {
				log.Printf("error change DNS: %v", err)
			}
			return m, tea.Quit
		}
	}

	return m, nil
}

// Отображение интерфейса
func (m Model) View() string {
	var s strings.Builder

	// Заголовок
	s.WriteString(m.styler.titleStyle.Render("amoDNS") + "\n\n")

	// Информация о соединении
	s.WriteString(m.styler.infoStyle.Render(fmt.Sprintf("WI-FI Connection Name: %s", m.activeConnName)) + "\n")
	s.WriteString(m.styler.infoStyle.Render(fmt.Sprintf("Current DNS: %s", m.activeConnDNS)) + "\n\n")

	// Список выбора DNS
	s.WriteString(m.styler.instructionStyle.Render("Select DNS:") + "\n\n")

	for i, choice := range m.choices {
		var checkbox string
		lineStyle := m.styler.defaultStyle

		if m.cursor == i {
			checkbox = m.styler.cursorStyle.Render("[ x ]") // Отображение курсора на выбранной опции
			lineStyle = m.styler.highlightStyle
		} else {
			checkbox = "[   ]" // Не выбранная опция
		}

		s.WriteString(fmt.Sprintf("%s %s\n", checkbox, lineStyle.Render(string(choice))))
	}

	// Инструкции
	s.WriteString("\n")
	s.WriteString(m.styler.instructionStyle.Render("Press Enter to select, q to quit.") + "\n")

	if m.selected != "" {
		s.WriteString("\n")
		s.WriteString(m.styler.selectedDNSStyle.Render(fmt.Sprintf("Selected DNS: %s", m.selected)) + "\n")
	}

	return s.String()
}
