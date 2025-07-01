package pouch

import (
	"fmt"
	"os/exec"
)


func DeleteFile(containerID, filePath string) error {
	cmd := exec.Command("docker", "exec", containerID, "rm", "-f", filePath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("file remove err: %v, output: %s", err, string(out))
	}
	return nil
}
