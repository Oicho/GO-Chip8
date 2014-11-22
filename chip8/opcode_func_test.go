package chip8

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func createBasicMem() *Memory {
	var mem = &Memory{}
	mem.Init()
	return mem
}

// TODO check init func

func TestClearScreen(t *testing.T) {
	// Adapt
	m := createBasicMem()
	m.Screen[0][0] = true
	m.Screen[3][5] = true

	// Act
	Dispatch(m, 0x00E0)

	// Assert
	assert.False(t, m.Screen[0][0], "Reset screen")
	assert.False(t, m.Screen[3][5], "Reset screen")
	assert.Equal(t, uint16(0x202), m.PC, "Go to the next instruction")
}

func TestReturnFromSubRoutine(t *testing.T) {
	// Adapt
	m := createBasicMem()
	m.SP = 1
	m.CallStack[0] = 0x321

	// Act
	Dispatch(m, 0x00EE)

	// Assert
	assert.Equal(t, uint16(0x321), m.PC, "Return from subroutine")
	assert.Equal(t, 0, m.SP, "Decrement stack pointer")
}

func TestRca(t *testing.T) {
	// Adapt
	m := createBasicMem()

	// Act
	Dispatch(m, 0x0442)

	// Assert
	assert.Equal(t, uint16(0x202), m.PC, "Go to the next instruction")
}

func Test1NNN_OK(t *testing.T) {
	// Adapt
	m := createBasicMem()

	// Act
	Dispatch(m, 0x1242)

	// Assert
	assert.Equal(t, uint16(0x242), m.PC, "PC didn't move")
}

func Test2NNN_OK(t *testing.T) {
	// Adapt
	m := createBasicMem()

	// Act
	Dispatch(m, 0x2FFF)

	// Assert
	assert.Equal(t, 1, m.SP, "Increment stack pointer")
	assert.Equal(t, uint16(0x200), m.CallStack[0], "Return position")
	assert.Equal(t, uint16(0xFFF), m.PC, "PC didn't move")
}

func Test3XNN_Skip(t *testing.T) {
	// Adapt
	m := createBasicMem()
	m.V[5] = 0x42

	// Act
	Dispatch(m, 0x3542)

	// Assert
	assert.Equal(t, uint16(0x204), m.PC, "Haven't skip the instruction")
}

func Test3XNN_NoSkip(t *testing.T) {
	// Adapt
	m := createBasicMem()
	m.V[5] = 0x42

	// Act
	Dispatch(m, 0x35FF)

	// Assert
	assert.Equal(t, uint16(0x202), m.PC, "Haven't skip the instruction")
}

func Test4XNN_NoSkip(t *testing.T) {
	// Adapt
	m := createBasicMem()
	m.V[5] = 0x42

	// Act
	Dispatch(m, 0x45FF)

	// Assert
	assert.Equal(t, uint16(0x204), m.PC, "Have skip the instruction")
}

func Test4XNN_Skip(t *testing.T) {
	// Adapt
	m := createBasicMem()
	m.V[5] = 0x42

	// Act
	Dispatch(m, 0x4542)

	// Assert
	assert.Equal(t, uint16(0x202), m.PC, "Haven't skip the instruction")
}

func Test5XY0_Skip(t *testing.T) {
	// Adapt
	m := createBasicMem()
	m.V[2] = 0x23
	m.V[3] = 0x23

	// Act
	Dispatch(m, 0x5230)

	// Assert
	assert.Equal(t, uint16(0x204), m.PC, "Haven't skip the instruction")
}

func Test5XY0_NoSkip(t *testing.T) {
	// Adapt
	m := createBasicMem()
	m.V[2] = 0x23
	m.V[3] = 0x42

	// Act
	Dispatch(m, 0x5230)

	// Assert
	assert.Equal(t, uint16(0x202), m.PC, "Skip the instruction")
}

func Test6XNN_SimpleSet(t *testing.T) {
	// Adapt
	m := createBasicMem()
	m.V[2] = 0x23

	// Act
	Dispatch(m, 0x6242)

	// Assert
	assert.Equal(t, uint16(0x202), m.PC, "Move to the next instruction")
	assert.Equal(t, byte(0x42), m.V[2], "Set Register")
}

func Test6XNN_SameValueSet(t *testing.T) {
	// Adapt
	m := createBasicMem()
	m.V[2] = 0x20

	// Act
	Dispatch(m, 0x6220)

	// Assert
	assert.Equal(t, uint16(0x202), m.PC, "Move to the next instruction")
	assert.Equal(t, byte(0x20), m.V[2], "Set Register")
}

func Test7XNN(t *testing.T) {
	// Adapt
	m := createBasicMem()
	m.V[2] = 0x00

	// Act
	Dispatch(m, 0x7230)

	// Assert
	assert.Equal(t, uint16(0x202), m.PC, "Move to the next instruction")
	assert.Equal(t, byte(0x30), m.V[2], "Add to Register")
}

func Test7XNN_Overflow(t *testing.T) {
	// Adapt
	m := createBasicMem()
	m.V[2] = 0xFF

	// Act
	Dispatch(m, 0x7210)

	// Assert
	assert.Equal(t, uint16(0x202), m.PC, "Move to the next instruction")
	assert.Equal(t, byte(0xF), m.V[2], "Add to Register with overflow")
}

// TODO remove me
func Test8(t *testing.T) {
	// Adapt
	m := createBasicMem()

	// Act
	Dispatch(m, 0x8)

	// Assert
	assert.Equal(t, uint16(0x202), m.PC, "Move to the next instruction")
}

func Test8XY0(t *testing.T) {
	// Adapt
	m := createBasicMem()
	m.V[0x3] = 0x32
	m.V[0xF] = 0xA0

	// Act
	Dispatch(m, 0x83F0)

	// Assert
	assert.Equal(t, uint16(0xA0), m.V[0x3], "Changed VX")
	assert.Equal(t, uint16(0xA0), m.V[0xF], "Unchanged VY")
	assert.Equal(t, uint16(0x202), m.PC, "Move to the next instruction")
}

func Test8XY1_FF(t *testing.T) {
	// Adapt
	m := createBasicMem()
	m.V[0x4] = 0xF0
	m.V[0x5] = 0x0F

	// Act
	Dispatch(m, 0x8451)

	// Assert
	assert.Equal(t, uint16(0xFF), m.V[0x4], "Changed VX")
	assert.Equal(t, uint16(0x0F), m.V[0x5], "Unchanged VY")
	assert.Equal(t, uint16(0x202), m.PC, "Move to the next instruction")
}

func Test8XY1_0F(t *testing.T) {
	// Adapt
	m := createBasicMem()
	m.V[0x4] = 0x00
	m.V[0x5] = 0x0F

	// Act
	Dispatch(m, 0x8451)

	// Assert
	assert.Equal(t, uint16(0x0F), m.V[0x4], "Changed VX")
	assert.Equal(t, uint16(0x0F), m.V[0x5], "Unchanged VY")
	assert.Equal(t, uint16(0x202), m.PC, "Move to the next instruction")
}

func Test8XY2_0(t *testing.T) {
	// Adapt
	m := createBasicMem()
	m.V[0x2] = 0x0F
	m.V[0xF] = 0xF0

	// Act
	Dispatch(m, 0x82F2)

	// Assert
	assert.Equal(t, uint16(0x00), m.V[0x2], "Changed VX")
	assert.Equal(t, uint16(0xF0), m.V[0xF], "Unchanged VY")
	assert.Equal(t, uint16(0x202), m.PC, "Move to the next instruction")
}

func Test8XY2_0F(t *testing.T) {
	// Adapt
	m := createBasicMem()
	m.V[0x2] = 0x0F
	m.V[0xF] = 0xFF

	// Act
	Dispatch(m, 0x82F2)

	// Assert
	assert.Equal(t, uint16(0x0F), m.V[0x2], "Changed VX")
	assert.Equal(t, uint16(0xFF), m.V[0xF], "Unchanged VY")
	assert.Equal(t, uint16(0x202), m.PC, "Move to the next instruction")
}

func Test9XY0_NoSkip(t *testing.T) {
	// Adapt
	m := createBasicMem()
	m.V[2] = 0x23
	m.V[3] = 0x23

	// Act
	Dispatch(m, 0x9230)

	// Assert
	assert.Equal(t, uint16(0x202), m.PC, "Haven't skip the instruction")
}

func Test9XY0_Skip(t *testing.T) {
	// Adapt
	m := createBasicMem()
	m.V[2] = 0x23
	m.V[3] = 0x42

	// Act
	Dispatch(m, 0x9230)

	// Assert
	assert.Equal(t, uint16(0x204), m.PC, "Skip the instruction")
}

func TestANNN(t *testing.T) {
	// Adapt
	m := createBasicMem()
	m.I = 42

	// Act
	Dispatch(m, 0xA153)

	// Assert
	assert.Equal(t, uint16(0x153), m.I, "Set I")
}

func TestBNNN_with0(t *testing.T) {
	// Adapt
	m := createBasicMem()
	m.V[0] = 0x30

	// Act
	Dispatch(m, 0xB000)

	// Assert
	assert.Equal(t, uint16(0x30), m.PC, "Move PC")
}

func TestBNNN_withV0_0(t *testing.T) {
	// Adapt
	m := createBasicMem()
	m.V[0] = 0x0

	// Act
	Dispatch(m, 0xB100)

	// Assert
	assert.Equal(t, uint16(0x100), m.PC, "Move PC")
}

func TestBNNN_Overflow(t *testing.T) {
	// Adapt
	m := createBasicMem()
	m.V[0] = 0x4

	// Act
	Dispatch(m, 0xBFFF)

	// Assert
	assert.Equal(t, uint16(0x3), m.PC, "Move PC")
}

func TestCXNN(t *testing.T) {
	// Adapt
	m := createBasicMem()
	m.V[0xF] = 0xFF

	// Act
	Dispatch(m, 0xCF00)

	// Assert
	assert.Equal(t, uint16(0x202), m.PC, "Move PC")
	assert.Equal(t, 0, m.V[0xF], "And operator")
}

func TestENNN(t *testing.T) {
	// Adapt
	m := createBasicMem()

	// Act
	Dispatch(m, 0xEFFF)

	// Assert
	assert.Equal(t, uint16(0x202), m.PC, "Move to the next instruction")
}

func TestFX07(t *testing.T) {
	// Adapt
	m := createBasicMem()
	m.V[3] = 1
	m.DelayTimer = 0x42
	// Act
	Dispatch(m, 0xF307)

	// Assert
	assert.Equal(t, m.DelayTimer, m.V[3], "Set VX to delay timer")
	assert.Equal(t, uint16(0x202), m.PC, "Move to the next instruction")
}

func TestFX15(t *testing.T) {
	// Adapt
	m := createBasicMem()
	m.V[0x7] = 0xFF

	// Act
	Dispatch(m, 0xF715)

	// Assert
	assert.Equal(t, uint16(0xFF), m.DelayTimer, "Delay timer sets")
	assert.Equal(t, uint16(0x202), m.PC, "Move to the next instruction")
}

func TestFX18(t *testing.T) {
	// Adapt
	m := createBasicMem()
	m.V[0x7] = 0xF0

	// Act
	Dispatch(m, 0xF718)

	// Assert
	assert.Equal(t, uint16(0xF0), m.SoundTimer, "Sound timer sets")
	assert.Equal(t, uint16(0x202), m.PC, "Move to the next instruction")
}

func TestFX1E(t *testing.T) {
	// Adapt
	m := createBasicMem()
	m.V[0x7] = 0x30
	m.I = 0x3

	// Act
	Dispatch(m, 0xF71E)

	// Assert
	assert.Equal(t, uint16(0x33), m.I, "I operator sets")
	assert.Equal(t, uint16(0x202), m.PC, "Move to the next instruction")
}

func TestXyExtractor(t *testing.T) {
	// Adapt
	opcode := uint16(0x1234)

	// Act
	x, y := xyExtractor(opcode)

	// Assert
	assert.Equal(t, 2, x, "X value")
	assert.Equal(t, 3, y, "Y value")
}

func TestXyExtractor_0(t *testing.T) {
	// Adapt
	opcode := uint16(0x0034)

	// Act
	x, y := xyExtractor(opcode)

	// Assert
	assert.Equal(t, 0, x, "X value")
	assert.Equal(t, 3, y, "Y value")
}
