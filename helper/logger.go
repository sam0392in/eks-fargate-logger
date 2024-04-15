package helper

import (
	"github.com/sirupsen/logrus"
)

func Logger() *logrus.Logger {
	var Logger *logrus.Logger
	Logger = logrus.New()
	Logger.Formatter = &logrus.JSONFormatter{}
	Logger.Level = logrus.InfoLevel
	return Logger
}
