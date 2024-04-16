// Rotate an image pressing the 'A' hardware button
package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"time"

	"github.com/vwhitteron/go-pirateaudio/buttons"
	"github.com/vwhitteron/go-pirateaudio/display"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <img-path>\n", os.Args[0])
		os.Exit(1)
	}

	dsp, err := display.Init()
	if err != nil {
		panic(err)
	}
	defer dsp.Close()

	dsp.FillScreen(color.RGBA{R: 0, G: 0, B: 0, A: 0})

	var rotation display.Rotation
	rotation = 0
	buttons.OnButtonAPressed(func() {
		dsp.FillScreen(color.RGBA{R: 0, G: 0, B: 0, A: 0})
		// Rotate before pushing pixels, so the image appears rotated
		dsp.Rotate(rotation)
		img, err := os.Open(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		defer img.Close()
		dsp.DrawImage(img)
		rotation++
		if rotation > 3 {
			rotation = 0
		}
	})

	for {
		time.Sleep(1)
	}
}
