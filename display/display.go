package display

import (
	"github.com/RFloTeo/power-spy/resources"
	tea "github.com/charmbracelet/bubbletea"
	"time"
)

type Model struct {
	stats    map[string]resources.DockerStats
	duration time.Duration
}

type TickMsg time.Time

func (m Model) Init() tea.Cmd {
	return tea.Every(m.duration, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case TickMsg:
		keys := m.getKeys()
		m.stats = resources.GetStats(keys)
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "r":
			m.stats = resources.Refresh()
		}
	}
	return m, nil
}

func (m Model) View() string {
	return ""
}

func (m Model) getKeys() []string {
	keys := make([]string, len(m.stats))
	i := 0
	for key := range m.stats {
		keys[i] = key
		i++
	}
	return keys
}
