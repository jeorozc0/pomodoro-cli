package cmd

import (
	"fmt"

	"github.com/jeorozc0/pomodoro-cli/internal"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pomodoro",
	Short: "A simple pomodor tool for the terminal",
	RunE:  rootRunE,
}

func rootRunE(cmd *cobra.Command, args []string) error {
	// Run the menu
	workTimeSelection, restTimeSelection, err := internal.RunMenu()
	if err != nil {
		return fmt.Errorf("failed to run menu: %w", err)
	}

	// Parse selections
	workDuration, err := parseWork(workTimeSelection)
	if err != nil {
		return fmt.Errorf("invalid work time: %w", err)
	}

	restDuration, err := parseRest(restTimeSelection)
	if err != nil {
		return fmt.Errorf("invalid rest time: %w", err)
	}

	// Clear screen (optional)
	clearScreen()

	// Run work timer
	fmt.Println("Starting timer")
	if err := internal.RunPomodoro(4, // totalCycles
		workDuration, // workDuration
		restDuration, // shortBreakDuration
		15); err != nil {
		return fmt.Errorf("failed to run work timer: %w", err)
	}

	fmt.Println("Timer sessions completed.")
	return nil
}

func parseWork(selection string) (int, error) {
	switch selection {
	case "15 mins":
		return 1, nil
	case "25 mins":
		return 25, nil
	case "45 mins":
		return 45, nil
	default:
		return 0, fmt.Errorf("unknown time selection: %s", selection)
	}
}

func parseRest(selection string) (int, error) {
	switch selection {
	case "5 mins":
		return 1, nil
	case "10 mins":
		return 10, nil
	case "15 mins":
		return 15, nil
	default:
		return 0, fmt.Errorf("unknown time selection: %s", selection)
	}
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func Execute() error {
	return rootCmd.Execute()
}
