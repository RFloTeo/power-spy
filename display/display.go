package display

import (
	"bytes"
	"fmt"
	"github.com/RFloTeo/power-spy/resources"
	tea "github.com/charmbracelet/bubbletea"
	"text/tabwriter"
	"time"
)

type Model struct {
	Stats    map[string]resources.DockerStats
	Duration time.Duration
}

type TickMsg time.Time

func (m Model) Init() tea.Cmd {
	return tea.Every(m.Duration, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case TickMsg:
		keys := m.getKeys()
		m.Stats = resources.GetStats(keys)
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "r":
			newStats, err := resources.Refresh()
			if err == nil {
				m.Stats = newStats
			}
		}
	}
	return m, tea.Every(m.Duration, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m Model) View() string {
	//TODO: Adapt to power usage and meaning of metrics when you find out
	var b bytes.Buffer
	tab := tabwriter.NewWriter(&b, 5, 4, 1, '\t', 0)
	s := "q-Quit r-Refresh\n"
	fmt.Fprintln(tab, "ID\tCPU\tNetIn\tNetOut\tPower")
	for key, stats := range m.Stats {
		fmt.Fprintf(tab, "%s\t%d\t%d\t%d\n", key, stats.CPU, stats.NetworkIn, stats.NetworkOut)
	}
	tab.Flush()
	s += b.String()
	return s
}

func (m Model) getKeys() []string {
	keys := make([]string, len(m.Stats))
	i := 0
	for key := range m.Stats {
		keys[i] = key
		i++
	}
	return keys
}
