package main

import (
	"github.com/Oicho/GO-Chip8/chip8"
	"github.com/Oicho/GO-Chip8/myLogger"
	termbox "github.com/nsf/termbox-go"

	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Usage: program filepath")
		return
	}
	myLogger.Init(true)
	var mem = chip8.Memory{}
	mem.Init()

	err := termbox.Init()
	var romPath string
	var pause bool

	if (len(os.Args) == 3) {
		pause = true
		romPath = os.Args[2]
	} else {
		pause = false
		romPath = os.Args[1]
	}
	mem.LoadRom(romPath)
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
					mem.PrintMemoryValues()
					mem.Iterate()
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
				mem.PrintMemoryValues()
				mem.Iterate()
				termbox.Flush()
				time.Sleep(10 * time.Millisecond)
			}
		}
	}
}
