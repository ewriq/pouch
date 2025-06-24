package pouch

import (
	"fmt"
	"os/exec"
	"strings"
)

func ListFiles(containerID, containerPath string) ([]string, error) {
	cmd := exec.Command("docker", "exec", containerID, "ls", "-l", containerPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("File listing err: %v, %s", err, string(out))
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")

	return lines, nil
}
