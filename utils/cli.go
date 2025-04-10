package utils

import "fmt"

func ClearLines(lineCount int) {
	for i := 0; i < lineCount; i++ {
		fmt.Print("\033[1A\033[K")
	}
}

func ConsoleWarn(format string, a ...any) {
	fmt.Printf("\033[33m⚠️  Warning\033[0m: %s\n", fmt.Sprintf(format, a...))
}

func ConsoleError(format string, a ...any) {
	fmt.Printf("\033[31m❌ Error\033[0m: %s\n", fmt.Sprintf(format, a...))
}
