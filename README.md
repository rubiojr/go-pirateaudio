## Pirate Audio Go Module

Go module to control Pimoroni's Pirate Audio LCD and buttons.

⚠️   Highly Experimental, API subject to change, may set your device on fire ⚠️ 

![gadget.jpg](gadget.jpg)

Head over to [the official Pirate Audio GitHub repository](https://github.com/pimoroni/pirate-audio) to configure your Raspberry to be able to use this library. You'll need the `config.txt` changes mentioned in the README.

## ST7789

The driver package for the 240x240px [Pirate Audio display](https://shop.pimoroni.com/products/pirate-audio-headphone-amp).

Heavily based on the [TinyGo](https://github.com/tinygo-org/drivers/tree/e376785596dc8269f3e8aa42a9bf75fb1457febc/st7789) driver, modified to work with mainline Go, using [periph.io](https://periph.io) to interface with the Raspberry PI (SPI/GPIO).

Also used the [Python driver](https://github.com/pimoroni/st7789-python) by [Philip Howard](https://github.com/Gadgetoid) as a reference.

### Drawing an Image on the display

```Go
// Display a rotated image the display
package main

import (
	"fmt"
	"image/color"
	"log"
	"os"

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
```

### Controlling the hardware buttons (A,B,X,Y)

```Go
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
```

### Combining HW buttons and image display

![](images/rotate.gif)

```Go
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
```

### Drawing text

![](docs/images/matrix.gif)

```Go
package main

import (
	"time"

	"github.com/vwhitteron/go-pirateaudio/textview"
)

func main() {
	opts := textview.DefaultOpts
	opts.FGColor = textview.GREEN
	tv := textview.NewWithOptions(opts)
	tv.Draw("")
	time.Sleep(3 * time.Second)
	tv.DrawChars("Wake up, Neo...")
	time.Sleep(3 * time.Second)
	tv.DrawChars("The Matrix has you...")
	time.Sleep(3 * time.Second)
	tv.DrawChars("Follow the white rabbit.")
}
```
