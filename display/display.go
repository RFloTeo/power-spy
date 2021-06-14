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
	MuW        int
}

type TickMsg time.Time

func (m Model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, tea.Tick(m.Duration, func(t time.Time) tea.Msg {
		return TickMsg(t)
	}))
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case TickMsg:
		keys := m.getKeys()
		newStats := resources.GetStats(keys)
		for newKey := range newStats {
			m.Stats[newKey] = newStats[newKey]
		}
		if resources.PowerOn {
			m.MuW = resources.GetMicroWatts()
		}
		go resources.RecordStats(m.Containers, m.Stats, m.MuW)
		return m, tea.Tick(m.Duration, func(t time.Time) tea.Msg {
			return TickMsg(t)
		})

	case tea.KeyMsg:
		if m.Text.Focused() && msg.String() != " " {
			var cmd tea.Cmd
			m.Text, cmd = m.Text.Update(msg)
			return m, cmd
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
			err := resources.ToggleRecording(m.Filter, len(m.Containers))
			m.ToggleFail = err != nil
		case " ": // Focus/Blur textbox
			if !resources.IsRecording {
				if m.Text.Focused() {
					m.Text.Blur()
					m.Filter = m.Text.Value()
				} else {
					_ = m.Text.Reset()
					m.Filter = ""
					m.Text.Focus()
				}

			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	var b bytes.Buffer
	tab := tabwriter.NewWriter(&b, 5, 4, 3, ' ', 0)

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

		if strings.Contains(c.Id, m.Filter) || strings.Contains(c.Image, m.Filter) {
			id := c.Id[:16]
			mem := float64(stats.Memory) / 1048576.0
			netIn := getNetUnits(stats.NetworkIn)
			netOut := getNetUnits(stats.NetworkOut)
			fmt.Fprintf(tab, "%s\t%s\t%.3f MiB\t%.2f%%\t%.2f%%\t%s\t%s\n", id, c.Image, mem, stats.MemoryPercent, stats.CPU, netIn, netOut)
		}
	}
	tab.Flush()
	s += b.String()

	// Power reading
	if resources.PowerOn {
		s += fmt.Sprintf("Power: %.2f W\n", float64(m.MuW)/1000000)
	}

	// text box
	s += "\n" + m.Text.View() + "\n\n"

	return s
}

func getNetUnits(i int) string {
	num := float64(i)
	units := []string{"B", "kB", "MB", "GB"}
	times := 0
	for num >= 1000.0 {
		num /= 1000.0
		times++
	}
	return fmt.Sprintf("%.2f %s", num, units[times])
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
