package image

import (
	"fmt"
	"os/exec"
)

func Pull(image string) error {
	cmd := exec.Command("docker", "pull", image)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("pull error: %v\n%s", err, string(output))
	}

	return nil
}
