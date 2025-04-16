package utils

import (
	"fmt"
	"io"
	"os"

	"github.com/nfwGytautas/appy-cli/shared"
)

type console struct {
	dbg io.Writer
	err io.Writer
}

type debugWriter struct {
}

type errorWriter struct {
}

func (w *debugWriter) Write(p []byte) (n int, err error) {
	Console.Debug("%s", p)
	return len(p), nil
}

func (w *errorWriter) Write(p []byte) (n int, err error) {
	Console.Debug("%s", p)
	return len(p), nil
}

func (c *console) Header() {
	fmt.Println("------------------------------------------------------------------------------------------------")
	fmt.Printf("                                   Appy CLI: %s, verbose: %t\n", shared.Version, Verbose)
	fmt.Println("------------------------------------------------------------------------------------------------")
}

func (c *console) ClearLines(lineCount int) {
	for i := 0; i < lineCount; i++ {
		fmt.Print("\033[1A\033[K")
	}
}

func (c *console) ClearEntireConsole() {
	fmt.Print("\033[H\033[2J")
	c.Header()
}

func (c *console) Info(format string, a ...any) {
	fmt.Printf("\033[32m%s\033[0m", fmt.Sprintf(format, a...))
}

func (c *console) Warn(format string, a ...any) {
	fmt.Printf("\033[33m⚠️ %s ⚠️\033[0m", fmt.Sprintf(format, a...))
}

func (c *console) Error(format string, a ...any) {
	fmt.Printf("\033[31m❌ %s ❌\033[0m", fmt.Sprintf(format, a...))
}

func (c *console) Debug(format string, a ...any) {
	if Verbose {
		fmt.Printf("\033[90m%s\033[0m", fmt.Sprintf(format, a...))
	}
}

func (c *console) InfoLn(format string, a ...any) {
	fmt.Printf("\033[32m%s\033[0m\n", fmt.Sprintf(format, a...))
}

func (c *console) WarnLn(format string, a ...any) {
	fmt.Printf("\033[33m⚠️ %s ⚠️\033[0m\n", fmt.Sprintf(format, a...))
}

func (c *console) ErrorLn(format string, a ...any) {
	fmt.Printf("\033[31m❌ %s ❌\033[0m\n", fmt.Sprintf(format, a...))
}

func (c *console) DebugLn(format string, a ...any) {
	if Verbose {
		fmt.Printf("\033[90m%s\033[0m\n", fmt.Sprintf(format, a...))
	}
}

func (c *console) Fatal(err error) {
	c.ErrorLn("%s", err)
	os.Exit(1)
}

func (c *console) DebugWriter() io.Writer {
	return c.dbg
}

func (c *console) ErrorWriter() io.Writer {
	return c.err
}
