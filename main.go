// You can edit this code!
// Click here and start typing.
package main

import (
	"time"

	"github.com/Oicho/GO-Chip8/chip8"
	"github.com/Oicho/GO-Chip8/graphics"
	"github.com/nsf/termbox-go"
)

func main() {
	var mem = chip8.Memory{}
	mem.Init()
	mem.Screen[30][1] = true
	mem.Screen[31][12] = true
	const coldef = termbox.ColorDefault
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()
	termbox.Flush()
loop:
	for {
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		select {
		case ev := <-eventQueue:
			if ev.Type == termbox.EventKey && ev.Key == termbox.KeyEsc {
				break loop
			}
		default:
			graphics.PrintScreen(mem.Screen)
			termbox.Flush()
			time.Sleep(10 * time.Millisecond)
		}
	}
}
