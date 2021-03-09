package log

import (
	"sync"

	_log "log"

	"github.com/sirupsen/logrus"
)

var (
	logger *logrus.Logger
	once   sync.Once
)

//Init initializes the global logger instance
func Init(options ...Option) error {
	logger = logrus.New()
	for _, option := range options {
		if err := option.apply(logger); err != nil {
			return err
		}
	}
	return nil
}

//Logger gets the global logger instance
func Logger() *logrus.Logger {
	once.Do(func() {
		if logger == nil {
			_log.Println("logger 未初始化, 使用默认配置")
			_ = Init()
		}
	})
	return logger
}
