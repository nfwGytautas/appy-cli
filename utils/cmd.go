package utils

import (
	"fmt"
	"os/exec"
	"strings"
)

func RunCommand(dir string, action string) error {
	cmd := exec.Command(strings.Split(action, " ")[0], strings.Split(action, " ")[1:]...)

	Console.DebugLn(cmd.String())

	if dir != "" {
		cmd.Dir = dir
	}

	cmd.Stdout = Console.DebugWriter()
	cmd.Stderr = Console.DebugWriter()
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run command: %v", err)
	}

	return nil
}

func RunTools(file string, tools []string) error {
	for _, tool := range tools {
		toolArgs := strings.Split(tool, " ")
		toolArgs = append(toolArgs, file)
		cmd := exec.Command(toolArgs[0], toolArgs[1:]...)
		cmd.Stdout = Console.DebugWriter()
		cmd.Stderr = Console.DebugWriter()
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to run command: %v", err)
		}
	}

	return nil
}
