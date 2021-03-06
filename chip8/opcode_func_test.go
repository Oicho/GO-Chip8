package chip8

import (
	"strconv"
	"testing"

	"github.com/Oicho/GO-Chip8/myLogger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type OpcodeTestSuite struct {
	suite.Suite
}

var MockCheckInputs = CheckInputs

func (suite *OpcodeTestSuite) SetupTest() {
	myLogger.Init(true)
	CheckInputs = MockCheckInputs
}

func (suite *OpcodeTestSuite) TestClearScreen() {
	// Adapt
	m := createBasicMem()
	m.Screen[0][0] = true
	m.Screen[3][5] = true

	// Act
	m.Decode(0x00E0)

	// Assert
	assert.False(suite.T(), m.Screen[0][0], "Reset screen")
	assert.False(suite.T(), m.Screen[3][5], "Reset screen")
	assert.Equal(suite.T(), uint16(0x202), m.PC, "Go to the next instruction")
}

func (suite *OpcodeTestSuite) TestReturnFromSubRoutine() {
	// Adapt
	m := createBasicMem()
	m.SP = 1
	m.CallStack[0] = 0x321

	// Act
	m.Decode(0x00EE)

	// Assert
	assert.Equal(suite.T(), uint16(0x321), m.PC, "Return from subroutine")
	assert.Equal(suite.T(), 0, m.SP, "Decrement stack pointer")
}

func (suite *OpcodeTestSuite) TestRca() {
	// Adapt
	m := createBasicMem()

	// Act
	m.Decode(0x0442)

	// Assert
	assert.Equal(suite.T(), uint16(0x202), m.PC, "Go to the next instruction")
}

func (suite *OpcodeTestSuite) Test1NNN_OK() {
	// Adapt
	m := createBasicMem()

	// Act
	m.Decode(0x1242)

	// Assert
	assert.Equal(suite.T(), uint16(0x242), m.PC, "PC didn't move")
}

func (suite *OpcodeTestSuite) Test2NNN_OK() {
	// Adapt
	m := createBasicMem()

	// Act
	m.Decode(0x2FFF)

	// Assert
	assert.Equal(suite.T(), 1, m.SP, "Increment stack pointer")
	assert.Equal(suite.T(), uint16(0x200), m.CallStack[0], "Return position")
	assert.Equal(suite.T(), uint16(0xFFF), m.PC, "PC didn't move")
}

func (suite *OpcodeTestSuite) Test3XNN_Skip() {
	// Adapt
	m := createBasicMem()
	m.V[5] = 0x42

	// Act
	m.Decode(0x3542)

	// Assert
	assert.Equal(suite.T(), uint16(0x204), m.PC, "Haven't skip the instruction")
}

func (suite *OpcodeTestSuite) Test3XNN_NoSkip() {
	// Adapt
	m := createBasicMem()
	m.V[5] = 0x42

	// Act
	m.Decode(0x35FF)

	// Assert
	assert.Equal(suite.T(), uint16(0x202), m.PC, "Haven't skip the instruction")
}

func (suite *OpcodeTestSuite) Test4XNN_NoSkip() {
	// Adapt
	m := createBasicMem()
	m.V[5] = 0x42

	// Act
	m.Decode(0x45FF)

	// Assert
	assert.Equal(suite.T(), uint16(0x204), m.PC, "Have skip the instruction")
}

func (suite *OpcodeTestSuite) Test4XNN_Skip() {
	// Adapt
	m := createBasicMem()
	m.V[5] = 0x42

	// Act
	m.Decode(0x4542)

	// Assert
	assert.Equal(suite.T(), uint16(0x202), m.PC, "Haven't skip the instruction")
}

func (suite *OpcodeTestSuite) Test5XY0_Skip() {
	// Adapt
	m := createBasicMem()
	m.V[2] = 0x23
	m.V[3] = 0x23

	// Act
	m.Decode(0x5230)

	// Assert
	assert.Equal(suite.T(), uint16(0x204), m.PC, "Haven't skip the instruction")
}

func (suite *OpcodeTestSuite) Test5XY0_NoSkip() {
	// Adapt
	m := createBasicMem()
	m.V[2] = 0x23
	m.V[3] = 0x42

	// Act
	m.Decode(0x5230)

	// Assert
	assert.Equal(suite.T(), uint16(0x202), m.PC, "Skip the instruction")
}

func (suite *OpcodeTestSuite) Test6XNN_SimpleSet() {
	// Adapt
	m := createBasicMem()
	m.V[2] = 0x23

	// Act
	m.Decode(0x6242)

	// Assert
	assert.Equal(suite.T(), uint16(0x202), m.PC, "Move to the next instruction")
	assert.Equal(suite.T(), byte(0x42), m.V[2], "Set Register")
}

func (suite *OpcodeTestSuite) Test6XNN_SameValueSet() {
	// Adapt
	m := createBasicMem()
	m.V[2] = 0x20

	// Act
	m.Decode(0x6220)

	// Assert
	assert.Equal(suite.T(), uint16(0x202), m.PC, "Move to the next instruction")
	assert.Equal(suite.T(), byte(0x20), m.V[2], "Set Register")
}

func (suite *OpcodeTestSuite) Test7XNN() {
	// Adapt
	m := createBasicMem()
	m.V[2] = 0x00

	// Act
	m.Decode(0x7230)

	// Assert
	assert.Equal(suite.T(), uint16(0x202), m.PC, "Move to the next instruction")
	assert.Equal(suite.T(), byte(0x30), m.V[2], "Add to Register")
}

func (suite *OpcodeTestSuite) Test7XNN_Overflow() {
	// Adapt
	m := createBasicMem()
	m.V[2] = 0xFF

	// Act
	m.Decode(0x7210)

	// Assert
	assert.Equal(suite.T(), uint16(0x202), m.PC, "Move to the next instruction")
	assert.Equal(suite.T(), byte(0xF), m.V[2], "Add to Register with overflow")
}

// TODO remove me
func (suite *OpcodeTestSuite) Test8() {
	// Adapt
	m := createBasicMem()

	// Act
	m.Decode(0x8)

	// Assert
	assert.Equal(suite.T(), uint16(0x202), m.PC, "Move to the next instruction")
}

func (suite *OpcodeTestSuite) Test8XY0() {
	// Adapt
	m := createBasicMem()
	m.V[0x3] = 0x32
	m.V[0xF] = 0xA0

	// Act
	m.Decode(0x83F0)

	// Assert
	assert.Equal(suite.T(), uint16(0xA0), m.V[0x3], "Changed VX")
	assert.Equal(suite.T(), uint16(0xA0), m.V[0xF], "Unchanged VY")
	assert.Equal(suite.T(), uint16(0x202), m.PC, "Move to the next instruction")
}

func (suite *OpcodeTestSuite) Test8XY1_FF() {
	// Adapt
	m := createBasicMem()
	m.V[0x4] = 0xF0
	m.V[0x5] = 0x0F

	// Act
	m.Decode(0x8451)

	// Assert
	assert.Equal(suite.T(), uint16(0xFF), m.V[0x4], "Changed VX")
	assert.Equal(suite.T(), uint16(0x0F), m.V[0x5], "Unchanged VY")
	assert.Equal(suite.T(), uint16(0x202), m.PC, "Move to the next instruction")
}

func (suite *OpcodeTestSuite) Test8XY1_0F() {
	// Adapt
	m := createBasicMem()
	m.V[0x4] = 0x00
	m.V[0x5] = 0x0F

	// Act
	m.Decode(0x8451)

	// Assert
	assert.Equal(suite.T(), uint16(0x0F), m.V[0x4], "Changed VX")
	assert.Equal(suite.T(), uint16(0x0F), m.V[0x5], "Unchanged VY")
	assert.Equal(suite.T(), uint16(0x202), m.PC, "Move to the next instruction")
}

func (suite *OpcodeTestSuite) Test8XY2_0() {
	// Adapt
	m := createBasicMem()
	m.V[0x2] = 0x0F
	m.V[0xF] = 0xF0

	// Act
	m.Decode(0x82F2)

	// Assert
	assert.Equal(suite.T(), uint16(0x00), m.V[0x2], "Changed VX")
	assert.Equal(suite.T(), uint16(0xF0), m.V[0xF], "Unchanged VY")
	assert.Equal(suite.T(), uint16(0x202), m.PC, "Move to the next instruction")
}

func (suite *OpcodeTestSuite) Test8XY2_0F() {
	// Adapt
	m := createBasicMem()
	m.V[0x2] = 0x0F
	m.V[0xF] = 0xFF

	// Act
	m.Decode(0x82F2)

	// Assert
	assert.Equal(suite.T(), uint16(0x0F), m.V[0x2], "Changed VX")
	assert.Equal(suite.T(), uint16(0xFF), m.V[0xF], "Unchanged VY")
	assert.Equal(suite.T(), uint16(0x202), m.PC, "Move to the next instruction")
}

func (suite *OpcodeTestSuite) Test8XY3_00() {
	// Adapt
	m := createBasicMem()
	m.V[0x2] = 0xFF
	m.V[0x3] = 0xFF

	// Act
	m.Decode(0x8233)

	// Assert
	assert.Equal(suite.T(), uint16(0x00), m.V[0x2], "Changed VX")
	assert.Equal(suite.T(), uint16(0xFF), m.V[0x3], "Unchanged VY")
	assert.Equal(suite.T(), uint16(0x202), m.PC, "Move to the next instruction")
}

func (suite *OpcodeTestSuite) Test8XY4_Simple_add() {
	// Adapt
	m := createBasicMem()
	m.V[2] = 0x3
	m.V[3] = 0x2
	// Act
	m.Decode(0x8234)

	// Assert
	assert.Equal(suite.T(), uint16(5), m.V[0x2], "Changed VX")
	assert.Equal(suite.T(), uint16(2), m.V[0x3], "Unchanged VY")
	assert.Equal(suite.T(), uint16(0), m.V[0xF], "No carry flag")
	assert.Equal(suite.T(), uint16(0x202), m.PC, "Move to the next instruction")
}

func (suite *OpcodeTestSuite) Test8XY4_Overflow() {
	// Adapt
	m := createBasicMem()
	m.V[2] = 0xFF
	m.V[3] = 0x2
	// Act
	m.Decode(0x8234)

	// Assert
	assert.Equal(suite.T(), uint16(1), m.V[0x2], "Changed VX")
	assert.Equal(suite.T(), uint16(2), m.V[0x3], "Unchanged VY")
	assert.Equal(suite.T(), uint16(1), m.V[0xF], "Carry flag")
	assert.Equal(suite.T(), uint16(0x202), m.PC, "Move to the next instruction")
}

func (suite *OpcodeTestSuite) Test8XY5_Simple_sub() {
	// Adapt
	m := createBasicMem()
	m.V[2] = 0x5
	m.V[3] = 0x2
	// Act
	m.Decode(0x8235)

	// Assert
	assert.Equal(suite.T(), uint16(3), m.V[0x2], "Changed VX")
	assert.Equal(suite.T(), uint16(2), m.V[0x3], "Unchanged VY")
	assert.Equal(suite.T(), uint16(0), m.V[0xF], "No carry flag")
	assert.Equal(suite.T(), uint16(0x202), m.PC, "Move to the next instruction")
}

func (suite *OpcodeTestSuite) Test8XY5_Overflow() {
	// Adapt
	m := createBasicMem()
	m.V[2] = 0x5
	m.V[3] = 0x7
	// Act
	m.Decode(0x8235)

	// Assert
	assert.Equal(suite.T(), uint16(0xFE), m.V[0x2], "Changed VX")
	assert.Equal(suite.T(), uint16(7), m.V[0x3], "Unchanged VY")
	assert.Equal(suite.T(), uint16(1), m.V[0xF], "No carry flag")
	assert.Equal(suite.T(), uint16(0x202), m.PC, "Move to the next instruction")
}

func (suite *OpcodeTestSuite) Test8XY6_no_flag() {
	// Adapt
	m := createBasicMem()
	m.V[0xE] = 2

	// Act
	m.Decode(0x8EE6)

	// Assert
	assert.Equal(suite.T(), uint16(1), m.V[0xE], "Unchanged VY")
	assert.Equal(suite.T(), uint16(0), m.V[0xF], "No carry flag")
	assert.Equal(suite.T(), uint16(0x202), m.PC, "Move to the next instruction")
}

func (suite *OpcodeTestSuite) Test8XY6_flag() {
	// Adapt
	m := createBasicMem()
	m.V[0xE] = 9

	// Act
	m.Decode(0x8EE6)

	// Assert
	assert.Equal(suite.T(), uint16(4), m.V[0xE], "Unchanged VY")
	assert.Equal(suite.T(), uint16(1), m.V[0xF], "No carry flag")
	assert.Equal(suite.T(), uint16(0x202), m.PC, "Move to the next instruction")
}

func (suite *OpcodeTestSuite) Test8XY6_different_XY() {
	// Adapt
	m := createBasicMem()
	m.V[0xE] = 9

	// Act
	m.Decode(0x8EA6)

	// Assert
	assert.Equal(suite.T(), uint16(4), m.V[0xE], "Unchanged VY")
	assert.Equal(suite.T(), uint16(1), m.V[0xF], "No carry flag")
	assert.Equal(suite.T(), uint16(0x202), m.PC, "Move to the next instruction")
}

func (suite *OpcodeTestSuite) Test8XY6_zero() {
	// Adapt
	m := createBasicMem()
	m.V[0xE] = 1

	// Act
	m.Decode(0x8EE6)

	// Assert
	assert.Equal(suite.T(), uint16(0), m.V[0xE], "Unchanged VY")
	assert.Equal(suite.T(), uint16(1), m.V[0xF], "No carry flag")
	assert.Equal(suite.T(), uint16(0x202), m.PC, "Move to the next instruction")
}

func (suite *OpcodeTestSuite) Test8XY7_Simple_sub() {
	// Adapt
	m := createBasicMem()
	m.V[2] = 5
	m.V[3] = 10
	// Act
	m.Decode(0x8237)

	// Assert
	assert.Equal(suite.T(), uint16(5), m.V[0x2], "Changed VX")
	assert.Equal(suite.T(), uint16(10), m.V[0x3], "Unchanged VY")
	assert.Equal(suite.T(), uint16(0), m.V[0xF], "No carry flag")
	assert.Equal(suite.T(), uint16(0x202), m.PC, "Move to the next instruction")
}

func (suite *OpcodeTestSuite) Test8XY7_Overflow() {
	// Adapt
	m := createBasicMem()
	m.V[2] = 10
	m.V[3] = 4
	// Act
	m.Decode(0x8237)

	// Assert
	assert.Equal(suite.T(), uint16(0xFA), m.V[0x2], "Changed VX")
	assert.Equal(suite.T(), uint16(4), m.V[0x3], "Unchanged VY")
	assert.Equal(suite.T(), uint16(1), m.V[0xF], "No carry flag")
	assert.Equal(suite.T(), uint16(0x202), m.PC, "Move to the next instruction")
}

func (suite *OpcodeTestSuite) Test8XYE_no_flag() {
	// Adapt
	m := createBasicMem()
	m.V[0xA] = 3

	// Act
	m.Decode(0x8AAE)

	// Assert
	assert.Equal(suite.T(), uint16(6), m.V[0xA], "Changed VX")
	assert.Equal(suite.T(), uint16(0), m.V[0xF], "No carry flag")
	assert.Equal(suite.T(), uint16(0x202), m.PC, "Move to the next instruction")
}
func (suite *OpcodeTestSuite) Test8XYE_different_XY() {
	// Adapt
	m := createBasicMem()
	m.V[0xA] = 3

	// Act
	m.Decode(0x8ACE)

	// Assert
	assert.Equal(suite.T(), uint16(6), m.V[0xA], "Changed VX")
	assert.Equal(suite.T(), uint16(0), m.V[0xF], "No carry flag")
	assert.Equal(suite.T(), uint16(0x202), m.PC, "Move to the next instruction")
}

func (suite *OpcodeTestSuite) Test8XYE_flag() {
	// Adapt
	m := createBasicMem()
	m.V[0xA] = 0xC0

	// Act
	m.Decode(0x8AAE)

	// Assert
	assert.Equal(suite.T(), uint16(0x80), m.V[0xA], "Changed VX")
	assert.Equal(suite.T(), uint16(1), m.V[0xF], "No carry flag")
	assert.Equal(suite.T(), uint16(0x202), m.PC, "Move to the next instruction")
}

func (suite *OpcodeTestSuite) Test8XYE_zero() {
	// Adapt
	m := createBasicMem()
	m.V[0xA] = 0

	// Act
	m.Decode(0x8AAE)

	// Assert
	assert.Equal(suite.T(), uint16(0x0), m.V[0xA], "Changed VX")
	assert.Equal(suite.T(), uint16(0), m.V[0xF], "No carry flag")
	assert.Equal(suite.T(), uint16(0x202), m.PC, "Move to the next instruction")
}

func (suite *OpcodeTestSuite) Test9XY0_NoSkip() {
	// Adapt
	m := createBasicMem()
	m.V[2] = 0x23
	m.V[3] = 0x23

	// Act
	m.Decode(0x9230)

	// Assert
	assert.Equal(suite.T(), uint16(0x202), m.PC, "Haven't skip the instruction")
}

func (suite *OpcodeTestSuite) Test9XY0_Skip() {
	// Adapt
	m := createBasicMem()
	m.V[2] = 0x23
	m.V[3] = 0x42

	// Act
	m.Decode(0x9230)

	// Assert
	assert.Equal(suite.T(), uint16(0x204), m.PC, "Skip the instruction")
}

func (suite *OpcodeTestSuite) TestANNN() {
	// Adapt
	m := createBasicMem()
	m.I = 42

	// Act
	m.Decode(0xA153)

	// Assert
	assert.Equal(suite.T(), uint16(0x153), m.I, "Set I")
}

func (suite *OpcodeTestSuite) TestBNNN_with0() {
	// Adapt
	m := createBasicMem()
	m.V[0] = 0x30

	// Act
	m.Decode(0xB000)

	// Assert
	assert.Equal(suite.T(), uint16(0x30), m.PC, "Move PC")
}

func (suite *OpcodeTestSuite) TestBNNN_withV0_0() {
	// Adapt
	m := createBasicMem()
	m.V[0] = 0x0

	// Act
	m.Decode(0xB100)

	// Assert
	assert.Equal(suite.T(), uint16(0x100), m.PC, "Move PC")
}

func (suite *OpcodeTestSuite) TestBNNN_Overflow() {
	// Adapt
	m := createBasicMem()
	m.V[0] = 0x4

	// Act
	m.Decode(0xBFFF)

	// Assert
	assert.Equal(suite.T(), uint16(0x3), m.PC, "Move PC")
}

func (suite *OpcodeTestSuite) TestCXNN() {
	// Adapt
	m := createBasicMem()
	m.V[0xF] = 0xFF

	// Act
	m.Decode(0xCF00)

	// Assert
	assert.Equal(suite.T(), uint16(0x202), m.PC, "Move PC")
	assert.Equal(suite.T(), 0, m.V[0xF], "And operator")
}

func (suite *OpcodeTestSuite) TestDXYNN_noHeight() {
	// Adapt
	m := createBasicMem()

	// Act
	m.Decode(0xD230)

	//	Assert
	for x, superArr := range m.Screen {
		for y, b := range superArr {
			assert.False(suite.T(), b, "Pixel at x:"+strconv.Itoa(x)+", y:"+strconv.Itoa(y)+" modified")
		}
	}
	assert.Equal(suite.T(), uint16(0x202), m.PC, "Move to the next instruction")
}

func (suite *OpcodeTestSuite) TestDXYNN_emptySprites() {
	// Adapt
	m := createBasicMem()
	// empty memory
	m.I = 0x300

	// Act
	m.Decode(0xD23F)

	//	Assert
	for x, superArr := range m.Screen {
		for y, b := range superArr {
			assert.False(suite.T(), b, "Pixel at x:"+strconv.Itoa(x)+", y:"+strconv.Itoa(y)+" modified")
		}
	}
	assert.Equal(suite.T(), uint16(0x202), m.PC, "Move to the next instruction")
}

func (suite *OpcodeTestSuite) TestDXYNN_font() {
	// Adapt
	m := createBasicMem()
	// empty memory
	m.I = 0x0

	// Act
	m.Decode(0xD23F)

	//	Assert
	res := false
	for _, superArr := range m.Screen {
		for _, b := range superArr {
			res = b || res
		}
	}
	assert.True(suite.T(), res, "At least one print")
	assert.Equal(suite.T(), uint16(0x202), m.PC, "Move to the next instruction")
}

func (suite *OpcodeTestSuite) TestDXYNN_flag() {
	// Adapt
	m := createBasicMem()
	// empty memory
	m.I = 0x0
	for x, superArr := range m.Screen {
		for y, _:= range superArr {

			m.Screen[x][y] = true
		}
	}


	// Act
	m.Decode(0xD23F)

	//	Assert
	assert.Equal(suite.T(), uint16(1), m.V[0xF], "Flag raised")
	assert.Equal(suite.T(), uint16(0x202), m.PC, "Move to the next instruction")
}

func (suite *OpcodeTestSuite) TestENNN() {
	// Adapt
	m := createBasicMem()

	// Act
	m.Decode(0xEFFF)

	// Assert
	assert.Equal(suite.T(), uint16(0x202), m.PC, "Move to the next instruction")
}

func (suite *OpcodeTestSuite) TestEX9E_Good_Key_pressed_then_Skip() {
	// Adapt
	m := createBasicMem()
	CheckInputs = func (m *Memory) (bool, byte){
		return true, 4
	}

	// Act
	m.V[0xB] = 4
	m.Decode(0xEB9E)

	// Assert
	assert.Equal(suite.T(), uint16(0x204), m.PC, "Move to the next instruction")
}

func (suite *OpcodeTestSuite) TestEX9E_No_Key_Pressed_then_no_skip() {
	// Adapt
	m := createBasicMem()
	CheckInputs = func (m *Memory) (bool, byte){
		return false, 4
	}

	// Act
	m.V[0xB] = 4
	m.Decode(0xEB9E)

	// Assert
	assert.Equal(suite.T(), uint16(0x202), m.PC, "Go to the next instruction")
}
func (suite *OpcodeTestSuite) TestEX9E_Different_Key_Pressed_then_no_skip() {
	// Adapt
	m := createBasicMem()
	CheckInputs = func (m *Memory) (bool, byte){
		return true, 4
	}

	// Act
	m.V[0xB] = 5
	m.Decode(0xEB9E)

	// Assert
	assert.Equal(suite.T(), uint16(0x202), m.PC, "Go to the next instruction")
}
func (suite *OpcodeTestSuite) TestEXA1_No_Key_Pressed_then_skip() {
	// Adapt
	m := createBasicMem()
	CheckInputs = func (m *Memory) (bool, byte){
		return false, 0
	}

	// Act
	m.V[0xB] = 5
	m.Decode(0xEBA1)

	// Assert
	assert.Equal(suite.T(), uint16(0x204), m.PC, "Skip the next instruction")
}
func (suite *OpcodeTestSuite) TestEXA1_Different_Key_Pressed_then_skip() {
	// Adapt
	m := createBasicMem()
	CheckInputs = func (m *Memory) (bool, byte){
		return true, 1
	}

	// Act
	m.V[0xB] = 5
	m.Decode(0xEBA1)

	// Assert
	assert.Equal(suite.T(), uint16(0x204), m.PC, "Skip the next instruction")
}
func (suite *OpcodeTestSuite) TestEXA1_Good_Key_Pressed_then_no_skip() {
	// Adapt
	m := createBasicMem()
	CheckInputs = func (m *Memory) (bool, byte){
		return true, 1
	}

	// Act
	m.V[0xB] = 5
	m.Decode(0xEBA1)

	// Assert
	assert.Equal(suite.T(), uint16(0x204), m.PC, "Skip the next instruction")
}

func (suite *OpcodeTestSuite) TestFX07() {
	// Adapt
	m := createBasicMem()
	m.V[3] = 1
	m.DelayTimer = 0x42
	// Act
	m.Decode(0xF307)

	// Assert
	assert.Equal(suite.T(), m.DelayTimer, m.V[3], "Set VX to delay timer")
	assert.Equal(suite.T(), uint16(0x202), m.PC, "Move to the next instruction")
}

func (suite *OpcodeTestSuite) TestFX0A() {
	// Adapt
	m := createBasicMem()
	CheckInputs = func (m *Memory) (bool, byte){
		return true, 4
	}
	// Act
	m.Decode(0xF10A)

	// Assert
	assert.Equal(suite.T(), uint16(0x202), m.PC, "Move to the next instruction")
	assert.Equal(suite.T(), 0, m.V[0], "Register not set")
	assert.Equal(suite.T(), 4, m.V[1], "Register set")
	for i := 2; i < 0x10; i++ {
		assert.Equal(suite.T(), 0, m.V[i], "Register not set")
	}
}

func (suite *OpcodeTestSuite) TestFX15() {
	// Adapt
	m := createBasicMem()
	m.V[0x7] = 0xFF

	// Act
	m.Decode(0xF715)

	// Assert
	assert.Equal(suite.T(), uint16(0xFF), m.DelayTimer, "Delay timer sets")
	assert.Equal(suite.T(), uint16(0x202), m.PC, "Move to the next instruction")
}

func (suite *OpcodeTestSuite) TestFX18() {
	// Adapt
	m := createBasicMem()
	m.V[0x7] = 0xF0

	// Act
	m.Decode(0xF718)

	// Assert
	assert.Equal(suite.T(), uint16(0xF0), m.SoundTimer, "Sound timer sets")
	assert.Equal(suite.T(), uint16(0x202), m.PC, "Move to the next instruction")
}

func (suite *OpcodeTestSuite) TestFX1E() {
	// Adapt
	m := createBasicMem()
	m.V[0x7] = 0x30
	m.I = 0x3

	// Act
	m.Decode(0xF71E)

	// Assert
	assert.Equal(suite.T(), uint16(0x33), m.I, "I operator sets")
	assert.Equal(suite.T(), uint16(0x202), m.PC, "Move to the next instruction")
}

func (suite *OpcodeTestSuite) TestFX29() {
	// Adapt
	m := createBasicMem()
	m.V[3] = 2

	// Act
	m.Decode(0xF329)

	// Assert
	assert.Equal(suite.T(), uint16(10), m.I, "I operator sets")
	assert.Equal(suite.T(), uint16(0x202), m.PC, "Move to the next instruction")
}

func (suite *OpcodeTestSuite) TestFX33_123() {
	// Adapt
	m := createBasicMem()
	m.I = 0x300
	m.V[3] = 123

	// Act
	m.Decode(0xF333)

	// Assert
	assert.Equal(suite.T(), uint16(1), m.Memory[0x300], "X?? value")
	assert.Equal(suite.T(), uint16(2), m.Memory[0x301], "?X? value")
	assert.Equal(suite.T(), uint16(3), m.Memory[0x302], "??X value")
	assert.Equal(suite.T(), uint16(0x202), m.PC, "Move to the next instruction")

}

func (suite *OpcodeTestSuite) TestFX33_003() {
	// Adapt
	m := createBasicMem()
	m.I = 0x300
	m.V[3] = 3

	// Act
	m.Decode(0xF333)

	// Assert
	assert.Equal(suite.T(), uint16(0), m.Memory[0x300], "X?? value")
	assert.Equal(suite.T(), uint16(0), m.Memory[0x301], "?X? value")
	assert.Equal(suite.T(), uint16(3), m.Memory[0x302], "??X value")
	assert.Equal(suite.T(), uint16(0x202), m.PC, "Move to the next instruction")
}
func (suite *OpcodeTestSuite) TestFX33_000() {
	// Adapt
	m := createBasicMem()
	m.I = 0x300
	m.V[3] = 0

	// Act
	m.Decode(0xF333)

	// Assert
	assert.Equal(suite.T(), uint16(0), m.Memory[0x300], "X?? value")
	assert.Equal(suite.T(), uint16(0), m.Memory[0x301], "?X? value")
	assert.Equal(suite.T(), uint16(0), m.Memory[0x302], "??X value")
	assert.Equal(suite.T(), uint16(0x202), m.PC, "Move to the next instruction")
}

func (suite *OpcodeTestSuite) TestFF55() {
	// Adapt
	m := createBasicMem()
	m.I = 0x300
	m.V[0x0] = 0x1
	m.V[0x1] = 0x2
	m.V[0x2] = 0x3
	m.V[0x3] = 0x4
	m.V[0x4] = 0x5
	m.V[0x5] = 0x6
	m.V[0x6] = 0x7
	m.V[0x7] = 0x8
	m.V[0x8] = 0x9
	m.V[0x9] = 0xA
	m.V[0xA] = 0xB
	m.V[0xB] = 0xC
	m.V[0xC] = 0xD
	m.V[0xD] = 0xE
	m.V[0xE] = 0xF
	m.V[0xF] = 0x10

	// Act
	m.Decode(0xFF55)

	// Assert
	assert.Equal(suite.T(), uint16(0x1), m.Memory[0x300], "V[] set")
	assert.Equal(suite.T(), uint16(0x2), m.Memory[0x301], "V[] set")
	assert.Equal(suite.T(), uint16(0x3), m.Memory[0x302], "V[] set")
	assert.Equal(suite.T(), uint16(0x4), m.Memory[0x303], "V[] set")
	assert.Equal(suite.T(), uint16(0x5), m.Memory[0x304], "V[] set")
	assert.Equal(suite.T(), uint16(0x6), m.Memory[0x305], "V[] set")
	assert.Equal(suite.T(), uint16(0x7), m.Memory[0x306], "V[] set")
	assert.Equal(suite.T(), uint16(0x8), m.Memory[0x307], "V[] set")
	assert.Equal(suite.T(), uint16(0x9), m.Memory[0x308], "V[] set")
	assert.Equal(suite.T(), uint16(0xA), m.Memory[0x309], "V[] set")
	assert.Equal(suite.T(), uint16(0xB), m.Memory[0x30A], "V[] set")
	assert.Equal(suite.T(), uint16(0xC), m.Memory[0x30B], "V[] set")
	assert.Equal(suite.T(), uint16(0xD), m.Memory[0x30C], "V[] set")
	assert.Equal(suite.T(), uint16(0xE), m.Memory[0x30D], "V[] set")
	assert.Equal(suite.T(), uint16(0xF), m.Memory[0x30E], "V[] set")
	assert.Equal(suite.T(), uint16(0x10), m.Memory[0x30F], "V[] set")

	assert.Equal(suite.T(), uint16(0x202), m.PC, "Move to the next instruction")
}

func (suite *OpcodeTestSuite) TestF055() {
	// Adapt
	m := createBasicMem()
	m.I = 0x300
	m.V[0x0] = 0x1
	m.V[0x1] = 0x2
	m.V[0x2] = 0x3
	m.V[0x3] = 0x4
	m.V[0x4] = 0x5
	m.V[0x5] = 0x6
	m.V[0x6] = 0x7
	m.V[0x7] = 0x8
	m.V[0x8] = 0x9
	m.V[0x9] = 0xA
	m.V[0xA] = 0xB
	m.V[0xB] = 0xC
	m.V[0xC] = 0xD
	m.V[0xD] = 0xE
	m.V[0xE] = 0xF
	m.V[0xF] = 0x10

	// Act
	m.Decode(0xF055)

	// Assert
	assert.Equal(suite.T(), uint16(0x1), m.Memory[0x300], "V[] set")
	assert.Equal(suite.T(), uint16(0x0), m.Memory[0x301], "V[] set")
	assert.Equal(suite.T(), uint16(0x0), m.Memory[0x302], "V[] set")
	assert.Equal(suite.T(), uint16(0x0), m.Memory[0x303], "V[] set")
	assert.Equal(suite.T(), uint16(0x0), m.Memory[0x304], "V[] set")
	assert.Equal(suite.T(), uint16(0x0), m.Memory[0x305], "V[] set")
	assert.Equal(suite.T(), uint16(0x0), m.Memory[0x306], "V[] set")
	assert.Equal(suite.T(), uint16(0x0), m.Memory[0x307], "V[] set")
	assert.Equal(suite.T(), uint16(0x0), m.Memory[0x308], "V[] set")
	assert.Equal(suite.T(), uint16(0x0), m.Memory[0x309], "V[] set")
	assert.Equal(suite.T(), uint16(0x0), m.Memory[0x30A], "V[] set")
	assert.Equal(suite.T(), uint16(0x0), m.Memory[0x30B], "V[] set")
	assert.Equal(suite.T(), uint16(0x0), m.Memory[0x30C], "V[] set")
	assert.Equal(suite.T(), uint16(0x0), m.Memory[0x30D], "V[] set")
	assert.Equal(suite.T(), uint16(0x0), m.Memory[0x30E], "V[] set")
	assert.Equal(suite.T(), uint16(0x0), m.Memory[0x30F], "V[] set")

	assert.Equal(suite.T(), uint16(0x202), m.PC, "Move to the next instruction")
}

func (suite *OpcodeTestSuite) TestFF65_empty_memory() {
	// Adapt
	m := createBasicMem()
	m.I = 0x300

	// Act
	m.Decode(0xFF65)

	// Assert
	for i := 0; i < 0x10; i++{
		assert.Equal(suite.T(), uint16(0x00), m.V[i], "V[0x" + myLogger.ByteToString(byte(i)) + "] not set")
	}

	assert.Equal(suite.T(), uint16(0x202), m.PC, "Move to the next instruction")
}

func (suite *OpcodeTestSuite) TestFF65() {
	// Adapt
	m := createBasicMem()
	m.I = 0x300
	var val byte = 0
	for i := 0x300; i < 0x410; i++ {
		m.Memory[i] = val
		val++
	}

	// Act
	m.Decode(0xFF65)

	// Assert
	val = 0
	for i := 0; i < 0x10; i++{
		assert.Equal(suite.T(), uint16(val), m.V[i], "V[0x" + myLogger.ByteToString(byte(i)) + "] not set")
		val ++
	}

	assert.Equal(suite.T(), uint16(0x202), m.PC, "Move to the next instruction")
}

func (suite *OpcodeTestSuite) TestFB65() {
	// Adapt
	m := createBasicMem()
	m.I = 0x300
	var val byte = 0
	for i := 0x300; i < 0x410; i++ {
		m.Memory[i] = val
		val++
	}

	// Act
	m.Decode(0xFB65)

	// Assert
	val = 0
	for i := 0; i < 0xC; i++{
		assert.Equal(suite.T(), uint16(val), m.V[i], "V[0x" + myLogger.ByteToString(byte(i)) + "] not set")
		val ++
	}

	assert.Equal(suite.T(), uint16(0x202), m.PC, "Move to the next instruction")
}

func (suite *OpcodeTestSuite) TestXyExtractor() {
	// Adapt
	opcode := uint16(0x1234)

	// Act
	x, y := xyExtractor(opcode)

	// Assert
	assert.Equal(suite.T(), 2, x, "X value")
	assert.Equal(suite.T(), 3, y, "Y value")
}

func (suite *OpcodeTestSuite) TestXyExtractor_0() {
	// Adapt
	opcode := uint16(0x0034)

	// Act
	x, y := xyExtractor(opcode)

	// Assert
	assert.Equal(suite.T(), 0, x, "X value")
	assert.Equal(suite.T(), 3, y, "Y value")
}

func TestOpcodeTestSuite(t *testing.T) {
	suite.Run(t, new(OpcodeTestSuite))
}
