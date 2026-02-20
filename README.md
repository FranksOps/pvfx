# pvfx

`pvfx` is a CLI tool written in Go to interact with Linux Logical Volume Manager (LVM) and XFS filesystems. It provides an intuitive interface for system administrators and DevOps engineers to manage, inspect, and grow volumes without having to remember the intricate syntax of the underlying `lvm2` and `xfs_growfs` toolsets.

> **Status:** Work in Progress. Currently being developed as a learning project to build idiomatic CLI tools in Go.

## Features (Current and Planned)

- **status:** View LVM structures (PVs, VGs, LVs), XFS filesystem usage, block device layouts, and raw `/proc/diskstats` disk I/O metrics—all formatted neatly.
- **grow (WIP):** Grow LVM logical volumes and automatically extend the underlying XFS filesystem.

## Prerequisites

Because `pvfx` calls out to standard system utilities, you must have the following packages installed (they are present by default on almost all RHEL/Ubuntu-based enterprise systems):

- `lvm2` (provides `pvs`, `vgs`, `lvs`)
- `xfsprogs` (provides XFS utilities)
- `util-linux` (provides `lsblk`, `df`)

## Installation

Assuming you have Go 1.22+ installed:

```bash
git clone https://github.com/FranksOps/pvfx.git
cd pvfx
go build -o pvfx .
```

To install it directly to your `$GOPATH/bin`:

```bash
go install github.com/FranksOps/pvfx@latest
```

## Usage

```bash
# View help and available commands
./pvfx help

# View status of storage topologies and disk IO
./pvfx status

# (WIP) Grow a volume
./pvfx grow
```

## Development & Go Learning Notes

This project is built following idiomatic Go practices for system utilities:

- **Standard Library First:** We avoid CGO (e.g., `libblkid` or `libdevmapper` C bindings) to ensure the resulting binary is completely statically linked and easy to deploy across different Linux distributions.
- **Cobra for CLI:** We use `spf13/cobra` to handle subcommands, flags, and help text generation, which is the industry standard for Go CLI applications (used by Kubernetes, Docker, etc.).
- **Direct System Calls:** When reading system metrics, we prefer parsing files like `/proc/diskstats` using `bufio.Scanner` rather than relying on external binaries like `iostat`.
- **Robust Command Execution:** When we do need to call external commands (like `lvs` or `df`), we use `os/exec` with explicit arguments to prevent shell injection vulnerabilities.
