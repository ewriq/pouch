package pouch

import (
	"fmt"
	"os/exec"
	"strings"
)

func GetContainerStats(containerID string) (map[string]string, error) {
	stats := make(map[string]string)

	cmd := exec.Command("docker", "stats", "--no-stream", "--format",
		"CPU={{.CPUPerc}}, MEM={{.MemUsage}}, NET={{.NetIO}}, DISK={{.BlockIO}}", containerID)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("stats error: %v\noutput: %s", err, out)
	}

	lines := strings.Split(strings.TrimSpace(string(out)), ", ")
	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(strings.ToLower(parts[0]))
			val := strings.TrimSpace(parts[1])
			stats[key] = val
		}
	}

	return stats, nil
}
