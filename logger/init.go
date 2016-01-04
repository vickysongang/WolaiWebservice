package logger

import (
	"github.com/cihub/seelog"

	"WolaiWebservice/config"
)

func Initialize() error {
	var err error

	logger, err := seelog.LoggerFromConfigAsFile(config.Env.Seelog.Config)
	if err != nil {
		return err
	}
	seelog.ReplaceLogger(logger)

	return nil
}
