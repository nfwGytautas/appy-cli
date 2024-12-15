package logic

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func getCurrentDirectory() string {
	dir, err := os.Getwd()
	errorCheckAndPanic(err)

	return dir
}

func ensureDirectoryIsEmpty(path string) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	if len(files) > 0 {
		return os.ErrExist
	}

	return nil
}

func errorCheckAndPanic(err error) {
	if err != nil {
		panic(err)
	}
}

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
	pathComponents := strings.Split(path, "/")

	prefix := ""
	for _, component := range pathComponents {
		if _, err := os.Stat(prefix + component); os.IsNotExist(err) {
			err := os.Mkdir(prefix+component, 0755)
			if err != nil {
				return err
			}
		}
		prefix += component + "/"
	}

	return nil
}

func ensureDirectories(paths []string) error {
	for _, path := range paths {
		err := ensureDirectory(path)
		if err != nil {
			return err
		}
	}

	return nil
}

func goFmtFile(path string) error {
	return exec.Command("gofmt", "-w", path).Run()
}

func goImportFile(path string) error {
	return exec.Command("goimports", "-w", path).Run()
}

func goPipeline(path string) error {
	err := goImportFile(path)
	if err != nil {
		return err
	}

	err = goFmtFile(path)
	if err != nil {
		return err
	}

	return nil
}

func beutifyContent(content string) string {
	content = strings.ReplaceAll(content, "\n\n", "\n")

	return content
}

func chainErrorChecks(errs ...error) {
	for _, err := range errs {
		errorCheckAndPanic(err)
	}
}

func walkFilesInDirectory(path string, callback func(string, string)) error {
	walkFn := func(s string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			callback(s, filename(s))
		}

		return nil
	}

	// Check if directory
	// info, err := os.Stat(path)
	// if err != nil {
	// 	if info.IsDir() {
	// 		err = filepath.WalkDir(path, walkFn)
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// }

	return filepath.WalkDir(path, walkFn)
}

func readFileLines(path string) ([]string, error) {
	fileIO, err := os.OpenFile(path, os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	defer fileIO.Close()

	rawBytes, err := io.ReadAll(fileIO)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(rawBytes), "\n")
	return lines, nil
}

func filename(path string) string {
	return strings.Split(filepath.Base(path), ".")[0]
}

func getUserInput(prompt string, val *string) error {
	fmt.Print(prompt)
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		return err
	}

	*val = input

	return nil
}

func stringValueIsOneOf(value string, values []string) bool {
	for _, v := range values {
		if value == v {
			return true
		}
	}

	return false
}

func capitalizeFirstLetter(in string) string {
	return strings.ToUpper(in[:1]) + strings.ToLower(in[1:])
}

func uniqueStrings(in []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	for _, entry := range in {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}

	return list
}

func getGoModule() string {
	goMod, err := os.ReadFile("go.mod")
	if err != nil {
		return ""
	}

	lines := strings.Split(string(goMod), "\n")
	for _, line := range lines {
		if strings.Contains(line, "module") {
			return strings.Split(line, " ")[1]
		}
	}

	return ""
}
