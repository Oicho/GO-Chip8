package chip8

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	// Adapt
	var mem = &Memory{}

	// Act
	mem.Init()

	// Assert
	for i := 0; i < 80; i++ {
		assert.Equal(t, chip8Fontset[i], mem.Memory[i], "Fontset is loaded")
	}
	for i := range mem.Screen {
		for j := range mem.Screen[i] {
			assert.False(t, mem.Screen[i][j], "Screen not clean")
		}
	}
	assert.Equal(t, uint16(0x200), mem.PC, "Pc not set")
	assert.Equal(t, len(mem.Screen), 64, "Screen width")
	assert.Equal(t, len(mem.Screen[0]), 32, "Screen height")
}

func TestLoadRom_BadPath(t *testing.T) {
	// Adapt
	m := createBasicMem()

	// Act
	err := m.LoadRom("not_found")

	// Assert
	assert.NotNil(t, err, "Error raised")
}

func TestLoadRom_Goodfile(t *testing.T) {
	// Adapt
	m := createBasicMem()

	// Act
	err := m.LoadRom("../rom/TICTAC")

	// Assert
	assert.Nil(t, err, "Error raised")
	file, _ := os.Open("../rom/TICTAC")
	defer file.Close()
	data := make([]byte, 0xE00)
	nbBytes, _ := file.Read(data)
	i := 0
	for ; i < nbBytes; i++ {
		assert.Equal(t, data[i], m.Memory[i+0x200], "Should load Rom data")
	}
	for ; i < 0xE00; i++ {
		assert.Equal(t, 0, m.Memory[i+0x200], "Should be null")
	}

}

func TestFetch(t *testing.T) {
	m := createBasicMem()
	m.Memory[m.PC] = 0x01
	m.Memory[m.PC+1] = 0x23
	assert.Equal(t, uint16(0x0123), m.Fetch(), "Simple opcode fetching")
}
