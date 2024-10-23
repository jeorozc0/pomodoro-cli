/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/jeorozc0/pomodoro-cli/cmd"
)

func main() {
	sampleRate := beep.SampleRate(44100) // Common sample rate
	speaker.Init(sampleRate, sampleRate.N(time.Second/10))
	cmd.Execute()
}
