package buttons

import (
	"fmt"
	"log"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"
)

func init() {
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}
}

func OnButtonAPressed(fn func()) {
	onButtonPressed(5, fn)
}

func OnButtonBPressed(fn func()) {
	onButtonPressed(6, fn)
}

func OnButtonXPressed(fn func()) {
	onButtonPressed(16, fn)
}

func OnButtonYPressed(fn func()) {
	onButtonPressed(24, fn)
}

func onButtonPressed(n int, fn func()) {
	go func() {
		p := gpioreg.ByName(fmt.Sprintf("GPIO%d", n))
		if err := p.In(gpio.PullUp, gpio.FallingEdge); err != nil {
			log.Fatal(err)
		}
		for {
			p.WaitForEdge(-1)
			fn()
		}
	}()
}
