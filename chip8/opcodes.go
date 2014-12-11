package chip8

import (
	"math/rand"

	"github.com/Oicho/GO-Chip8/graphics"
)

var mainFunctionArray = [0x10]func(*Memory, uint16){ZeroDispatcher, OneJumpTo, TwoCallSubRoutine, ThreeEqSkip, FourNeqSkip, FiveEqSkip, SixSetRegister, SevenAddToRegister, EightDispatcher, NineNeqSkip, ASetAddressRegister, BJumpToV0, CSetToRandomNumber, DWrapsOnScreen, EDispatcher, FDispatcher}
var eightFunctionArray = [0xF]func(*Memory, uint16){EightZeroSet, EightOneORSet, EightTwoANDSet, EightThreeXORSet, EightFourAdd, EightFiveSub, EightSixRightShift, EightSevenMinus, nil, nil, nil, nil, nil, nil, EightFourteenLeftShift}
var fFunctionMap = map[uint16]func(*Memory, uint16){7: FSetVXtoDelayTimer, 0x0A: FWaitKeyPress, 0x15: FSetDelayTimerToVX, 0x18: FSetSoundTimerToVX, 0x1E: FAddVXToI, 0x29: FGoToSprite, 0x33: FBCD, 0x55: FWriteMemory, 0x65: FReadMemory}
var r = rand.New(rand.NewSource(99))

// ZeroDispatcher is the 0??? opcodes dispatcher
func ZeroDispatcher(m *Memory, opcode uint16) {
	switch opcode {
	case 0x00E0:
		ZeroClearScreen(m, opcode)
	case 0x00EE:
		ZeroReturnFromSubRoutine(m, opcode)
	default:
		m.PC += 2
		// TODO show something funny
	}
}

// ZeroClearScreen is 00E0 opcode
// which clear the screen
func ZeroClearScreen(m *Memory, opcode uint16) {
	for i, arr := range m.Screen {
		for j := range arr {
			m.Screen[i][j] = false
		}
	}
	graphics.PrintScreen(m.Screen)
	m.PC += 2
}

// ZeroReturnFromSubRoutine is the 00EE opcode
// which return from a subroutine
func ZeroReturnFromSubRoutine(m *Memory, opcode uint16) {
	m.SP--
	m.PC = m.CallStack[m.SP]
}

// OneJumpTo is the 1NNN opcode which jump to the NNN address
func OneJumpTo(m *Memory, opcode uint16) {
	m.PC = opcode & 0x0FFF
}

// TwoCallSubRoutine is the 2NNN opcode
// which call the subroutine at the NNN address
func TwoCallSubRoutine(m *Memory, opcode uint16) {
	m.CallStack[m.SP] = m.PC
	m.SP++
	m.PC = opcode & 0x0FFF
}

// ThreeEqSkip is the 3XNN opcode
// which skip the next instruction if VX equals NN
func ThreeEqSkip(m *Memory, opcode uint16) {
	vx := m.V[(opcode&0x0F00)>>8]
	if vx == byte(opcode&0x00FF) {
		m.PC += 4
	} else {
		m.PC += 2
	}
}

// FourNeqSkip is the 4XNN opcode
// which skips the next instruction if VX not equals NN
func FourNeqSkip(m *Memory, opcode uint16) {
	vx := m.V[(opcode&0x0F00)>>8]
	if vx != byte(opcode&0x00FF) {
		m.PC += 4
	} else {
		m.PC += 2
	}
}

// FiveEqSkip is the 5XY0 opcode
// which skips the next instruction if VX equals VY.
func FiveEqSkip(m *Memory, opcode uint16) {
	vx := m.V[(opcode&0x0F00)>>8]
	vy := m.V[(opcode&0x00F0)>>4]
	if vx == vy {
		m.PC += 4
	} else {
		m.PC += 2
	}
}

// SixSetRegister is the 6XNN opcode
// which set VX to NN
func SixSetRegister(m *Memory, opcode uint16) {
	m.V[(opcode&0x0F00)>>8] = byte(opcode & 0x00FF)
	m.PC += 2
}

// SevenAddToRegister is the 7XNN opcode
// which add NN to VX
func SevenAddToRegister(m *Memory, opcode uint16) {
	m.V[(opcode&0x0F00)>>8] += byte(opcode & 0x00FF)
	m.PC += 2
}

// EightDispatcher is the dispatcher for 8XY? opcodes
func EightDispatcher(m *Memory, opcode uint16) {
	code := opcode & 0x000F
	eightFunctionArray[code](m, opcode)
	m.PC += 2
}

// -------------- 8XY? opcde-----------\\

// EightZeroSet is the 8XY0 opcode
// which sets VX to the value of VY
func EightZeroSet(m *Memory, opcode uint16) {
	m.V[(opcode&0x0F00)>>8] = m.V[(opcode&0x00F0)>>4]
}

// EightOneORSet is the 8XY1 opcode
// which sets VX to VX OR VY
func EightOneORSet(m *Memory, opcode uint16) {
	m.V[(opcode&0x0F00)>>8] = m.V[(opcode&0x0F00)>>8] | m.V[(opcode&0x00F0)>>4]
}

// EightTwoANDSet is the 8XY2 opcode
// which sets VX to VX AND VY
func EightTwoANDSet(m *Memory, opcode uint16) {
	m.V[(opcode&0x0F00)>>8] = m.V[(opcode&0x0F00)>>8] & m.V[(opcode&0x00F0)>>4]
}

// EightThreeXORSet is the 8XY3 opcode
// which sets VX to VX XOR VY
func EightThreeXORSet(m *Memory, opcode uint16) {
	m.V[(opcode&0x0F00)>>8] = m.V[(opcode&0x0F00)>>8] ^ m.V[(opcode&0x00F0)>>4]
}

// EightFourAdd is the 8XY4 opcode
// which Adds VY to VX. VF is set to 1 when there's a carry, and to 0 when there isn't
func EightFourAdd(m *Memory, opcode uint16) {
	x, y := xyExtractor(opcode)
	if m.V[y] > 0xff-m.V[x] {
		m.V[0xF] = 1
	} else {
		m.V[0xF] = 0
	}
	m.V[x] += m.V[y]
}

// EightFiveSub is the 8XY5 opcode
// which set VX to VX-VY
func EightFiveSub(m *Memory, opcode uint16) {
	x, y := xyExtractor(opcode)
	if m.V[x] > m.V[y] {
		m.V[0xF] = 1
	} else {
		m.V[0xF] = 0
	}
	m.V[x] = m.V[x] - m.V[y]
}

// EightSixRightShift is the 8XY6 opcode
// which shifts VX right by one
func EightSixRightShift(m *Memory, opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	m.V[0xF] = 1 & m.V[x]
	m.V[x] = m.V[x] >> 1
}

// EightSevenMinus is the 8XY7 opcode
// which set VX to VY-VX
func EightSevenMinus(m *Memory, opcode uint16) {
	x, y := xyExtractor(opcode)
	if m.V[y] > m.V[x] {
		m.V[0xF] = 1
	} else {
		m.V[0xF] = 0
	}
	m.V[x] = m.V[y] - m.V[x]
}

// EightFourteenLeftShift is the 8XYE opcode
// which shifts VX left by one
func EightFourteenLeftShift(m *Memory, opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	if 0x80&m.V[x] == 0 {
		m.V[0xF] = 0
	} else {
		m.V[0xF] = 1
	}
	m.V[x] = m.V[x] << 1
}

// NineNeqSkip is the 9XY0 opcode
// which skips the next instruction if VX doesn't equal VY
func NineNeqSkip(m *Memory, opcode uint16) {
	if m.V[(opcode&0x0F00)>>8] != m.V[(opcode&0x00F0)>>4] {
		m.PC += 4
	} else {
		m.PC += 2
	}
}

// ASetAddressRegister is the ANNN opcode
// which set the Address register I to NNN
func ASetAddressRegister(m *Memory, opcode uint16) {
	m.I = opcode & 0x0FFF
	m.PC += 2
}

// BJumpToV0 is the BNNN opcode
// which jump to the address V0 + NNN
func BJumpToV0(m *Memory, opcode uint16) {
	m.PC = uint16(m.V[0]) + (opcode & 0x0FFF)
	if m.PC >= 0x1000 {
		m.PC = m.PC - 0x1000
	}
}

// CSetToRandomNumber is the CXNN opcode
// which set VX to a random number and NN
func CSetToRandomNumber(m *Memory, opcode uint16) {
	m.V[(opcode&0x0F00)>>8] = byte((opcode & 0x00FF)) & byte(r.Int63n(0x100))
	m.PC += 2
}

// DWrapsOnScreen is the DXYN opcode
// which draw sprites
func DWrapsOnScreen(m *Memory, opcode uint16) {
	x, y := xyExtractor(opcode)
	vx := uint16(m.V[x])
	vy := uint16(m.V[y])
	height := 0x000F & opcode
	m.V[0xF] = 0
	for py := uint16(0); py < height; py++ {
		pixel := m.Memory[m.I+py]
		for px := uint16(0); px < 8; px++ {
			if (pixel & (0x80 >> px)) != 0 {
				if m.Screen[px+vx][py+vy] {
					m.V[0xF] = 1
				}
				m.Screen[px+vx][py+vy] = !m.Screen[px+vx][py+vy]
			}
		}
	}
	graphics.PrintScreen(m.Screen)
	m.PC += 2
}

// EDispatcher is the E??? opcodes dispatcher
func EDispatcher(m *Memory, opcode uint16) {
	switch opcode & 0x00FF {
	case 0x9E:
		ESkipIfKeyPress(m, opcode)
	case 0xA1:
		ESkipIfKeyNotPress(m, opcode)
	default:
		// TODO show error
	}
	m.PC += 2
}

// ESkipIfKeyPress is the EX9E opcode
// which skip the next instruction if the key stored in VX is pressed
func ESkipIfKeyPress(m *Memory, opcode uint16) {
	if m.Key[m.V[(opcode&0x0F00)>>8]] {
		m.PC += 2
	}
}

// ESkipIfKeyNotPress is the EXA1 opcode
// which skip the next instruction if the key stored in VX is not pressed
func ESkipIfKeyNotPress(m *Memory, opcode uint16) {
	if !m.Key[m.V[(opcode&0x0F00)>>8]] {
		m.PC += 2
	}
}

// FDispatcher is the dispatcher for FNNN opcodes
func FDispatcher(m *Memory, opcode uint16) {
	fFunctionMap[opcode&0x00FF](m, opcode)
	m.PC += 2
}

// FSetVXtoDelayTimer is the FX07 opcode
// which sets VX to the value of the delay timer
func FSetVXtoDelayTimer(m *Memory, opcode uint16) {
	m.V[(opcode&0x0F00)>>8] = m.DelayTimer
}

// FWaitKeyPress is the FX0A opcode
// which wait a key press and then stores it in VX
func FWaitKeyPress(m *Memory, opcode uint16) {
	m.V[(opcode&0x0F00)>>8] = m.WaitForInput()
}

// FSetDelayTimerToVX is the FX15 opcode
// which sets the delay timer to VX
func FSetDelayTimerToVX(m *Memory, opcode uint16) {
	m.DelayTimer = m.V[(opcode&0x0F00)>>8]
}

// FSetSoundTimerToVX is the FX18 opcode
// which sets the sound timer to VX
func FSetSoundTimerToVX(m *Memory, opcode uint16) {
	m.SoundTimer = m.V[(opcode&0x0F00)>>8]
}

// FAddVXToI is the FX1E opcode
// which adds VX to I
func FAddVXToI(m *Memory, opcode uint16) {
	m.I += uint16(m.V[(opcode&0x0F00)>>8])
}

// FGoToSprite is the FX29 opcode
// which sets I to the location of the sprite for the character in VX
func FGoToSprite(m *Memory, opcode uint16) {
	m.I = uint16(5) * uint16(m.V[(0x0F00&opcode)>>8])
}

// FBCD is the FX33 opcode
// which stores the Binary-coded decimal representation of VX,
// with the most significant of three digits at the address in I,
// the middle digit at I plus 1,
// and the least significant digit at I plus 2.
func FBCD(m *Memory, opcode uint16) {
	vx := m.V[(0x0F00&opcode)>>8]
	m.Memory[m.I] = vx / 100
	m.Memory[m.I+1] = (vx / 10) % 10
	m.Memory[m.I+2] = (vx % 100) % 10
}

// FWriteMemory is the FX55 opcode
// which stores V0 to VX in memory starting at address I
func FWriteMemory(m *Memory, opcode uint16) {
	vx := (opcode & 0x0F00) >> 8
	for p := uint16(0); p <= vx; p++ {
		m.Memory[m.I+p] = m.V[p]
	}
}

// FReadMemory is the FX65 opcode
// which fills V0 to VX with values from memory starting at address I
func FReadMemory(m *Memory, opcode uint16) {
	vx := (opcode & 0x0F00) >> 8
	for p := uint16(0); p <= vx; p++ {
		m.V[p] = m.Memory[m.I+p]
	}
}

func xyExtractor(opcode uint16) (x uint16, y uint16) {
	x = (opcode & 0x0F00) >> 8
	y = (opcode & 0x00F0) >> 4
	return
}
