package chip8

import (
	"fmt"
	"os"
	"strings"

	"github.com/nsf/termbox-go"
)

// Memory represents the internal memory of the CHIP-8 emulator
type Memory struct {
	I          uint16
	PC         uint16
	SP         uint16
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

// LoadRom load a rom in the memory
func (m *Memory) LoadRom(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		// TODO log
		return err
	}

	defer file.Close()
	data := make([]byte, 0xE00)
	nbBytes, err := file.Read(data)

	if err != nil {
		// TODO log
		return err
	}

	for i := 0; i < nbBytes; i++ {
		m.Memory[i+0x200] = data[i]
	}

	return nil
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
	// TODO call fun array
	fmt.Println(opcode)
}

// Iterate does one cycle of a chip8
func (m *Memory) Iterate() {
	opcode := m.Fetch()
	m.Decode(opcode)
	// Decode
	// Execute
}

// WaitForInput wait for an input and then returns it
func (m *Memory) WaitForInput() byte {
	for {
		b, i := m.CheckInputs()
		if b {
			return i
		}
		// TODO sleep ?
	}
}

// CheckInputs verify if there is a key pressed
// and then set the Key  array accordingly
func (m *Memory) CheckInputs() (bool, byte) {

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
			return true, 0
		case "4":
			m.Key[1] = true
			return true, 1
		case "5":
			m.Key[2] = true
			return true, 2
		case "6":
			m.Key[3] = true
			return true, 3
		case "e":
			m.Key[4] = true
			return true, 4
		case "r":
			m.Key[5] = true
			return true, 5
		case "t":
			m.Key[6] = true
			return true, 6
		case "y":
			m.Key[7] = true
			return true, 7
		case "d":
			m.Key[8] = true
			return true, 8
		case "f":
			m.Key[9] = true
			return true, 9
		case "g":
			m.Key[10] = true
			return true, 10
		case "h":
			m.Key[11] = true
			return true, 11
		case "c":
			m.Key[12] = true
			return true, 12
		case "v":
			m.Key[13] = true
			return true, 13
		case "b":
			m.Key[14] = true
			return true, 14
		case "n":
			m.Key[15] = true
			return true, 15
		}
	}
	return false, 0
}
