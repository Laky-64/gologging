package gologging

import (
	"errors"
	"fmt"
	"github.com/Laky-64/gologging/types"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"golang.org/x/term"
	"os"
	"strings"
	"sync/atomic"
	"time"
)

func (ctx *Logger) internalLog(level Level, message ...any) {
	if atomic.LoadInt32(&ctx.level) > int32(level) {
		return
	}
	if atomic.LoadUint32(&ctx.isDiscard) != 0 {
		return
	}
	capitalize := func(s string) string {
		return strings.ToUpper(s[:1]) + s[1:]
	}
	var errMessages []string
	for _, x := range message {
		parsed := fmt.Sprintf("%v", x)
		if len(parsed) > 0 {
			errMessages = append(errMessages, strings.ReplaceAll(parsed, "\r", ""))
		}
	}
	if len(errMessages) == 0 {
		return
	}

	iconStyle := ctx.s.IconStyle
	messageStyle := ctx.s.MsgStyle

	switch level {
	case DebugLevel:
		iconStyle = iconStyle.Background(lipgloss.Color("#2f5d77")).
			Foreground(lipgloss.Color("#bbbbbb")).
			SetString("D")
		messageStyle = messageStyle.Foreground(lipgloss.Color("#279999"))
	case InfoLevel:
		iconStyle = iconStyle.Background(lipgloss.Color("#698759")).
			Foreground(lipgloss.Color("#e9f5e6")).
			Bold(true).
			SetString("I")
		messageStyle = messageStyle.Foreground(lipgloss.Color("#abc022"))
	case WarnLevel:
		iconStyle = iconStyle.Background(lipgloss.Color("#bbb527")).
			SetString("W")
		messageStyle = messageStyle.Foreground(lipgloss.Color("#bbb527"))
	case ErrorLevel, FatalLevel:
		iconStyle = iconStyle.Background(lipgloss.Color("#cf5b56")).
			SetString("E")
		messageStyle = messageStyle.Foreground(lipgloss.Color("#cf5b56"))
	}

	termWidth, _, _ := term.GetSize(int(os.Stderr.Fd()))
	tagStyle := ctx.s.TagStyle
	packageStyle := ctx.s.PkgStyle
	if termWidth > MinTermWidth {
		tagStyle = tagStyle.Width(tagStyle.GetWidth() * termWidth / MinTermWidth)
		packageStyle = packageStyle.Width(packageStyle.GetWidth() * termWidth / MinTermWidth)
	}

	messageWidth := termWidth - lipgloss.Width(ctx.s.TimeStyle.Render(ctx.timeFormat))
	messageWidth -= lipgloss.Width(tagStyle.Render())
	messageWidth -= lipgloss.Width(packageStyle.Render())
	messageWidth -= lipgloss.Width(iconStyle.Render())
	messageContainerStyle := lipgloss.NewStyle().
		Width(messageWidth).
		Renderer(ctx.re)

	var mainDetails *types.CallerInfo
	startSkips := 0
	ignoreInfo, _ := getInfo(0)
	for {
		details, err := getInfo(startSkips)
		if err == nil && details.PackageName != ignoreInfo.PackageName {
			mainDetails = details
			break
		} else if errors.Is(err, goRoutineFunction) {
			break
		}
		if startSkips > 10 {
			break
		}
		startSkips++
	}
	var packageName string
	if mainDetails != nil {
		packageName = mainDetails.PackageName
	} else {
		packageName = "unknown"
	}
	matches := tagRgx.FindStringSubmatch(strings.Join(errMessages, " "))
	tagName := matches[2]
	if len(tagName) == 0 {
		packageInfo := strings.Split(packageName, ".")
		tagName = packageInfo[len(packageInfo)-1]
	}
	if tagName == "main" {
		tagStyle = tagStyle.Foreground(lipgloss.Color("#ab91ba"))
	}
	tagName = capitalize(tagName)
	var lines []string
	skips := startSkips - 1
	if level == FatalLevel {
		for {
			skips++
			subDetails, runtimeErr := getInfo(skips)
			if runtimeErr != nil {
				break
			}
			if mainDetails == nil || mainDetails.PackageName != subDetails.PackageName {
				subDetails.FuncName = subDetails.PackageName + "." + subDetails.FuncName
			}
			lines = append(
				lines,
				lipgloss.JoinHorizontal(
					lipgloss.Left,
					messageStyle.Render(fmt.Sprintf("at %s(", subDetails.FuncName)),
					ctx.s.FileStyle.Render(
						lipgloss.JoinHorizontal(
							lipgloss.Left,
							ansi.SetHyperlink(fmt.Sprintf("%s:%d", subDetails.FilePath, subDetails.Line)),
							fmt.Sprintf("%s:%d", subDetails.FileName, subDetails.Line),
							ansi.ResetHyperlink(),
						),
					),
					messageStyle.Render(")"),
				),
			)
		}
	}
	var logName string
	if len(ctx.loggerName) > 0 {
		logName = fmt.Sprintf("[%s] ", ctx.loggerName)
	}
	blocks := []string{
		messageContainerStyle.Render(
			messageStyle.Render(
				lipgloss.JoinHorizontal(
					lipgloss.Left,
					logName,
					capitalize(matches[3]),
				),
			),
		),
	}
	if len(lines) > 0 {
		blocks = append(blocks, messageContainerStyle.Render(ctx.s.BlockStyle.Render(lipgloss.JoinVertical(lipgloss.Top, lines...))))
	}
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.b.WriteString(
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			ctx.s.TimeStyle.Render(time.Now().Format(ctx.timeFormat)),
			tagStyle.Render(ansi.Truncate(tagName, tagStyle.GetWidth()-3, "...")),
			packageStyle.Render(ansi.Truncate(packageName, packageStyle.GetWidth()-3, "...")),
			iconStyle.Render(),
			lipgloss.JoinVertical(
				lipgloss.Top,
				blocks...,
			),
		),
	)
	ctx.b.WriteRune('\n')
	_, _ = ctx.b.WriteTo(ctx.w)
	if level == FatalLevel {
		os.Exit(1)
	}
}
