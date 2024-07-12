package main

import (
	"bufio"
	"fmt"
	"os"

	// "os/exec"
	// "strings"
	tea "github.com/charmbracelet/bubbletea"
  "github.com/charmbracelet/lipgloss"
)

var version = "0.01"
var revision = "1"

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type model struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
}

func initialModel(texts []string) model {
	return model{
		choices:  texts,
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "esc":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		case "enter":
			_, ok := m.selected[m.cursor]
			print(ok)
      for el, _ := range m.selected {
        w := bufio.NewWriter(os.Stdout)
        defer w.Flush()
        w.WriteString(m.choices[el] + "\n")
      }
      return m, tea.Quit
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	s := ""

	// Iterate over our choices
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
  texts := []string {}
  for scanner.Scan() {
    texts = append(texts, scanner.Text())
  }

	p := tea.NewProgram(initialModel(texts), tea.WithOutput(os.Stderr), tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

}
