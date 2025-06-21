package events

import (
	"fmt"
	"os/exec"
	"strings"
)

func Restart(id string) (string, error) {
	cmd := exec.Command("docker", "restart", id)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("restart error: %v\noutput: %s", err, out)
	}
	return strings.TrimSpace(string(out)), nil
}
