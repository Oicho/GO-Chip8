package myLogger

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
)

var (
	// Trace used for memory dump
	Trace *log.Logger
	// Info used for standard log
	Info *log.Logger
	// Warning used for small error output
	Warning *log.Logger
	// Error used for programm breaking error output
	Error *log.Logger
)

type mlog log.Logger

var verbose = false

// InfoVerbosePrint print on Info output
// only if myLogger is initialize with true
func InfoVerbosePrint(s string) {
	if verbose {
		Info.Println(s)
	}
}

// TraceVerbosePrint print on Trace output
// only if myLogger is initialize with true
func TraceVerbosePrint(s string) {
	if verbose {
		Trace.Println(s)
	}
}

// WarningVerbosePrint print on Warning output
// only if myLogger is initialize with true
func WarningVerbosePrint(s string) {
	if verbose {
		Warning.Println(s)
	}
}

// ErrorVerbosePrint print on Error output
// only if myLogger is initialize with true
func ErrorVerbosePrint(s string) {
	if verbose {
		Error.Println(s)
	}
}

// InfoPrint print on Info output
func InfoPrint(s string) {
	Info.Println(s)
}

// TracePrint print on Trace output
func TracePrint(s string) {
	Trace.Println(s)
}

// WarningPrint print on Warning output
func WarningPrint(s string) {
	Warning.Println(s)
}

// ErrorPrint print on Error output
func ErrorPrint(s string) {
	Error.Println(s)
}

// Uint16ToString convert a uint16 to its hexadecimal representation
func Uint16ToString(i uint16) string {
	return hex.EncodeToString([]byte{byte(i >> 8), byte(i & 0x00FF)})
}

// ByteToString convert a uint16 to its hexadecimal representation
func ByteToString(i byte) string {
	return hex.EncodeToString([]byte{i})
}

// Init Initialize the logger output and set the verbose flag
func Init(b bool) error {
	verbose = b
	var logpath = os.Getenv("GOPATH") + "/src/github.com/Oicho/GO-Chip8/log/GO-Chip8.out"
	f, err := os.OpenFile(logpath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("log Failed")
		return err
	}
	Trace = log.New(f,
		"TRACE: ",
		log.Ltime|log.Lshortfile)

	Info = log.New(f,
		"INFO: ",
		log.Ltime|log.Lshortfile)

	Warning = log.New(f,
		"WARNING: ",
		log.Ltime|log.Lshortfile)

	Error = log.New(f,
		"ERROR: ",
		log.Ltime|log.Lshortfile)
	return nil
}
