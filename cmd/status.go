// Package cmd provides the CLI commands for pvfx.
//
// This file implements the 'status' command which displays
// current LVM and XFS filesystem information.
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show LVM and XFS filesystem status",
	Long: `Display current physical volumes, volume groups, logical volumes,
and XFS filesystem information including size and free space.

This command runs 'pvs', 'vgs', 'lvs', and 'df' to gather information.`,
	RunE: runStatus,
}

// runStatus executes the status command.
// We parse command output rather than using an LVM library to keep
// dependencies minimal and handle various LVM configurations.
func runStatus(cmd *cobra.Command, args []string) error {
	fmt.Println("=== Block Devices & Filesystems ===")
	if err := runBlockDevices(); err != nil {
		return fmt.Errorf("failed to get block devices: %w", err)
	}

	fmt.Println("\n=== Physical Volumes ===")
	if err := runPVS(); err != nil {
		return fmt.Errorf("failed to get PV info: %w", err)
	}

	fmt.Println("\n=== Volume Groups ===")
	if err := runVGS(); err != nil {
		return fmt.Errorf("failed to get VG info: %w", err)
	}

	fmt.Println("\n=== Logical Volumes ===")
	if err := runLVS(); err != nil {
		return fmt.Errorf("failed to get LV info: %w", err)
	}

	fmt.Println("\n=== XFS Filesystems ===")
	if err := runXFSInfo(); err != nil {
		return fmt.Errorf("failed to get XFS info: %w", err)
	}

	fmt.Println("\n=== Performance Metrics (Disk I/O) ===")
	if err := runDiskStats(); err != nil {
		return fmt.Errorf("failed to get disk stats: %w", err)
	}

	return nil
}

// runPVS runs 'pvs' and prints the output.
// exec.Command returns a Cmd struct - we use CombinedOutput to
// capture both stdout and stderr in one call.
func runPVS() error {
	// Context: exec.Command separates the program name from its arguments.
	// This is more secure than shell injection because the command
	// and args are passed as a slice, not interpreted by a shell.
	cmd := exec.Command("pvs", "--units", "g", "--noheadings", "-o", "pv_name,vg_name,pv_size,pv_free,pv_used")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("pvs failed: %w\n%s", err, out)
	}
	printFormattedTable(out, []string{"PV", "VG", "Size", "Free", "Used"})
	return nil
}

// runVGS runs 'vgs' to show volume group information.
func runVGS() error {
	cmd := exec.Command("vgs", "--units", "g", "--noheadings", "-o", "vg_name,vg_size,vg_free,vg_extent_size,vg_extent_count,vg_free_count")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("vgs failed: %w\n%s", err, out)
	}
	printFormattedTable(out, []string{"VG", "Size", "Free", "Extent Size", "Total Extents", "Free Extents"})
	return nil
}

// runLVS runs 'lvs' to show logical volume information.
func runLVS() error {
	cmd := exec.Command("lvs", "--units", "g", "--noheadings", "-o", "lv_name,vg_name,lv_size,lv_path,devices")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("lvs failed: %w\n%s", err, out)
	}
	printFormattedTable(out, []string{"LV", "VG", "Size", "Path", "Devices"})
	return nil
}

// runXFSInfo runs 'df' filtered to XFS filesystems to show mount points and usage.
func runXFSInfo() error {
	cmd := exec.Command("df", "-Th", "-t", "xfs")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("df failed: %w\n%s", err, out)
	}
	// df output is already formatted, just print it
	fmt.Print(string(out))
	return nil
}

// printFormattedTable takes raw command output and column headers,
// splits the output into lines, and prints a formatted table.
// This is a simple implementation - for production tools you might
// use a table formatting library.
func printFormattedTable(data []byte, headers []string) {
	// Idiomatic Go: Use text/tabwriter from the standard library to 
	// align columns automatically. It's much more robust than manual fmt.Printf.
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	
	// Print headers
	fmt.Fprintln(w, strings.Join(headers, "\t"))
	// Print separator line
	sep := make([]string, len(headers))
	for i := range sep {
		sep[i] = strings.Repeat("-", len(headers[i]))
	}
	fmt.Fprintln(w, strings.Join(sep, "\t"))
	
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		
		// Convert whitespace gaps into tabs for tabwriter
		fields := strings.Fields(line)
		
		// If we have more fields than headers, combine the extras into the last column
		if len(fields) > len(headers) {
			lastIdx := len(headers) - 1
			fields[lastIdx] = strings.Join(fields[lastIdx:], " ")
			fields = fields[:len(headers)]
		}
		
		fmt.Fprintln(w, strings.Join(fields, "\t"))
	}
	// Idiomatic Go: Always check the scanner for errors after the loop.
	// It catches IO errors that happen during reading.
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error reading command output: %v\n", err)
	}
	w.Flush()
}

// runBlockDevices executes lsblk. It displays raw storage topologies perfectly
// on almost all Linux distributions (showing RAID, LVM, partitions).
func runBlockDevices() error {
	// Idiomatic Go CLI: Executing standard system commands is often preferred over
	// linking CGO libraries (like libblkid) to maintain static binaries and
	// cross-compilation capability.
	cmd := exec.Command("lsblk", "-o", "NAME,MAJ:MIN,RM,SIZE,RO,TYPE,FSTYPE,MOUNTPOINT")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("lsblk failed: %w\n%s", err, out)
	}
	fmt.Print(string(out))
	return nil
}

// runDiskStats reads from /proc/diskstats directly to gather performance metrics.
func runDiskStats() error {
	// Idiomatic Go: For system metrics, parsing /proc directly is preferred over
	// calling external binaries like 'iostat'. It avoids an external dependency
	// (sysstat package) which might not be installed on every server.
	file, err := os.Open("/proc/diskstats")
	if err != nil {
		return fmt.Errorf("open /proc/diskstats: %w", err)
	}
	// Idiomatic Go: defer ensures the file is closed when the function returns,
	// preventing file descriptor leaks even if errors occur below.
	defer file.Close()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "DEVICE\tREADS_COMPLETED\tWRITES_COMPLETED\tIO_IN_PROGRESS")
	fmt.Fprintln(w, "------\t---------------\t----------------\t--------------")

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		// /proc/diskstats format has at least 14 fields for block devices.
		// Field 2 is device name, 3 is reads completed, 7 is writes completed, 11 is I/Os currently in progress.
		if len(fields) >= 14 {
			devName := fields[2]
			// Skip loop devices and ram disks to reduce noise
			if strings.HasPrefix(devName, "loop") || strings.HasPrefix(devName, "ram") {
				continue
			}
			
			reads := fields[3]
			writes := fields[7]
			ioInProgress := fields[11]
			
			// Fprintf handles writing directly to our tabwriter
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", devName, reads, writes, ioInProgress)
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("reading /proc/diskstats: %w", err)
	}
	w.Flush()
	return nil
}

func init() {
	// Add status command to root command.
	// Cobra automatically adds this command to the root command's
	// subcommands list via the init() function.
	rootCmd.AddCommand(statusCmd)

	// Add flags specific to status command
	// --json could be added here for machine-readable output
}