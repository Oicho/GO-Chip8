// You can edit this code!
// Click here and start typing.
package main

import (
	"github.com/nsf/termbox-go"
	"time"
)

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

func (m *Memory) Init() {
	m.PC = 0x200
	m.Screen = make([][]bool, 64)
	for i, _ := range m.Screen {
		m.Screen[i] = make([]bool, 32)
	}
}

func (m *Memory) CheckInputs() {
	const coldef = termbox.ColorDefault

	ev := termbox.PollEvent()
	if ev.Type == termbox.EventKey {
		switch key := ev.Key; key {
		case termbox.KeyCtrlA:
			termbox.Close()
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

	event_queue := make(chan termbox.Event)
	go func() {
		for {
			event_queue <- termbox.PollEvent()
		}
	}()
	termbox.Flush()
loop:
	for {
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		select {
		case ev := <-event_queue:
			tbprint(0, 10, coldef, coldef, string(ev.Ch))
			termbox.Flush()

			if ev.Type == termbox.EventKey && ev.Key == termbox.KeyEsc {
				break loop
			}
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}
