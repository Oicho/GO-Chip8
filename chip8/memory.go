package chip8

import (
	"fmt"
	"strings"

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
	CallStack  [256]uint16
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

// Fetch get an opcode from memory and then return it
func (m *Memory) Fetch() uint16 {
	opcode := uint16(m.Memory[m.PC]) << 8
	opcode += uint16(m.Memory[m.PC+1])
	m.PC += 2
	return opcode
}

// Decode does stuff
func (m *Memory) Decode(opcode uint16) {
	fmt.Println(opcode)
}

// Iterate does one cycle of a chip8
func (m *Memory) Iterate() {
	opcode := m.Fetch()
	m.Decode(opcode)
	// Decode
	// Execute
}

// CheckInputs verify if there is a key pressed
// and then set the Key  array accordingly
func (m *Memory) CheckInputs() {
	for i := range m.Key {
		m.Key[i] = false
	}
	ev := termbox.PollEvent()
	if ev.Type == termbox.EventKey {
		str := string(ev.Ch)
		if str >= "A" {
			str = strings.ToLower(str)
		}
		switch str {
		case "3":
			m.Key[0] = true
		case "4":
			m.Key[1] = true
		case "5":
			m.Key[2] = true
		case "6":
			m.Key[3] = true
		case "e":
			m.Key[4] = true
		case "r":
			m.Key[5] = true
		case "t":
			m.Key[6] = true
		case "y":
			m.Key[7] = true
		case "d":
			m.Key[8] = true
		case "f":
			m.Key[9] = true
		case "g":
			m.Key[10] = true
		case "h":
			m.Key[11] = true
		case "c":
			m.Key[12] = true
		case "v":
			m.Key[13] = true
		case "b":
			m.Key[14] = true
		case "n":
			m.Key[15] = true
		}
	}
}
