package internal

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type timerModel struct {
	duration time.Duration
	elapsed  time.Duration
	done     bool
}

func initialTimerModel(duration time.Duration) timerModel {
	return timerModel{
		duration: duration,
		elapsed:  0,
		done:     false,
	}
}

func (m timerModel) Init() tea.Cmd {
	return tickCmd()
}

func (m timerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		m.elapsed += time.Second
		if m.elapsed >= m.duration {
			m.done = true
			return m, tea.Quit
		}
		return m, tickCmd()
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m timerModel) View() string {
	remaining := m.duration - m.elapsed
	if remaining < 0 {
		remaining = 0
	}
	return fmt.Sprintf("Time remaining: %s\n\nPress Ctrl+C or q to quit.\n", remaining.String())
}

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(time.Second)
		return tickMsg(time.Now())
	}
}

func RunTimer(duration time.Duration) error {
	p := tea.NewProgram(initialTimerModel(duration))
	if _, err := p.Run(); err != nil {
		return err
	}
	fmt.Println("Timer completed.")
	return nil
}
