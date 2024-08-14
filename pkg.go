package gologging

import (
	"bytes"
	"io"
	"os"
	"sync"
	"sync/atomic"
)

var (
	loggersRegistry   = sync.Map{}
	registry          = sync.Map{}
	defaultLogger     atomic.Pointer[Logger]
	defaultLoggerOnce sync.Once
)

func defaultInstance() *Logger {
	dl := defaultLogger.Load()
	if dl == nil {
		defaultLoggerOnce.Do(func() {
			defaultLogger.CompareAndSwap(
				nil, &Logger{
					level:      int32(InfoLevel),
					mu:         &sync.RWMutex{},
					b:          bytes.Buffer{},
					timeFormat: "2006-01-02 15:04:05",
				},
			)
			defaultLogger.Load().SetOutput(os.Stderr)
		})
		dl = defaultLogger.Load()
	}
	return dl
}

func GetLogger(name string) *Logger {
	if logger, ok := loggersRegistry.Load(name); ok {
		return logger.(*Logger)
	}
	logger := defaultInstance()
	logger.loggerName = name
	loggersRegistry.Store(name, logger)
	return logger
}

func SetLevel(level Level) {
	defaultInstance().SetLevel(level)
}

func GetLevel() Level {
	return defaultInstance().GetLevel()
}

func SetOutput(w io.Writer) {
	defaultInstance().SetOutput(w)
}

func Debug(message ...any) {
	defaultInstance().Debug(message...)
}

func Info(message ...any) {
	defaultInstance().Info(message...)
}

func Warn(message ...any) {
	defaultInstance().Warn(message...)
}

func Error(message ...any) {
	defaultInstance().Error(message...)
}

func Fatal(message ...any) {
	defaultInstance().Fatal(message...)
}
