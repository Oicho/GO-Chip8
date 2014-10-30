// You can edit this code!
// Click here and start typing.
package main

import (
	"strings"
	"time"

	"github.com/nsf/termbox-go"
)

// Memory represents the internal memory of the CHIP-8 emulator
type Memory struct {
	I          uint16
	PC         uint16
	SP         uint16
	VF         bool
	DelayTimer byte
	SoundTimer byte
	Key        [16]bool
	Screen     [][]bool
	V          [16]byte
	CallStack  [256]byte
	Memory     [4096]byte
}

// Init must be called when right after you create a Memory variable
// it initialize the screen array and set some values
func (m *Memory) Init() {
	m.PC = 0x200
	m.Screen = make([][]bool, 64)
	for i := range m.Screen {
		m.Screen[i] = make([]bool, 32)
	}
}

// CheckInputs verify if there is a key pressed
// and then set the Key  array accordingly
func (m *Memory) CheckInputs() {
	for i := range m.Key {
		m.Key[i] = false
	}
	const coldef = termbox.ColorDefault
	ev := termbox.PollEvent()
	if ev.Type == termbox.EventKey {
		str := string(ev.Ch)
		if str >= "A" {
			str = strings.ToLower(str)
		}
		switch str {
		case "1":
			m.Key[0] = true
		case "2":
			m.Key[1] = true
		case "3":
			m.Key[2] = true
		case "4":
			m.Key[3] = true
		case "q":
			m.Key[4] = true
		case "w":
			m.Key[5] = true
		case "e":
			m.Key[6] = true
		case "r":
			m.Key[7] = true
		case "a":
			m.Key[8] = true
		case "s":
			m.Key[9] = true
		case "d":
			m.Key[10] = true
		case "f":
			m.Key[11] = true
		case "z":
			m.Key[12] = true
		case "x":
			m.Key[13] = true
		case "c":
			m.Key[14] = true
		case "v":
			m.Key[15] = true
		}
	}

}

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}
func main() {
	var mem = Memory{}
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
			tbprint(0, 10, coldef, coldef, str)
			termbox.Flush()

			if ev.Type == termbox.EventKey && ev.Key == termbox.KeyEsc {
				break loop
			}
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}
