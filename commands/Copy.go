package pouch

import (
	"fmt"
	"os/exec"
)

func CopyToContainer(ID, HostPath, TargetPath string) error {
	cmd := exec.Command("docker", "cp", HostPath, fmt.Sprintf("%s:%s", ID, TargetPath))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker cp hatası: %v, çıktı: %s", err, string(out))
	}
	return nil
}
