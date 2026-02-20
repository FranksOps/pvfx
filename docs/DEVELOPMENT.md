# Developer Documentation

Welcome to `pvfx`. This is a CLI tool built in Go to inspect and modify Linux block devices, LVM topologies, and XFS filesystems.

As a systems engineer learning Go, this tool is meant to be a robust, practical introduction to idiomatic Go.

## Project Structure

This project follows a standard Go project layout, particularly optimized for CLI tools using Cobra:

- `main.go`: The absolute bare-minimum entrypoint for the compiled binary. All it does is call `cmd.Execute()`.
- `cmd/`: This directory contains all the CLI commands, which is the standard pattern for Go CLI tools.
    - `root.go`: Initializes the Cobra CLI framework and defines global flags.
    - `status.go`: Inspects the system (PVs, VGs, LVs, XFS mount points, block devices, `/proc/diskstats`).
    - `grow.go`: (WIP) Will provide the logic to resize LVM and XFS objects.

## Go Concepts to Keep in Mind

If you are a sysadmin used to Bash or Python, these are the core Go concepts utilized in this tool:

### 1. `os/exec` instead of backticks/os.system()
In Go, calling out to Bash commands is not done via an unsafe shell interpreter. Instead, `exec.Command("pvs", "--units", "g")` creates an array of arguments passed directly to the `execve` syscall. This prevents shell injection vulnerabilities.

### 2. Tabwriter for CLI Output
We don't try to manually align columns using spaces and `fmt.Printf`. The Go standard library provides `text/tabwriter`, which buffers the output, calculates the maximum width of each column, and flushes it perfectly aligned.

### 3. Reading `/proc`
Instead of adding dependencies to run `iostat` or `sysstat`, we read directly from `/proc/diskstats`. Go's `os.Open` and `bufio.Scanner` are extremely efficient at streaming line-by-line metrics from pseudo-filesystems.

### 4. `defer`
Whenever we open a file or a network connection, we immediately write `defer file.Close()` on the next line. This guarantees the file descriptor is closed when the function exits, even if an error is returned halfway through the function.

### 5. `RunE` instead of `Run`
In Cobra, we use `RunE` (which returns an `error`) instead of `Run` (which doesn't). This forces us to bubble up our errors (`return fmt.Errorf("failed to do X: %w", err)`) to the top level, rather than calling `os.Exit(1)` deep inside a utility function. This is critical for testability and idiomatic error handling.

## Adding a New Command

If you want to add a new command (e.g., `shrink`), you can create a new file `cmd/shrink.go`:

```go
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var shrinkCmd = &cobra.Command{
	Use:   "shrink",
	Short: "Shrink a volume",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Shrinking...")
		return nil
	},
}

func init() {
    // This init() function automatically runs on startup
    // and registers this command with the root CLI.
	rootCmd.AddCommand(shrinkCmd)
}
```