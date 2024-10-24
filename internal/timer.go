package internal

import (
	"fmt"
	"log"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
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

func initialPomodoroModel(totalCycles, workDuration, shortBreakDuration, longBreakDuration int) pomodoroModel {
	return pomodoroModel{
		totalCycles:        totalCycles,
		currentCycle:       1,
		workDuration:       time.Duration(workDuration) * time.Minute,       // Correctly set work duration
		shortBreakDuration: time.Duration(shortBreakDuration) * time.Minute, // Correctly set short break duration
		longBreakDuration:  15 * time.Minute,                                // Correctly set long break duration
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
			if p.elapsed >= p.currentDuration() {
				p = p.nextState()
				if p.state == stateFinished {
					return p, tea.Quit
				}
			}
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

func (m pomodoroModel) View() string {
	remaining := m.currentDuration() - m.elapsed
	if remaining < 0 {
		remaining = 0
	}

	minutes := int(remaining.Minutes())
	seconds := int(remaining.Seconds()) % 60

	status := "Running"
	if m.paused {
		status = "Paused"
	}

	var stateStr string
	switch m.state {
	case stateWork:
		stateStr = "Work"
	case stateShortBreak:
		stateStr = "Short Break"
	case stateLongBreak:
		stateStr = "Long Break"
	}

	return fmt.Sprintf(
		"%s - Cycle %d/%d\nTime remaining: %02d:%02d [%s]\n\nPress 'p' to pause/resume, Ctrl+C or 'q' to quit.\n",
		stateStr, m.currentCycle, m.totalCycles, minutes, seconds, status)
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
	go func() {
		if err := playSound(); err != nil {
			fmt.Fprintf(os.Stderr, "Error playing sound: %v\n", err)
		}
	}()
	return p
}

func playSound() error {
	f, err := os.Open("Censor Beep Sound.mp3")
	if err != nil {
		return fmt.Errorf("couldn't open file: %s", err)
	}
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	<-done
	return nil
}

func RunPomodoro(totalCycles, workDuration, shortBreakDuration, longBreakDuration int) error {
	p := tea.NewProgram(initialPomodoroModel(totalCycles, workDuration, shortBreakDuration, longBreakDuration))
	if _, err := p.Run(); err != nil {
		return err
	}
	fmt.Println("Pomodoro session completed.")
	return nil
}
