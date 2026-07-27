package utils

import (
	"fmt"
	"os"
	"strings"
)

const (
	colorReset  = "\033[0m"
	colorBold   = "\033[1m"
	colorDim    = "\033[2m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
)

var useColor = true

func init() {
	if os.Getenv("NO_COLOR") != "" {
		useColor = false
	}
}

func c(color, text string) string {
	if !useColor {
		return text
	}
	return color + text + colorReset
}

func Print(args ...interface{}) {
	fmt.Fprint(os.Stdout, args...)
}

func Println(args ...interface{}) {
	fmt.Fprintln(os.Stdout, args...)
}

func Printf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stdout, format, args...)
}

func Info(msg string) {
	fmt.Fprintln(os.Stdout, c(colorBlue, "ℹ")+"  "+msg)
}

func Success(msg string) {
	fmt.Fprintln(os.Stdout, c(colorGreen, "✔")+"  "+msg)
}

func Warn(msg string) {
	fmt.Fprintln(os.Stdout, c(colorYellow, "⚠")+"  "+msg)
}

func Error(msg string) {
	fmt.Fprintln(os.Stderr, c(colorRed, "✗")+"  "+msg)
}

func Fatal(msg string, err error) {
	Error(msg + ": " + err.Error())
	os.Exit(1)
}

func Debug(msg string) {
	fmt.Fprintln(os.Stdout, c(colorGray, "…")+"  "+msg)
}

func Step(title string) {
	fmt.Fprintln(os.Stdout)
	fmt.Fprintln(os.Stdout, c(colorBold+colorCyan, "── "+title+" "))
	fmt.Fprintln(os.Stdout, c(colorGray, strings.Repeat("─", 40)))
}
