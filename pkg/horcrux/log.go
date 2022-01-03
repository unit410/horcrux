package horcrux

import (
	"fmt"
	"os"
)

func logf(format string, v ...interface{}) {
	if stdoutIsTerminal() {
		fmt.Printf(format, v...)
	}
}
func logln(v ...interface{}) {
	if stdoutIsTerminal() {
		fmt.Println(v...)
	}
}

func warnf(format string, v ...interface{}) {
	format = fmt.Sprintf("Warn: %s", format)
	if stdoutIsTerminal() {
		fmt.Printf(format, v...)
	}
}

var forceTerminal bool

func stdoutIsTerminal() bool {
	if forceTerminal {
		return true
	}
	o, _ := os.Stdout.Stat()
	if (o.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
		// Printing to a Terminal
		return true
	} else {
		// Not printing to a terminal
		return false
	}
}
