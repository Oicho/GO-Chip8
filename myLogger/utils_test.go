package myLogger

import (
	"testing"
	"github.com/stretchr/testify/assert"

	"os"
)

func TestUint16ToString(t *testing.T) {
	// Adapt
	var i uint16 = 0x1234

	// Act
	s := Uint16ToString(i)

	// Assert
	assert.Equal(t, "1234", s, "fail convert")
}

func TestUint16ToString_Zero(t *testing.T) {
	// Adapt
	var i uint16 = 0

	// Act
	s := Uint16ToString(i)

	// Assert
	assert.Equal(t, "0000", s, "fail convert")
}

func TestByteToString(t *testing.T) {
	// Adapt
	var b byte = 255

	// Act
	s := ByteToString(b)

	// Assert
	assert.Equal(t, "ff", s, "fail convert")
}

func TestByteToString_Zero(t *testing.T) {
	// Adapt
	var b byte = 0

	// Act
	s := ByteToString(b)

	// Assert
	assert.Equal(t, "00", s, "fail convert")
}

func TestInit(t *testing.T) {
	// Act
	err := Init(true)

	// Assert
	assert.Nil(t, err, "Return Nil")
	assert.True(t, verbose, "Verbose mode setup")
	assert.NotNil(t, Trace, "Trace setup")
	assert.NotNil(t, Info, "Info setup")
	assert.NotNil(t, Warning, "Warning setup")
	assert.NotNil(t, Error, "Error setup")
}
func TestInit_Fail(t *testing.T) {
	// Adapt
	logpath = os.Getenv("GOPATH") + "/test/Tasd"

	// Act
	err := Init(true)

	// Assert
	assert.NotNil(t, err, "Error thrown")
}
