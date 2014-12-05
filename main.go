// You can edit this code!
// Click here and start typing.
package main

import (
	"strings"
	"time"

	"github.com/Oicho/GO-Chip8/chip8"
	"github.com/Oicho/GO-Chip8/myLogger"

	"github.com/nsf/termbox-go"
)

func main() {
	myLogger.Init(true)
	var mem = chip8.Memory{}
	mem.Init()
	const coldef = termbox.ColorDefault
	err := termbox.Init()
	romPath := "./rom/IBM"
	mem.LoadRom(romPath)
	pause := false
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
			if ev.Type == termbox.EventKey && ev.Key == termbox.KeySpace {
				pause = !pause
			}
			if ev.Type == termbox.EventKey {
				str := string(ev.Ch)
				if str >= "A" {
					str = strings.ToLower(str)
				}
				switch str {
				case "s":
					myLogger.WarningPrint("Now in step by step mode")
					mem.Iterate()
					mem.PrintMemoryValues()
					termbox.Flush()
					break
				case "a":
					myLogger.InfoPrint("Reloading/pausing emulator")
					mem = chip8.Memory{}
					mem.Init()
					mem.LoadRom(romPath)
					pause = true
					termbox.Flush()
					break
				case "q":
					myLogger.InfoPrint("Dump")
					break
				}
			}
		default:
			if !pause {
				mem.Iterate()
				mem.PrintMemoryValues()
				termbox.Flush()
				time.Sleep(10 * time.Millisecond)
			}
		}
	}
}
