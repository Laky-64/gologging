package gologging

import (
	"fmt"
	"github.com/Laky-64/gologging/types"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"golang.org/x/term"
	"math"
	"os"
	"strings"
	"time"
)

func internalLog(level Level, message ...any) {
	if level < currentLevel {
		return
	}
	capitalize := func(s string) string {
		return strings.ToUpper(s[:1]) + s[1:]
	}
	var errMessages []string
	for _, x := range message {
		switch x.(type) {
		case error:
			errMessages = append(errMessages, x.(error).Error())
		case string:
			errMessages = append(errMessages, x.(string))
		case bool:
			if x.(bool) {
				errMessages = append(errMessages, "true")
			} else {
				errMessages = append(errMessages, "false")
			}
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			errMessages = append(errMessages, fmt.Sprintf("%d", x))
		case float32, float64:
			errMessages = append(errMessages, fmt.Sprintf("%f", x))
		default:
			if x == nil {
				continue
			}
			errMessages = append(errMessages, fmt.Sprintf("%v", x))
		}
	}
	if len(errMessages) == 0 {
		return
	}
	termWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))
	termWidth = int(math.Max(float64(termWidth), 100))

	timeStyle := lipgloss.NewStyle().
		MarginRight(2)
	tagStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#62c6b7")).
		Width(int(math.Max(float64(termWidth)*0.08, 10)))
	packageStyle := lipgloss.NewStyle().
		Width(int(math.Max(float64(termWidth)*0.15, 15)))
	messageStyle := lipgloss.NewStyle()
	blockContainer := lipgloss.NewStyle().
		MarginLeft(1)
	fileStyle := lipgloss.NewStyle().
		Margin(0, 1).
		Foreground(lipgloss.Color("#61afe1"))
	iconStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#000000")).
		Padding(0, 1).
		MarginRight(1)

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

	var mainDetails *types.CallerInfo
	startSkips := 2
	for {
		details, err := getInfo(startSkips)
		if err == nil {
			mainDetails = details
			break
		}
		if startSkips > 10 {
			break
		}
		startSkips++
	}

	matches := tagRgx.FindStringSubmatch(strings.Join(errMessages, " "))
	tagName := matches[2]
	if len(tagName) == 0 {
		packageInfo := strings.Split(mainDetails.PackageName, ".")
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
			if mainDetails.PackageName != subDetails.PackageName {
				subDetails.FuncName = subDetails.PackageName + "." + subDetails.FuncName
			}
			lines = append(
				lines,
				lipgloss.JoinHorizontal(
					lipgloss.Top,
					messageStyle.Render(fmt.Sprintf("at %s(", subDetails.FuncName)),
					fileStyle.Render(
						lipgloss.JoinHorizontal(
							lipgloss.Top,
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
	blocks := []string{
		messageStyle.Render(capitalize(matches[3])),
	}
	if len(lines) > 0 {
		blocks = append(blocks, blockContainer.Render(lipgloss.JoinVertical(lipgloss.Top, lines...)))
	}
	fmt.Println(
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			timeStyle.Render(time.Now().Format("2006-01-02 15:04:05")),
			tagStyle.Render(ansi.Truncate(tagName, tagStyle.GetWidth()-3, "...")),
			packageStyle.Render(ansi.Truncate(mainDetails.PackageName, packageStyle.GetWidth()-3, "...")),
			iconStyle.Render(),
			lipgloss.JoinVertical(
				lipgloss.Top,
				blocks...,
			),
		),
	)
	if level == FatalLevel {
		os.Exit(1)
	}
}
