package pouch
import (
	"fmt"
	"os/exec"
)

func Start(id string) error {
	cmd := exec.Command("docker", "start", id)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("start error: %v, output: %s", err, string(out))
	}
	return nil
}