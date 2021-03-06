package chip8

import (
	"github.com/Oicho/GO-Chip8/graphics"
	"github.com/Oicho/GO-Chip8/myLogger"
	termbox "github.com/nsf/termbox-go"
	"os"
	"strconv"
	"strings"
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

var chip8Fontset = [80]byte{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80} // F

// Init must be called when right after you create a Memory variable
// it initialize the screen array and set some values
func (m *Memory) Init() {
	myLogger.InfoPrint("Initiating a chip8")
	for i := 0; i < 80; i++ {
		m.Memory[i] = chip8Fontset[i]
	}
	m.PC = 0x200
	m.Screen = make([][]bool, 64)
	for i := range m.Screen {
		m.Screen[i] = make([]bool, 32)
	}
	myLogger.Info.Print("Finished init chip8")
}

// LoadRom load a rom in the memory
func (m *Memory) LoadRom(filePath string) error {
	myLogger.InfoPrint("Loading a ROM")
	file, err := os.Open(filePath)
	if err != nil {
		myLogger.Error.Println("file:<" + filePath + "> not found")
		return err
	}

	defer file.Close()
	data := make([]byte, 0xE00)
	nbBytes, err := file.Read(data)

	if err != nil {
		myLogger.Error.Println("Couldn't read the file")
		return err
	}
	for i := 0; i < nbBytes; i++ {
		m.Memory[i+0x200] = data[i]
	}
	myLogger.Info.Println("ROM loading done")
	return nil
}

// Fetch get an opcode from memory and then return it
func (m *Memory) Fetch() uint16 {
	opcode := uint16(m.Memory[m.PC]) << 8
	opcode += uint16(m.Memory[m.PC+1])
	return opcode
}

// Decode does stuff
func (m *Memory) Decode(opcode uint16) {
	mainFunctionArray[(0xF000&opcode)>>12](m, opcode)
}

// Iterate does one cycle of a chip8
func (m *Memory) Iterate() {
	myLogger.InfoVerbosePrint("Fetching at 0x" + myLogger.Uint16ToString(m.PC))
	opcode := m.Fetch()
	myLogger.InfoVerbosePrint("Executing 0x"+ myLogger.Uint16ToString(opcode))
	m.Decode(opcode)
}

// WaitForInput wait for an input and then returns it
func WaitForInput(m *Memory) byte {
	for {
		b, i := CheckInputs(m)
		if b {
			return i
		}
		// TODO sleep ?
	}
}

// CheckInputs verify if there is a key pressed
// and then set the Key  array accordingly
var CheckInputs = func (m *Memory) (bool, byte) {
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

// PrintMemoryValues print chip8 state value
func (m *Memory) PrintMemoryValues() {
	graphics.PrintScreen(m.Screen)
	height := 0
	width := 130
	graphics.PrintString(width,
		height,
		termbox.ColorDefault,
		termbox.ColorDefault,
		"PC="+myLogger.Uint16ToString(m.PC))
	height++
	graphics.PrintString(width,
		height,
		termbox.ColorDefault,
		termbox.ColorDefault,
		"I="+myLogger.Uint16ToString(m.I))
	height++
	for i := 0; i < 0x10; i++ {
		graphics.PrintString(width,
			height,
			termbox.ColorDefault,
			termbox.ColorDefault,
			"V["+strconv.Itoa(i)+"]="+myLogger.ByteToString(m.V[i]))
		height++
	}

}
