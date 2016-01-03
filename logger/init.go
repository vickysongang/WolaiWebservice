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
