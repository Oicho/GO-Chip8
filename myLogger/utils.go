package myLogger

import (
	"fmt"
	"log"
	"os"
)

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

type mlog log.Logger

var verbose = false

func InfoVerbosePrint(s string) {
	if verbose {
		Info.Println(s)
	}
}

func TraceVerbosePrint(s string) {
	if verbose {
		Trace.Println(s)
	}
}

func WarningVerbosePrint(s string) {
	if verbose {
		Warning.Println(s)
	}
}
func ErrorVerbosePrint(s string) {
	if verbose {
		Error.Println(s)
	}
}
func InfoPrint(s string) {
	Info.Println(s)
}

func TracePrint(s string) {
	Trace.Println(s)
}

func WarningPrint(s string) {
	Warning.Println(s)
}
func ErrorPrint(s string) {
	Error.Println(s)
}

func Init(b bool) {
	verbose = b
	var logpath = os.Getenv("GOPATH") + "/src/github.com/Oicho/GO-Chip8/log/GO-Chip8.out"
	f, err := os.OpenFile(logpath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("log Failed")
		return
	}
	Trace = log.New(f,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(f,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(f,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(f,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}
