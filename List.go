package pouch

import (
	"fmt"
	"os/exec"
)

func List() (string, error) {
	cmd := exec.Command("docker", "ps", "-a") 
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("list error: %v, output: %s", err, string(out))
	}
	
	return string(out), nil
}