package chip8

import (
	"github.com/Oicho/GO-Chip8/myLogger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type MemoryTestSuite struct {
	suite.Suite
}

func (suite *MemoryTestSuite) SetupTest() {
	myLogger.Init(true)
}

func createBasicMem() *Memory {
	var mem = &Memory{}
	mem.Init()
	return mem
}

func (suite *MemoryTestSuite) TestInit() {
	// Adapt
	var mem = &Memory{}

	// Act
	mem.Init()

	// Assert
	for i := 0; i < 80; i++ {
		assert.Equal(suite.T(), chip8Fontset[i], mem.Memory[i], "Fontset is loaded")
	}
	for i := range mem.Screen {
		for j := range mem.Screen[i] {
			assert.False(suite.T(), mem.Screen[i][j], "Screen not clean")
		}
	}
	assert.Equal(suite.T(), uint16(0x200), mem.PC, "Pc not set")
	assert.Equal(suite.T(), len(mem.Screen), 64, "Screen width")
	assert.Equal(suite.T(), len(mem.Screen[0]), 32, "Screen height")
}

func (suite *MemoryTestSuite) TestLoadRom_BadPath() {
	// Adapt
	m := createBasicMem()

	// Act
	err := m.LoadRom("not_found")

	// Assert
	assert.NotNil(suite.T(), err, "Error raised")
}

func (suite *MemoryTestSuite) TestLoadRom_Goodfile() {
	// Adapt
	m := createBasicMem()

	// Act
	err := m.LoadRom("../rom/TICTAC")

	// Assert
	assert.Nil(suite.T(), err, "Error raised")
	file, _ := os.Open("../rom/TICTAC")
	defer file.Close()
	data := make([]byte, 0xE00)
	nbBytes, _ := file.Read(data)
	i := 0
	for ; i < nbBytes; i++ {
		assert.Equal(suite.T(), data[i], m.Memory[i+0x200], "Should load Rom data")
	}
	for ; i < 0xE00; i++ {
		assert.Equal(suite.T(), 0, m.Memory[i+0x200], "Should be null")
	}

}

func (suite *MemoryTestSuite) TestFetch() {
	m := createBasicMem()
	m.Memory[m.PC] = 0x01
	m.Memory[m.PC+1] = 0x23
	assert.Equal(suite.T(), uint16(0x0123), m.Fetch(), "Simple opcode fetching")
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestMemoryTestSuite(t *testing.T) {
	suite.Run(t, new(MemoryTestSuite))
}
