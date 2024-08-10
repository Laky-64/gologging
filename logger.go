package gologging

func Debug(message ...any) {
	internalLog(DebugLevel, message...)
}

func Info(message ...any) {
	internalLog(InfoLevel, message...)
}

func Warn(message ...any) {
	internalLog(WarnLevel, message...)
}

func Error(message ...any) {
	internalLog(ErrorLevel, message...)
}

func Fatal(message ...any) {
	internalLog(FatalLevel, message...)
}
