package pouch

import (
	"fmt"
	"os/exec"
	"strings"
)

func Inspect(id string) (string, error) {
	cmd := exec.Command("docker", "inspect", id)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("inspect error: %v\noutput: %s", err, out)
	}
	return strings.TrimSpace(string(out)), nil
}