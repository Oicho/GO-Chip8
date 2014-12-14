package myLogger

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"os"
	"strings"
)

type LoggerTestSuite struct {
	suite.Suite
}

var tmpLogPath = os.Getenv("GOPATH") + "/src/github.com/Oicho/GO-Chip8/test/tmp.log"

func (suite *LoggerTestSuite)SetupTest() {
	logpath = tmpLogPath
}

func (suite *LoggerTestSuite)TearDownTest() {
	os.Remove(tmpLogPath)
}

func (suite *LoggerTestSuite)TestUint16ToString() {
	// Adapt
	var i uint16 = 0x1234

	// Act
	s := Uint16ToString(i)

	// Assert
	assert.Equal(suite.T(), "1234", s, "fail convert")
}

func (suite *LoggerTestSuite)TestUint16ToString_Zero() {
	// Adapt
	var i uint16 = 0

	// Act
	s := Uint16ToString(i)

	// Assert
	assert.Equal(suite.T(), "0000", s, "fail convert")
}

func (suite *LoggerTestSuite)TestByteToString() {
	// Adapt
	var b byte = 255

	// Act
	s := ByteToString(b)

	// Assert
	assert.Equal(suite.T(), "ff", s, "fail convert")
}

func (suite *LoggerTestSuite)TestByteToString_Zero() {
	// Adapt
	var b byte = 0

	// Act
	s := ByteToString(b)

	// Assert
	assert.Equal(suite.T(), "00", s, "fail convert")
}

func (suite *LoggerTestSuite)TestInit() {
	// Act
	err := Init(true)

	// Assert
	assert.Nil(suite.T(), err, "Return Nil")
	assert.True(suite.T(), verbose, "Verbose mode setup")
	assert.NotNil(suite.T(), Trace, "Trace setup")
	assert.NotNil(suite.T(), Info, "Info setup")
	assert.NotNil(suite.T(), Warning, "Warning setup")
	assert.NotNil(suite.T(), Error, "Error setup")
}

func (suite *LoggerTestSuite)TestInit_Fail() {
	// Adapt
	logpath = os.Getenv("GOPATH") + "/test/Tasd"

	// Act
	err := Init(true)

	// Assert
	assert.NotNil(suite.T(), err, "Error thrown")
}

func (suite *LoggerTestSuite)TestInfoPrint(){
	// Adapt
	err := Init(false)

	// Act
	InfoPrint("YOLO")

	// Assert
	b, str := printCheck()
	assert.Nil(suite.T(), err, "Return Nil")
	assert.True(suite.T(), b, "Print only one line")
	assert.True(suite.T(), strings.Contains(str, "INFO:"), "Print header")
	assert.True(suite.T(), strings.Contains(str, "YOLO"), "Print Message")
}

func (suite *LoggerTestSuite)TestTracePrint(){
	// Adapt
	err := Init(false)

	// Act
	TracePrint("YOLO")

	// Assert
	b, str := printCheck()
	assert.Nil(suite.T(), err, "Return Nil")
	assert.True(suite.T(), b, "Print only one line")
	assert.True(suite.T(), strings.Contains(str, "TRACE:"), "Print header")
	assert.True(suite.T(), strings.Contains(str, "YOLO"), "Print Message")
}

func (suite *LoggerTestSuite)TestWarningPrint(){
	// Adapt
	err := Init(false)

	// Act
	WarningPrint("asd")

	// Assert
	b, str := printCheck()
	assert.Nil(suite.T(), err, "Return Nil")
	assert.True(suite.T(), b, "Print only one line")
	assert.True(suite.T(), strings.Contains(str, "WARNING:"), "Print header")
	assert.True(suite.T(), strings.Contains(str, "asd"), "Print Message")
}

func (suite *LoggerTestSuite)TestErrorPrint(){
	// Adapt
	err := Init(false)

	// Act
	ErrorPrint("asd")

	// Assert
	b, str := printCheck()
	assert.Nil(suite.T(), err, "Return Nil")
	assert.True(suite.T(), b, "Print only one line")
	assert.True(suite.T(), strings.Contains(str, "ERROR:"), "Print header")
	assert.True(suite.T(), strings.Contains(str, "asd"), "Print Message")
}

func (suite *LoggerTestSuite)TestInfoVerbosePrint(){
	// Adapt
	err := Init(true)

	// Act
	InfoVerbosePrint("YOLO")

	// Assert
	b, str := printCheck()
	assert.Nil(suite.T(), err, "Return Nil")
	assert.True(suite.T(), b, "Print only one line")
	assert.True(suite.T(), strings.Contains(str, "INFO:"), "Print header")
	assert.True(suite.T(), strings.Contains(str, "YOLO"), "Print Message")
}

func (suite *LoggerTestSuite)TestTraceVerbosePrint(){
	// Adapt
	err := Init(true)

	// Act
	TraceVerbosePrint("YOLO")

	// Assert
	b, str := printCheck()
	assert.Nil(suite.T(), err, "Return Nil")
	assert.True(suite.T(), b, "Print only one line")
	assert.True(suite.T(), strings.Contains(str, "TRACE:"), "Print header")
	assert.True(suite.T(), strings.Contains(str, "YOLO"), "Print Message")
}

func (suite *LoggerTestSuite)TestWarningVerbosePrint(){
	// Adapt
	err := Init(true)

	// Act
	WarningVerbosePrint("asd")

	// Assert
	b, str := printCheck()
	assert.Nil(suite.T(), err, "Return Nil")
	assert.True(suite.T(), b, "Print only one line")
	assert.True(suite.T(), strings.Contains(str, "WARNING:"), "Print header")
	assert.True(suite.T(), strings.Contains(str, "asd"), "Print Message")
}

func (suite *LoggerTestSuite)TestErrorVerbosePrint(){
	// Adapt
	err := Init(true)

	// Act
	ErrorVerbosePrint("asd")

	// Assert
	b, str := printCheck()
	assert.Nil(suite.T(), err, "Return Nil")
	assert.True(suite.T(), b, "Print only one line")
	assert.True(suite.T(), strings.Contains(str, "ERROR:"), "Print header")
	assert.True(suite.T(), strings.Contains(str, "asd"), "Print Message")
}

func (suite *LoggerTestSuite)TestInfoVerbosePrint_without_verbose_activated(){
	// Adapt
	err := Init(false)

	// Act
	InfoVerbosePrint("YOLO")

	// Assert
	b, str := printCheck()
	assert.Nil(suite.T(), err, "Return Nil")
	assert.False(suite.T(), b, "Print only one line")
	assert.False(suite.T(), strings.Contains(str, "INFO:"), "Print header")
	assert.False(suite.T(), strings.Contains(str, "YOLO"), "Print Message")
}

func (suite *LoggerTestSuite)TestTraceVerbosePrint_without_verbose_activated(){
	// Adapt
	err := Init(false)

	// Act
	TraceVerbosePrint("YOLO")

	// Assert
	b, str := printCheck()
	assert.Nil(suite.T(), err, "Return Nil")
	assert.False(suite.T(), b, "Print only one line")
	assert.False(suite.T(), strings.Contains(str, "TRACE:"), "Print header")
	assert.False(suite.T(), strings.Contains(str, "YOLO"), "Print Message")
}

func (suite *LoggerTestSuite)TestWarningVerbosePrint_without_verbose_activated(){
	// Adapt
	err := Init(false)

	// Act
	WarningVerbosePrint("asd")

	// Assert
	b, str := printCheck()
	assert.Nil(suite.T(), err, "Return Nil")
	assert.False(suite.T(), b, "Print only one line")
	assert.False(suite.T(), strings.Contains(str, "WARNING:"), "Print header")
	assert.False(suite.T(), strings.Contains(str, "asd"), "Print Message")
}

func (suite *LoggerTestSuite)TestErrorVerbosePrint_without_verbose_activated(){
	// Adapt
	err := Init(false)

	// Act
	ErrorVerbosePrint("asd")

	// Assert
	b, str := printCheck()
	assert.Nil(suite.T(), err, "Return Nil")
	assert.False(suite.T(), b, "Print only one line")
	assert.False(suite.T(), strings.Contains(str, "ERROR:"), "Print header")
	assert.False(suite.T(), strings.Contains(str, "asd"), "Print Message")
}


func TestLoggerTestSuite(t *testing.T) {
	suite.Run(t, new(LoggerTestSuite))
}

func printCheck() (bool , string){
	f, err := os.Open(tmpLogPath)
	if err != nil {
		return false, ""
	}
	defer f.Close()
	arr := make([]byte, 4096)
	length, err :=  f.Read(arr)
	if err != nil {
		return false, ""
	}
	str := string(arr)
	countReturnLine := false
	for i :=0; i < length; i++{
		if (arr[i] == '\n') {
			if countReturnLine {
				return false, str
			}
			countReturnLine = true
		}
	}
	return countReturnLine, str
}
