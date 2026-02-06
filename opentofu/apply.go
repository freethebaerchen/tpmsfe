package opentofu

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func TofuApply(workingDir string) (string, error) {
	init := exec.Command("tofu", "-chdir="+workingDir, "init")
	cmd := exec.Command("tofu", "-chdir="+workingDir, "apply", "-auto-approve")

	var stdout, stderr bytes.Buffer
	init.Stdout = &stdout
	init.Stderr = &stderr
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := init.Run(); err != nil {
		return "", fmt.Errorf("error running init: %v\nstderr: %s", err, stderr.String())
	}

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("error running apply: %v\nstderr: %s", err, stderr.String())
	}

	outputs := parseOutput(stdout.String())
	return outputs, nil
}

func parseOutput(output string) string {
	lines := strings.Split(output, "\n")
	if len(lines) == 0 {
		return ""
	}

	var lastLine string
	for i := len(lines) - 1; i >= 0; i-- {
		if strings.TrimSpace(lines[i]) != "" {
			lastLine = lines[i]
			break
		}
	}

	if lastLine == "" {
		return ""
	}

	fields := strings.Fields(lastLine)
	if len(fields) < 3 {
		return ""
	}

	return strings.Trim(fields[2], `"`)
}

