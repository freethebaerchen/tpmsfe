package tests

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

func CheckStateForRandValue(projectPath string, filename string, value string) (bool, error) {
	tfstate := filepath.Join(projectPath, filename)
	f, err := os.Open(tfstate)
	if err != nil {
		return false, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), value) {
			return false, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	return true, nil
}
