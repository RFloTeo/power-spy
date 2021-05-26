package display

import (
	"bytes"
	"fmt"
	"github.com/RFloTeo/power-spy/resources"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"strings"
	"text/tabwriter"
	"time"
)

type Model struct {
	Containers []resources.Container
	Stats      map[string]resources.DockerStats
	Duration   time.Duration
	StopFail   bool
	ToggleFail bool
	Text       textinput.Model
	Filter     string
}

type TickMsg time.Time

func (m Model) Init() tea.Cmd {
	return tea.Every(m.Duration, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}
	switch msg := msg.(type) {
	case TickMsg:
		keys := m.getKeys()
		m.Stats = resources.GetStats(keys)

	case tea.KeyMsg:
		if m.Text.Focused() && msg.String() != " " {
			m.Text, cmd := m.Text.Update(msg)
			cmds = append(cmds, cmd)
			break
		}
		switch msg.String() {
		case "q": // Quit
			err := resources.StopRecording()
			if err == nil {
				return m, tea.Quit
			}
			log.Println("Failed to stop recording on quit")
			m.StopFail = true
		case "r": // Refresh
			if !resources.IsRecording {
				newContainers, newStats, err := resources.Refresh()
				if err == nil {
					m.Containers = newContainers
					m.Stats = newStats
				}
			}
		case "o": // Toggle recording
			err := resources.ToggleRecording()
			m.ToggleFail = err != nil
		case " ": // Focus/Blur textbox
			if !resources.IsRecording {
				if m.Text.Focused() {
					m.Text.Blur()
					m.Filter = m.Text.Value()
				} else {
					_ = m.Text.Reset()
					m.Text.Focus()
				}

			}
		}
	}
	cmds = append(cmds, tea.Every(m.Duration, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var b bytes.Buffer
	tab := tabwriter.NewWriter(&b, 5, 4, 1, '\t', 0)

	// General Info
	s := "q-Quit"
	if !resources.IsRecording {
		s += "; r-Refresh; space-Toggle Search Box"
	}
	s += "; o-Toggle Recording\n"
	if resources.IsRecording {
		s += "Currently recording\n"
	}
	if m.ToggleFail {
		s += "Warning: Failed to toggle recording!\n"
	}
	if m.StopFail {
		s += "Warning: Couldn't stop recording; exit aborted. Press q to retry\n"
	}

	// Stats table
	fmt.Fprintln(tab, "ID\tName\tMem\tMem%\tCPU\tNetIn\tNetOut")
	for _, c := range m.Containers {
		stats := m.Stats[c.Id]

		name := ""
		if len(c.Names) > 0 {
			name = c.Names[0]
		}
		if strings.Contains(c.Id, m.Filter) || strings.Contains(name, m.Filter) {
			fmt.Fprintf(tab, "%s\t%s\t%d\t%f\t%f\t%d\t%d\n", c.Id, name, stats.Memory, stats.MemoryPercent, stats.CPU, stats.NetworkIn, stats.NetworkOut)
		}
	}

	// Finish up
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
