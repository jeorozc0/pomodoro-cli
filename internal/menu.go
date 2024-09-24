package internal

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	workTime, restTime []string
	cursor             int // Cursor position in the current list
	selectedWorkTime   int // Selected index for workTime
	selectedRestTime   int // Selected index for restTime
	state              int // 0 for workTime selection, 1 for restTime selection
}

func initialModel() model {
	return model{
		workTime:         []string{"Short", "Medium", "Long"},
		restTime:         []string{"Short", "Medium", "Long"},
		selectedWorkTime: -1, // No selection yet
		selectedRestTime: -1, // No selection yet
		state:            0,  // Start with workTime selection
	}
}

func (m model) Init() tea.Cmd {
	// No initialization
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			listLength := m.currentListLength()
			if m.cursor < listLength-1 {
				m.cursor++
			}

		case "enter", " ":
			switch m.state {
			case 0:
				m.selectedWorkTime = m.cursor
				m.state = 1  // Move to restTime selection
				m.cursor = 0 // Reset cursor
			case 1:
				m.selectedRestTime = m.cursor
				return m, tea.Quit // Selections made, exit program
			}
		}
	}
	return m, nil
}

// Helper function to get the length of the current list
func (m model) currentListLength() int {
	if m.state == 0 {
		return len(m.workTime)
	}
	return len(m.restTime)
}

func (m model) View() string {
	var s string
	switch m.state {
	case 0:
		s += "Select Work Time:\n\n"
		s += m.renderList(m.workTime)
		s += "\nPress Enter to select work time.\n"
	case 1:
		s += "Select Rest Time:\n\n"
		s += m.renderList(m.restTime)
		s += "\nPress Enter to select rest time.\n"
	}
	return s
}

// Helper function to render a list with cursor
func (m model) renderList(items []string) string {
	var s string
	for i, item := range items {
		cursor := " " // No cursor
		if m.cursor == i {
			cursor = ">" // Cursor
		}
		s += fmt.Sprintf("%s %s\n", cursor, item)
	}
	return s
}

func RunMenu() (string, string, error) {
	p := tea.NewProgram(initialModel())

	// Run the program
	finalModel, err := p.Run()
	if err != nil {
		return "", "", err
	}

	m, ok := finalModel.(model)
	if !ok {
		return "", "", fmt.Errorf("could not assert model")
	}

	workTime := m.workTime[m.selectedWorkTime]
	restTime := m.restTime[m.selectedRestTime]

	return workTime, restTime, nil
}
