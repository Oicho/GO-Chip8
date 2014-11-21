package chip8

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func createBasicMem() *Memory {
	var mem = &Memory{}
	mem.Init()
	return mem
}

func Test1NNN_OK(t *testing.T) {
	// Adapt
	m := createBasicMem()

	// Act
	OneJumpTo(m, 0x1242)

	// Assert
	assert.Equal(t, uint16(0x242), m.PC, "PC didn't move")
}

func Test3XNN_Skip(t *testing.T) {
	// Adapt
	m := createBasicMem()
	m.V[5] = 0x42

	// Act
	ThreeEqSkip(m, 0x4542)

	// Assert
	assert.Equal(t, uint16(0x204), m.PC, "Haven't skip the instruction")
}

func Test3XNN_NoSkip(t *testing.T) {
	// Adapt
	m := createBasicMem()
	m.V[5] = 0x42

	// Act
	ThreeEqSkip(m, 0x45FF)

	// Assert
	assert.Equal(t, uint16(0x202), m.PC, "Haven't skip the instruction")
}

func Test4XNN_NoSkip(t *testing.T) {
	// Adapt
	m := createBasicMem()
	m.V[5] = 0x42

	// Act
	FourNeqSkip(m, 0x45FF)

	// Assert
	assert.Equal(t, uint16(0x204), m.PC, "Have skip the instruction")
}

func Test4XNN_Skip(t *testing.T) {
	// Adapt
	m := createBasicMem()
	m.V[5] = 0x42

	// Act
	FourNeqSkip(m, 0x4542)

	// Assert
	assert.Equal(t, uint16(0x202), m.PC, "Haven't skip the instruction")
}

func Test5XY0_Skip(t *testing.T) {
	// Adapt
	m := createBasicMem()
	m.V[2] = 0x23
	m.V[3] = 0x23

	// Act
	FiveEqSkip(m, 0x5230)

	// Assert
	assert.Equal(t, uint16(0x204), m.PC, "Haven't skip the instruction")
}

func Test5XY0_NoSkip(t *testing.T) {
	// Adapt
	m := createBasicMem()
	m.V[2] = 0x23
	m.V[3] = 0x42

	// Act
	FiveEqSkip(m, 0x5230)

	// Assert
	assert.Equal(t, uint16(0x202), m.PC, "Skip the instruction")
}

func Test6XNN_SimpleSet(t *testing.T) {
	// Adapt
	m := createBasicMem()
	m.V[2] = 0x23

	// Act
	SixSetRegister(m, 0x6242)

	// Assert
	assert.Equal(t, uint16(0x202), m.PC, "Move to the next instruction")
	assert.Equal(t, byte(0x42), m.V[2], "Set Register")
}

func Test6XNN_SameValueSet(t *testing.T) {
	// Adapt
	m := createBasicMem()
	m.V[2] = 0x20

	// Act
	SixSetRegister(m, 0x6220)

	// Assert
	assert.Equal(t, uint16(0x202), m.PC, "Move to the next instruction")
	assert.Equal(t, byte(0x20), m.V[2], "Set Register")
}

func Test7XNN(t *testing.T) {
	// Adapt
	m := createBasicMem()
	m.V[2] = 0x00

	// Act
	SevenAddToRegister(m, 0x7230)

	// Assert
	assert.Equal(t, uint16(0x202), m.PC, "Move to the next instruction")
	assert.Equal(t, byte(0x30), m.V[2], "Add to Register")
}

func Test7XNN_Overflow(t *testing.T) {
	// Adapt
	m := createBasicMem()
	m.V[2] = 0xFF

	// Act
	SevenAddToRegister(m, 0x7210)

	// Assert
	assert.Equal(t, uint16(0x202), m.PC, "Move to the next instruction")
	assert.Equal(t, byte(0xF), m.V[2], "Add to Register with overflow")
}

func Test9XY0_NoSkip(t *testing.T) {
	// Adapt
	m := createBasicMem()
	m.V[2] = 0x23
	m.V[3] = 0x23

	// Act
	NineNeqSkip(m, 0x5230)

	// Assert
	assert.Equal(t, uint16(0x202), m.PC, "Haven't skip the instruction")
}

func Test9XY0_Skip(t *testing.T) {
	// Adapt
	m := createBasicMem()
	m.V[2] = 0x23
	m.V[3] = 0x42

	// Act
	NineNeqSkip(m, 0x5230)

	// Assert
	assert.Equal(t, uint16(0x204), m.PC, "Skip the instruction")
}

func TestBNNN_with0(t *testing.T) {
	// Adapt
	m := createBasicMem()
	m.V[0] = 0x30

	// Act
	BJumpToV0(m, 0xB000)

	// Assert
	assert.Equal(t, uint16(0x30), m.PC, "Move PC")
}

func TestBNNN_withV0_0(t *testing.T) {
	// Adapt
	m := createBasicMem()
	m.V[0] = 0x0

	// Act
	BJumpToV0(m, 0xB100)

	// Assert
	assert.Equal(t, uint16(0x100), m.PC, "Move PC")
}

func TestBNNN_Overflow(t *testing.T) {
	// Adapt
	m := createBasicMem()
	m.V[0] = 0x4

	// Act
	BJumpToV0(m, 0xBFFF)

	// Assert
	assert.Equal(t, uint16(0x3), m.PC, "Move PC")
}

func TestCXNN(t *testing.T) {
	// Adapt
	m := createBasicMem()
	m.V[0xF] = 0xFF

	// Act
	CSetToRandomNumber(m, 0xCF00)

	// Assert

	assert.Equal(t, uint16(0x202), m.PC, "Move PC")
	assert.Equal(t, 0, m.V[0xF], "And operator")
}
