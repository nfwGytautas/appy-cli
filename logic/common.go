package logic

import (
	"os"
	"os/exec"
)

func createFile(path string, content string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

func ensureDirectory(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.Mkdir(path, 0755)
	}

	return nil
}

func goFmtFile(path string) error {
	return exec.Command("gofmt", "-w", path).Run()
}
