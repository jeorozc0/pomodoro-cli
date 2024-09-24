package cmd

import (
	"fmt"
	"time"

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
	workDuration, err := parseTime(workTimeSelection)
	if err != nil {
		return fmt.Errorf("invalid work time: %w", err)
	}

	restDuration, err := parseTime(restTimeSelection)
	if err != nil {
		return fmt.Errorf("invalid rest time: %w", err)
	}

	// Clear screen (optional)
	clearScreen()

	// Run work timer
	fmt.Printf("Work Timer: %s\n\n", workDuration)
	if err := internal.RunTimer(workDuration); err != nil {
		return fmt.Errorf("failed to run work timer: %w", err)
	}

	// Clear screen (optional)
	clearScreen()

	// Run rest timer
	fmt.Printf("Rest Timer: %s\n\n", restDuration)
	if err := internal.RunTimer(restDuration); err != nil {
		return fmt.Errorf("failed to run rest timer: %w", err)
	}

	fmt.Println("Timer sessions completed.")
	return nil
}

func parseTime(selection string) (time.Duration, error) {
	switch selection {
	case "Short":
		return 5 * time.Second, nil
	case "Medium":
		return 10 * time.Second, nil
	case "Long":
		return 15 * time.Second, nil
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
