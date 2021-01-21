// Rotate an image pressing the 'A' hardware button
package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"time"

	"github.com/rubiojr/go-pirateaudio/buttons"
	"github.com/rubiojr/go-pirateaudio/st7789"
	"periph.io/x/conn/v3/driver/driverreg"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/host/v3/bcm283x"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <img-path>\n", os.Args[0])
		os.Exit(1)
	}

	if _, err := driverreg.Init(); err != nil {
		log.Fatal(err)
	}

	p, err := spireg.Open("SPI0.1")
	if err != nil {
		log.Fatal(err)
	}
	defer p.Close()

	// Raspberry PI broadcom chipset
	fmt.Println(bcm283x.Present())

	// SPI port we're using
	fmt.Println(p.(spi.Port))

	// USE GPIO9 to send data/commands
	// https://pinout.xyz/pinout/pirate_audio_line_out#
	display, err := st7789.NewSPI(p.(spi.Port), gpioreg.ByName("GPIO9"), &st7789.DefaultOpts)
	if err != nil {
		panic(err)
	}

	display.FillScreen(color.RGBA{R: 0, G: 0, B: 0, A: 0})

	var rotation st7789.Rotation
	rotation = 0
	buttons.OnButtonAPressed(func() {
		display.FillScreen(color.RGBA{R: 0, G: 0, B: 0, A: 0})
		// Rotate before pushing pixels, so the image appears rotated
		display.SetRotation(rotation)
		img, err := os.Open(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		defer img.Close()
		display.DrawImage(img)
		rotation++
		if rotation > 3 {
			rotation = 0
		}
	})

	for {
		time.Sleep(1)
	}
}
