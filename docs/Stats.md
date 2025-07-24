

### `pouch.ContainerStats` Function Explanation

This Go function, `ContainerStats`, resides in the `pouch` package and provides a programmatic way to retrieve **real-time resource usage statistics** (CPU, memory, network, disk I/O) of a **running Docker container**. It executes the `docker stats` command with formatting options and parses the output into a Go `map`.

---

#### Code Block

```go
package pouch

import (
	"fmt"
	"os/exec"
	"strings"
)

// ContainerStats returns real-time resource usage stats of a running container
// such as CPU, memory, network, and disk I/O in a map format.
func ContainerStats(containerID string) (map[string]string, error) {
	stats := make(map[string]string)

	cmd := exec.Command("docker", "stats", "--no-stream", "--format",
		"CPU={{.CPUPerc}}, MEM={{.MemUsage}}, NET={{.NetIO}}, DISK={{.BlockIO}}", containerID)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("stats error: %v\noutput: %s", err, out)
	}

	// Example output: "CPU=0.15%, MEM=12.5MiB / 1.95GiB, NET=648B / 0B, DISK=0B / 8.19kB"
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
```

---

### How It Works

The function works in two main phases: **data collection** and **data parsing**.

1. **Data Collection (via docker stats)**:

   * `--no-stream`: ensures the command fetches a single snapshot rather than a live stream.
   * `--format`: customizes the output to a predictable, easily parsable format like `KEY=VALUE`.

2. **Data Parsing**:

   * Trims any trailing whitespace.
   * Splits the string by `, ` into stat segments (e.g., `["CPU=...", "MEM=..."]`).
   * Splits each segment by `=` using `SplitN` (in case the value contains `=`).
   * Normalizes the keys (e.g., `"CPU"` → `"cpu"`) and trims the values.
   * Populates the result `map[string]string`.

---

### Parameters

* `containerID (string)`: The name or ID of a **running** Docker container.

---

### Return Values

* `map[string]string`: A map containing the following keys:

  * `"cpu"` – CPU usage percentage.
  * `"mem"` – Memory usage (e.g., `12.5MiB / 1.95GiB`).
  * `"net"` – Network I/O (e.g., `648B / 0B`).
  * `"disk"` – Block I/O (e.g., `0B / 8.19kB`).
* `error` – Returns `nil` if successful. Returns a detailed error if Docker fails or the container is not found/running.

---

### Requirements & Notes

* Docker must be installed and accessible in the system’s `$PATH`.
* The target container must be **running**; otherwise, no stats will be returned.

---

### Example Usage

Here’s how to retrieve and display stats for a container named `my-service`:

```go
package main

import (
	"fmt"
	"log"
	"github.com/your-username/your-project/pouch" // Adjust the import path as needed
)

func main() {
	containerID := "my-service"

	stats, err := pouch.ContainerStats(containerID)
	if err != nil {
		log.Fatalf("Failed to get container stats: %v", err)
	}

	fmt.Printf("Live Stats for '%s':\n", containerID)
	fmt.Printf("  - CPU: %s\n", stats["cpu"])
	fmt.Printf("  - Memory: %s\n", stats["mem"])
	fmt.Printf("  - Network I/O: %s\n", stats["net"])
	fmt.Printf("  - Disk I/O: %s\n", stats["disk"])
}
```

