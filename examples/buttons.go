package main

import (
	"fmt"
	"time"

	"github.com/vwhitteron/go-pirateaudio/buttons"
)

func main() {
	buttons.OnButtonAPressed(func() {
		fmt.Println("Yo Dawg, A pressed")
	})

	buttons.OnButtonXPressed(func() {
		fmt.Println("Yo Dawg, X pressed")
	})

	buttons.OnButtonYPressed(func() {
		fmt.Println("Yo Dawg, Y pressed")
	})

	buttons.OnButtonBPressed(func() {
		fmt.Println("Yo Dawg, B pressed")
	})

	for {
		time.Sleep(1)
	}
}
