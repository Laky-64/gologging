package gologging

var currentLevel = InfoLevel

type Level int32

const (
	// DebugLevel is the debug level.
	DebugLevel Level = -4
	// InfoLevel is the info level.
	InfoLevel Level = 0
	// WarnLevel is the warn level.
	WarnLevel Level = 4
	// ErrorLevel is the error level.
	ErrorLevel Level = 8
	// FatalLevel is the fatal level.
	FatalLevel Level = 12
)

func SetLevel(level Level) {
	currentLevel = level
}
