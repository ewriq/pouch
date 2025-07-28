package container

import (
	"fmt"
	"os/exec"

)

func Inspect(id string) ([]byte, error) {
	cmd := exec.Command("docker", "inspect", id)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return []byte{} , fmt.Errorf("inspect error: %v\noutput: %s", err, out)
	}
	return out, nil
}