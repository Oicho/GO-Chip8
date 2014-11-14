package chip8

import (
	"math/rand"
)

var eightFunctionArray = make([]func(*Memory, uint16), 0xF)
var fFunctionMap = make(map[uint16]func(*Memory, uint16))
var r = rand.New(rand.NewSource(99))

// OneJumpTo is the 1NNN opcode which jump to the NNN address
func OneJumpTo(m *Memory, opcode uint16) {
	m.PC = opcode & 0xFFFF
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
	vx := m.V[(opcode&0x0F00)>>16]
	m.PC += 2
	if vx == byte(opcode&0x00FF) {
		m.PC += 2
	}
}

// FourNeqSkip is the 4XNN opcode
// which skip the next instruction if VX not equals NN
func FourNeqSkip(m *Memory, opcode uint16) {
	vx := m.V[(opcode&0x0F00)>>16]
	m.PC += 2
	if vx != byte(opcode&0x00FF) {
		m.PC += 2
	}
}

// FiveEqSkip is the 5XY0 opcode
// which skip the next instruction if VX not equals NN
func FiveEqSkip(m *Memory, opcode uint16) {
	vx := m.V[(opcode&0x0F00)>>16]
	vy := m.V[(opcode&0x00F0)>>8]
	m.PC += 2
	if vx == vy {
		m.PC += 2
	}
}

// SixSetRegister is the 6XNN opcode
// which set VX to NN
func SixSetRegister(m *Memory, opcode uint16) {
	m.V[(opcode&0x0F00)>>16] = byte(opcode & 0x00FF)
	m.PC += 2
}

// SevenAddToRegister is the 7XNN opcode
// which add NN to VX
func SevenAddToRegister(m *Memory, opcode uint16) {
	sum := uint16(m.V[(opcode&0x0F00)>>16]) + opcode&0x00FF
	if sum > 0x00FF {
		m.VF = true
		//TODO is this the standard behaviour ?
		m.V[(opcode & 0x0F00)] += byte(opcode & 0x00FF)
	} else {
		m.V[(opcode&0x0F00)>>16] = byte(sum)
	}
	m.PC += 2
}

// EightArithmeticOperations is the dispatcher for 8XY? opcodes
func EightArithmeticOperations(m *Memory, opcode uint16) {
	code := opcode & 0x000F
	eightFunctionArray[code](m, opcode)
	m.PC += 2
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
}

// CSetToRandomNumber is the CXNN opcode
// which set VX to a random number and NN
func CSetToRandomNumber(m *Memory, opcode uint16) {
	m.V[(opcode&0x0F00)>>16] = byte((opcode & 0x00FF)) & byte(r.Int63n(0x100))
	m.PC += 2
}

// FUtils is the dispatcher for FNNN opcodes
func FUtils(m *Memory, opcode uint16) {
	fFunctionMap[opcode&0x00FF](m, opcode)
	m.PC += 2
}