package internal

import (
	"log"
	"os"
)

var (
	Logger *log.Logger
)

func InitLogger(logFilePath string, isTestMode bool) {
	var logOutput *os.File
	var err error

	if isTestMode {
		logOutput = os.Stderr
	} else {
		logOutput, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal("Failed to open log file:", err)
		}
	}

	Logger = log.New(logOutput, "", log.Ldate|log.Ltime)
}
