package utils

import (
	"fmt"
	"os/exec"
)

func RunCommand(dir string, args ...string) error {
	// Pull repository
	cmd := exec.Command(args[0], args[1:]...)

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
