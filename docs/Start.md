### `pouch.Start` Function Explanation

This Go function, `Start`, is part of the `pouch` package. Its primary role is to start a stopped Docker container by executing the `docker start` command under the hood. It returns only an `error` value to indicate success or failure.

#### Code Block

```go
package pouch

import (
	"fmt"
	"os/exec"
)

// Start starts a stopped Docker container with the given ID or name.
func Start(id string) error {
	// Create the command: docker start [id]
	cmd := exec.Command("docker", "start", id)

	// Execute the command and capture combined output (stdout + stderr)
	out, err := cmd.CombinedOutput()
	if err != nil {
		// Return a detailed error message with both Go error and Docker output
		return fmt.Errorf("start error: %v, output: %s", err, string(out))
	}

	// Return nil if successful
	return nil
}
```

---

### Function Details

1. **Command Construction**:
   The line `exec.Command("docker", "start", id)` constructs a command equivalent to typing `docker start [containerID]` in the terminal.

2. **Command Execution and Output Capture**:
   `cmd.CombinedOutput()` runs the command and captures both `stdout` and `stderr`. If the command fails (e.g., if the container does not exist), Docker writes the error to `stderr`. On success, it prints the container ID or name to `stdout`.

3. **Error Handling**:
   The function checks whether the command failed using `if err != nil`. If Docker encounters a problem—such as the container not being found, or the Docker daemon being unavailable—the function wraps the error and the Docker output into a single, descriptive message using `fmt.Errorf`.

4. **Success Case**:
   If everything works, the function simply returns `nil`, indicating the operation succeeded. It does not return the container's ID or name—just the status of the operation.

---

### Parameters

* `id (string)`: The ID or name of the Docker container to start.

---

### Return Value

* `error`:

  * Returns `nil` if the container is started successfully.
  * Returns a detailed error if the container cannot be found, if Docker is not running, or if another problem occurs.

---

### Requirements

* Docker must be installed on the host machine.
* The `docker` command must be accessible via the system’s `PATH`.

---

### Example Usage

The following example demonstrates how to start a stopped container named `my-database-container` using the `pouch.Start` function:

```go
package main

import (
	"log"
	"github.com/your-username/your-project/pouch" // Replace with your actual import path
)

func main() {
	containerToStart := "my-database-container"

	log.Printf("Starting container '%s'...", containerToStart)

	// Call pouch.Start
	err := pouch.Start(containerToStart)
	if err != nil {
		log.Fatalf("Failed to start container: %v", err)
	}

	log.Printf("Container '%s' started successfully!", containerToStart)
}
```

