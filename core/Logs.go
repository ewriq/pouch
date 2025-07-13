package pouch

import (
	"fmt"
	"os/exec"
)

func Logs(id string) (string, error) {
	cmd := exec.Command("docker", "logs", id)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("logs error: %v, output: %s", err, string(out))
	}

	return string(out), nil
}
