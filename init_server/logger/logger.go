package logger

import (
	"io"
	"log"
	"os"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

const defaultLogFile = "./backend-server.log"

func InitLogger() *logrus.Logger {
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	customFormatter.ForceColors = true

	file, err := os.OpenFile(defaultLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	multi := io.MultiWriter(file, os.Stderr)

	logger = &logrus.Logger{
		Level:     logrus.InfoLevel,
		Formatter: customFormatter,
		Out:       multi,
	}
	return logger
}

func GetLogger() *logrus.Logger {
	if logger == nil {
		InitLogger()
	}
	return logger
}
