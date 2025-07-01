package pouch

import (
	"fmt"
	"os/exec"
)

func Exec(id string, command []string) (string, error) {
	args := append([]string{"exec", id}, command...)
	cmd := exec.Command("docker", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("exec error: %v\noutput: %s", err, out)
	}
	return string(out), nil
}
