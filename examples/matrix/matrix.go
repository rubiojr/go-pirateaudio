package main

import (
	"time"

	"github.com/rubiojr/go-pirateaudio/textview"
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
