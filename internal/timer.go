package internal

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type pomodoroModel struct {
	totalCycles        int           // Total number of Pomodoro cycles to complete
	currentCycle       int           // Current cycle number
	workDuration       time.Duration // Duration of the work period
	shortBreakDuration time.Duration // Duration of the short break
	longBreakDuration  time.Duration // Duration of the long break
	state              timerState    // Current state of the timer
	elapsed            time.Duration // Time elapsed in the current period
	paused             bool          // Whether the timer is paused
}

type timerState int

const (
	stateWork timerState = iota
	stateShortBreak
	stateLongBreak
	stateFinished
)

func initialPomodoroModel(totalCycles int, workDuration, shortBreakDuration, longBreakDuration time.Duration) pomodoroModel {
	return pomodoroModel{
		totalCycles:        totalCycles,
		currentCycle:       1,
		workDuration:       workDuration,
		shortBreakDuration: shortBreakDuration,
		longBreakDuration:  longBreakDuration,
		state:              stateWork,
		elapsed:            0,
		paused:             false,
	}
}

func (p pomodoroModel) Init() tea.Cmd {
	return tickCmd()
}

func (p pomodoroModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		if p.paused != true {
			p.elapsed += time.Second
		}
		return p, tickCmd()
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return p, tea.Quit
		case "p":
			p.paused = !p.paused
			return p, nil
		}
	}
	return p, nil
}

func (p pomodoroModel) View() string {
	return fmt.Sprintf("Time elapsed: %s\n\nPress Ctrl+C or q to quit.\n", p.elapsed.String())
}

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(time.Second)
		return tickMsg(time.Now())
	}
}

func (p pomodoroModel) currentDuration() time.Duration {
	switch p.state {
	case stateWork:
		return p.workDuration
	case stateShortBreak:
		return p.shortBreakDuration
	case stateLongBreak:
		return p.longBreakDuration
	default:
		return 0
	}
}

func (p pomodoroModel) nextState() pomodoroModel {
	p.elapsed = 0
	switch p.state {
	case stateWork:
		if p.currentCycle < p.totalCycles {
			p.state = stateShortBreak
		} else {
			p.state = stateLongBreak
		}
	case stateShortBreak:
		p.currentCycle++
		p.state = stateWork
	case stateLongBreak:
		p.state = stateFinished
	}
	return p
}

func RunPomodoro(totalCycles int, workDuration, shortBreakDuration, longBreakDuration time.Duration) error {
	p := tea.NewProgram(initialPomodoroModel(totalCycles, workDuration, shortBreakDuration, longBreakDuration))
	if _, err := p.Run(); err != nil {
		return err
	}
	fmt.Println("Pomodoro session completed.")
	return nil
}
