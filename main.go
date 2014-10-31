// You can edit this code!
// Click here and start typing.
package main

import (
	"strings"
	"time"

	"github.com/Oicho/GO-Chip8/chip8"
	"github.com/Oicho/GO-Chip8/graphics"

	"github.com/nsf/termbox-go"
)

func printKeys() {

}
func main() {
	var mem = chip8.Memory{}
	mem.Init()

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
			str := string(ev.Ch)
			if str >= "A" {
				str = strings.ToLower(str)
			}
			graphics.Tbprint(0, 10, coldef, coldef, str)

			termbox.Flush()

			if ev.Type == termbox.EventKey && ev.Key == termbox.KeyEsc {
				break loop
			}
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}
