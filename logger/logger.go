package logger

import (
	"github.com/cihub/seelog"
)

const (
	PATH = "/var/lib/poi/logs/config/seelog.xml"
)

func Initialize() error {
	var err error

	logger, err := seelog.LoggerFromConfigAsFile("/var/lib/poi/logs/config/seelog.xml")
	if err != nil {
		return err
	}
	seelog.ReplaceLogger(logger)

	return nil
}

func Trace(v ...interface{}) {
	seelog.Trace(v)
}

func Debug(v ...interface{}) {
	seelog.Debug(v)
}

func Info(v ...interface{}) {
	seelog.Info(v)
}

func Warn(v ...interface{}) error {
	return seelog.Warn(v)
}

func Error(v ...interface{}) error {
	return seelog.Error(v)
}

func Critical(v ...interface{}) error {
	return seelog.Critical(v)
}
