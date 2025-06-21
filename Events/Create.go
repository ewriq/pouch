package events

import (
	"fmt"
	"os/exec"
	"strings"
)

type CreateOptions struct {
	Name string 
	Image string
	Port string
	HostDataDir string
	Network string
	Hostname string
	UserUIDGID string
	MemoryLimit string
	EntryPoint string
	CPULimit float64
	EnvVars map[string]string
	Labels map[string]string
} 



func Create(opt CreateOptions) (string, error) {
	args := []string{"create"}

	if opt.Name != "" {
		args = append(args, "--name", opt.Name)
	}

	for k, v := range opt.EnvVars {
		args = append(args, "-e", fmt.Sprintf("%s=%s", k, v))
	}

	if opt.Port != "" {
		args = append(args, "-p", opt.Port+":3306")
	}

	if opt.HostDataDir != "" {
		args = append(args, "-v", opt.HostDataDir)
	}

	args = append(args, "--restart", "always")

	if opt.Network != "" {
		args = append(args, "--network", opt.Network)
	}

	if opt.Hostname != "" {
		args = append(args, "--hostname", opt.Hostname)
	}

	if opt.MemoryLimit != "" {
		args = append(args, "--memory", opt.MemoryLimit)
	}

	if opt.CPULimit > 0 {
		args = append(args, "--cpus", fmt.Sprintf("%.1f", opt.CPULimit))
	}

	for k, v := range opt.Labels {
		args = append(args, "--label", fmt.Sprintf("%s=%s", k, v))
	}

	if opt.UserUIDGID != "" {
		args = append(args, "--user", opt.UserUIDGID)
	}

	args = append(args, "-i", "-t")

	if opt.EntryPoint != "" {
		args = append(args, "--entrypoint", opt.EntryPoint)
	}

	if opt.Image == "" {
		return "", fmt.Errorf("image name is required")
	}

	args = append(args, opt.Image)


	cmd := exec.Command("docker", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("Container create error: %v\noutput: %s", err, string(out))
	}

	return strings.TrimSpace(string(out)), nil
}