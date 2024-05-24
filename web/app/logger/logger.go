package logger

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	Log *logrus.Logger
}

func NewLogger(path string, filename string) *Logger {
	logFile, err := os.OpenFile(path+filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("ERROR opening log file: %v", err)
	}
	log := &logrus.Logger{
		Out:       logFile,
		Formatter: &logrus.TextFormatter{},
		Level:     logrus.DebugLevel,
		ExitFunc:  os.Exit,
	}
	return &Logger{Log: log}
}
