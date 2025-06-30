package gologging

import (
	"bytes"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"io"
	"os"
	"sync"
	"sync/atomic"
)

type Level int32

const (
	// DebugLevel is the debug level.
	DebugLevel Level = 1 << iota
	// InfoLevel is the info level.
	InfoLevel
	// WarnLevel is the warn level.
	WarnLevel
	// ErrorLevel is the error level.
	ErrorLevel
	// FatalLevel is the fatal level.
	FatalLevel
)

const MinTermWidth = 200

type Logger struct {
	level      atomic.Int32
	mu         *sync.RWMutex
	b          bytes.Buffer
	w          io.Writer
	re         *lipgloss.Renderer
	isDiscard  uint32
	timeFormat string
	loggerName string
	s          *Styles
}

func (ctx *Logger) SetOutput(w io.Writer) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	if w == nil {
		w = os.Stderr
	}
	var isDiscard uint32
	if w == io.Discard {
		isDiscard = 1
	}
	atomic.StoreUint32(&ctx.isDiscard, isDiscard)
	if v, ok := registry.Load(w); ok {
		ctx.re = v.(*lipgloss.Renderer)
	} else {
		ctx.re = lipgloss.NewRenderer(w, termenv.WithColorCache(true))
		registry.Store(w, ctx.re)
	}
	ctx.w = w
	ctx.s = defaultStyles(ctx.re)
}

func (ctx *Logger) SetLevel(level Level) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.level.Store(int32(level))
}

func (ctx *Logger) GetLevel() Level {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	return Level(ctx.level.Load())
}

func (ctx *Logger) Debug(message ...any) {
	ctx.internalLog(DebugLevel, message...)
}

func (ctx *Logger) DebugF(message string, args ...any) {
	ctx.Debug(fmt.Sprintf(message, args...))
}

func (ctx *Logger) Info(message ...any) {
	ctx.internalLog(InfoLevel, message...)
}

func (ctx *Logger) InfoF(message string, args ...any) {
	ctx.Info(fmt.Sprintf(message, args...))
}

func (ctx *Logger) Warn(message ...any) {
	ctx.internalLog(WarnLevel, message...)
}

func (ctx *Logger) WarnF(message string, args ...any) {
	ctx.Warn(fmt.Sprintf(message, args...))
}

func (ctx *Logger) Error(message ...any) {
	ctx.internalLog(ErrorLevel, message...)
}

func (ctx *Logger) ErrorF(message string, args ...any) {
	ctx.Error(fmt.Sprintf(message, args...))
}

func (ctx *Logger) Fatal(message ...any) {
	ctx.internalLog(FatalLevel, message...)
}

func (ctx *Logger) FatalF(message string, args ...any) {
	ctx.Fatal(fmt.Sprintf(message, args...))
}
