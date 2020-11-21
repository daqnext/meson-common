package logger

import (
	"github.com/sirupsen/logrus"
)

var BaseLogger *logrus.Logger

func Debug(msg string, params ...interface{}) {
	if BaseLogger == nil {
		return
	}
	BaseLogger.WithFields(SliceToFields(params)).Debug(msg)
}

func Info(msg string, params ...interface{}) {
	if BaseLogger == nil {
		return
	}
	BaseLogger.WithFields(SliceToFields(params)).Info(msg)
}

func Warn(msg string, params ...interface{}) {
	if BaseLogger == nil {
		return
	}
	BaseLogger.WithFields(SliceToFields(params)).Warn(msg)
}

func Error(msg string, params ...interface{}) {
	if BaseLogger == nil {
		return
	}
	BaseLogger.WithFields(SliceToFields(params)).Error(msg)
}

func Fatal(msg string, params ...interface{}) {
	if BaseLogger == nil {
		return
	}
	BaseLogger.WithFields(SliceToFields(params)).Fatal(msg)
}
