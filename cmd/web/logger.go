package main

import (
	"log"
	"os"
)

const (
	InfoLevel    = "INFO"
	WarningLevel = "WARNING"
	ErrorLevel   = "ERROR"
)

var (
	infoLog    = log.New(os.Stdout, InfoLevel+"\t", log.Ldate|log.Ltime)
	warningLog = log.New(os.Stdout, WarningLevel+"\t", log.Ldate|log.Ltime|log.Llongfile)
	errorLog   = log.New(os.Stderr, ErrorLevel+"\t", log.Ldate|log.Ltime|log.Lshortfile)
)

func Info(input ...any) {
	infoLog.Print(input...)
}

func InfoLn(input ...any) {
	infoLog.Println(input...)
}

func InfoF(format string, input ...any) {
	infoLog.Printf(format, input...)
}

func Warn(input ...any) {
	warningLog.Print(input...)
}

func WarnLn(input ...any) {
	warningLog.Println(input...)
}

func WarnF(format string, input ...any) {
	warningLog.Printf(format, input...)
}

func Error(input ...any) {
	errorLog.Print(input...)
}

func ErrorLn(input ...any) {
	errorLog.Println(input...)
}

func ErrorF(format string, input ...any) {
	errorLog.Printf(format, input...)
}

