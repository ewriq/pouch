package pouch

import (
	"fmt"
	"os/exec"
	"strings"
)

func Remove(id string, force bool) (string, error) {
	args := []string{"rm"}
	if force {
		args = append(args, "-f")
	}

	args = append(args, id)
	cmd := exec.Command("docker", args...)
	out, err := cmd.CombinedOutput()

	if err != nil {
		return "", fmt.Errorf("remove error: %v\noutput: %s", err, out)
	}
	return strings.TrimSpace(string(out)), nil
}