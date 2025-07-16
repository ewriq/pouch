package container

import (
	"fmt"
	"os/exec"
)

func Stop(id string) error {
	cmd := exec.Command("docker", "stop", id)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("start error: %v, output: %s", err, string(out))
	}
	return nil
}