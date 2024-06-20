// Display a rotated image the display
package main

import (
	"fmt"
	"image/color"
	"log"
	"os"

	"github.com/rubiojr/go-pirateaudio/display"
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

	// Set the screen color to white
	dsp.FillScreen(color.RGBA{R: 0, G: 0, B: 0, A: 0})

	img, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer img.Close()

	// Rotate before pushing pixels, so the image appears rotated
	dsp.Rotate(display.ROTATION_180)
	dsp.DrawImage(img)
}
