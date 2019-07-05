package util

import (
	"fmt"

	"github.com/littletwolee/commons/file"
	"github.com/littletwolee/commons/logger"
)

var _logger logger.Ilogger

func Logger() logger.Ilogger {
	return _logger
}
func LoggerInit() {
	path := Config().App.Log.Path
	if err := file.GetFile().PathExists(path, true); err != nil {
		panic(fmt.Errorf("get log path error: %s", err.Error()))
	}
	_logger = logger.Log(path)
}
