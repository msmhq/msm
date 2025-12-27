# Minecraft Server Manager

A command-line tool for running and managing multiple Minecraft server instances on a single host.

**MSM v0.12.0** is a complete rewrite in Go, replacing the original bash init script with a modern, type-safe CLI.

## Contents

- [Requirements](#requirements)
- [Quick Start](#quick-start)
- [Installation](#installation)
- [Features](#features)
- [Running Multiple Servers](#running-multiple-servers)
- [Importing Existing Worlds](#importing-existing-worlds)
- [Configuration](#configuration)
- [Permissions](#permissions)
- [Common Commands](#common-commands)
- [RAM Disk Support](#ram-disk-support)
- [Fabric Mod Loader](#fabric-mod-loader)
- [Full Command Reference](#full-command-reference)
- [Shell Completion](#shell-completion)
- [Upgrading from Bash MSM](#upgrading-from-bash-msm)
- [Troubleshooting](#troubleshooting)
- [Development](#development)
- [License](#license)
- [Acknowledgements](#acknowledgements)
- [Support](#support)

## Requirements

### Platform Support

| Platform | Status | Notes |
|----------|--------|-------|
| Linux | Full support | Primary target platform |
| macOS | Full support | Requires GNU screen (`brew install screen`) |
| Windows | WSL only | Native Windows not supported |

MSM requires a Unix-like environment due to its reliance on GNU screen for process management, POSIX filesystem semantics, and `/dev/shm` for RAM disk features.

### Dependencies

| Dependency | Required | Purpose |
|------------|----------|---------|
| Go 1.21+ | Build only | Not needed if using pre-built binary |
| GNU screen | Yes | Console and process management |
| rsync | Optional | rsync backups, RAM disk sync |
| rdiff-backup | Optional | Incremental backups with rotation |

### Minecraft Version Support

MSM supports **Minecraft 1.7.0 and later** (Java Edition), including the new year-based versioning scheme (e.g., `26.1`, `26.2`). Servers running Minecraft 1.6.x or earlier are not supported.

## Quick Start

```bash
# 1. Install MSM (see Installation section for details)
make build && sudo make setup && sudo make install

# 2. Create a server (auto-configures port, EULA, etc.)
msm server create survival

# 3. Download the latest Minecraft server jar
msm jar download survival

# 4. Start the server
msm start survival

# 5. Attach to console (detach with Ctrl+A, D)
msm console survival
```

## Installation

### Prerequisites

**Linux (Debian/Ubuntu):**
```bash
sudo apt update
sudo apt install golang-go screen rsync openjdk-21-jre-headless
```

**Linux (RHEL/Fedora):**
```bash
sudo dnf install golang screen rsync java-21-openjdk-headless
```

**macOS:**
```bash
brew install go screen rsync openjdk@21
```

**Java version requirements:**

| Minecraft Version | Minimum Java |
|-------------------|--------------|
| 1.21+ | Java 21 |
| 1.18 - 1.20.x | Java 17 |
| 1.17 | Java 16 |
| 1.12 - 1.16.x | Java 8 |

### From Source

```bash
git clone https://github.com/msmhq/msm.git
cd msm
make build           # Requires Go 1.21+
sudo make setup      # Creates minecraft user and /opt/msm directories
sudo make install    # Installs binary, config, and cron job
```

### Pre-built Binaries (Recommended)

Download from the [releases page](https://github.com/msmhq/msm/releases) - no Go required.

```bash
# Clone repo for Makefile and config files
git clone https://github.com/msmhq/msm.git
cd msm

# Download binary to bin/ (Linux amd64 example)
mkdir -p bin
curl -L https://github.com/msmhq/msm/releases/latest/download/msm-linux-amd64 -o bin/msm
chmod +x bin/msm

# Set up user, directories, config, and cron
sudo make setup      # Creates minecraft user and /opt/msm directories
sudo make install    # Installs binary, config, and cron job
```

Optionally enable auto-start on boot:
```bash
sudo make systemd-install
```

For multi-user setups where multiple people manage their own servers, see [PERMISSIONS.md](PERMISSIONS.md)

### Verify Installation

```bash
msm version
msm config
```

## Features

- **Multi-instance management** - Run multiple Minecraft server instances on a single host
- **Server lifecycle** - Start, stop, restart with configurable delays and player warnings
- **World backups** - Zip, rsync, or rdiff-backup with automatic rotation
- **RAM disk support** - Load worlds into RAM for reduced lag
- **Jar management** - Centralize jars and share across multiple server instances
- **Fabric mod loader** - Native support for Fabric with version compatibility checks
- **Player management** - Allowlist, operators, bans (players and IPs)
- **Log rolling** - Automatic log compression and archival
- **Console access** - Direct access to server console via screen

## Running Multiple Servers

MSM is designed to run multiple Minecraft server instances simultaneously on a single host. Each MSM "server" is an independent Minecraft instance with its own world, port, and configuration.

It is generally designed for one world per server instance.

### How It Works

| What you create | What you get |
|-----------------|--------------|
| `msm server create survival` | Independent server on port 25565 |
| `msm server create creative` | Independent server on port 25566 |
| `msm server create skyblock` | Independent server on port 25567 |

When you create a server, MSM automatically:
- Assigns the next available port (starting from 25565)
- Generates `server.properties` with sensible defaults
- Accepts the Minecraft EULA on your behalf
- Creates the directory structure for worlds

### Starting Multiple Servers

```bash
# Start all servers at once
msm start

# Or start individually
msm start survival
msm start creative

# Check what's running
msm server list
```

```
NAME       STATUS   PORT   RAM SIZE  VERSION  FABRIC
survival   running  25565  2048M     1.21.4   yes
creative   running  25566  2048M     1.21.4   -
skyblock   stopped  25567  4096M     1.20.4   -
```

```bash
# Inspect worlds for a server, even though there's typically only one
msm worlds list survival
```

```
NAME            SERVER    LOCATION  VERSION
world           survival  RAM       1.21.4
world_nether    survival  disk      1.21.4
world_the_end   survival  disk      1.21.4
```

Players connect to different servers using their port: `yourhost:25565`, `yourhost:25566`, etc.

### Resource Planning

Each Minecraft server requires its own RAM allocation (configured in `server.conf`). Plan accordingly:

| Servers | RAM Each | Total Needed |
|---------|----------|--------------|
| 3 | 2GB | ~8GB (including OS overhead) |
| 5 | 4GB | ~24GB |

## Importing Existing Worlds

See [IMPORTING.md](IMPORTING.md) for detailed instructions on importing worlds from:
- Single-player saves
- Other Minecraft servers
- Modpacks with custom jars

Quick start:
```bash
msm server create survival
sudo rsync -r --chown=minecraft:minecraft /path/to/world/ \
    /opt/msm/servers/survival/worldstorage/world/
msm jar download survival
msm start survival
```

## Configuration

MSM reads its configuration from `/etc/msm.conf`. Override this path with `--config` flag or `MSM_CONF` environment variable.

### Key Settings

```bash
# User to run servers as
USERNAME="minecraft"

# Storage paths
SERVER_STORAGE_PATH="/opt/msm/servers"
JAR_STORAGE_PATH="/opt/msm/jars"

# Server defaults
DEFAULT_RAM="1024"
DEFAULT_STOP_DELAY="10"

# Minecraft server.properties defaults (used when creating new servers)
DEFAULT_SERVER_PORT="25565"
DEFAULT_RENDER_DISTANCE="10"
DEFAULT_MAX_PLAYERS="20"
DEFAULT_DIFFICULTY="normal"
DEFAULT_GAMEMODE="survival"

# Cron schedule
CRON_MAINTENANCE_HOUR="5"
CRON_ARCHIVE_RETENTION_DAYS="30"
```

### Per-Server Configuration

Each server can override defaults in `<server>/server.conf`:

```bash
USERNAME="minecraft"
RAM="2048"
JAR_PATH="paper.jar"
STOP_DELAY="30"
```

## Permissions

MSM enforces ownership-based permissions:

- Each server has an owner (the `USERNAME` in `server.conf`)
- You can only manage servers you own, unless you're root
- When you create a server, it's automatically owned by you
- The systemd service runs as root, so it can start all servers on boot

```bash
# Alice creates a server (owned by alice)
alice$ msm server create my-server

# Alice can manage it
alice$ msm start my-server

# Bob cannot
bob$ msm stop my-server
Error: permission denied: server "my-server" is owned by user "alice"
```

See [PERMISSIONS.md](PERMISSIONS.md) for the full guide including multi-user setup.

## Common Commands

### Server Management

```bash
msm server list              # List all servers
msm server create <name>     # Create new server (auto-assigns port)
msm server delete <name>     # Delete a server
msm server rename <old> <new>  # Rename a server
msm server init <name>       # Initialize missing config files (for imported worlds)
msm server config <server>   # Show per-server configuration
msm server config set <server> <key> <value>  # Set a config value
```

### Server Lifecycle

```bash
msm start                    # Start all servers
msm stop [--now]             # Stop all servers
msm start <server>           # Start specific server
msm stop <server> [--now]    # Stop specific server
msm restart <server>         # Restart server
msm status <server>          # Check if running
msm console <server>         # Attach to console
msm cmd <server> "command"   # Send console command
```

### Backups

```bash
msm worlds backup --all      # Backup all worlds (all servers)
msm worlds backup <server>   # Backup worlds for one server
msm backup <server>          # Full server backup (includes configs)
```

### Player Management

```bash
msm allowlist add <server> <player>
msm allowlist remove <server> <player>
msm op <server> <player>
msm kick <server> <player> [reason]
```

### System Administration

```bash
sudo msm setup               # Fix directory ownership and permissions
msm config                   # Show global configuration
msm server config <server>   # Show per-server configuration
msm server config set <server> <key> <value>  # Set per-server config value
msm cron install             # Install maintenance cron job
```

### Jar Management

```bash
# Show current jar
msm jar <server>

# Download vanilla Minecraft directly
msm jar download <server>           # Latest version
msm jar download <server> 1.21.11    # Specific version

# Jar groups for custom jars (Paper, Spigot, etc.)
msm jargroup create <name> <url>
msm jargroup getlatest <name>
msm jar link <server> <jargroup>
```

## RAM Disk Support

Load worlds into RAM (`/dev/shm`) for reduced I/O latency:

```bash
# Check RAM status for a world
msm worlds ram survival world

# Enable RAM for a world
msm worlds ram on survival world

# Disable RAM for a world
msm worlds ram off survival world

# Force sync to disk
msm worlds todisk survival
```

A background sync daemon automatically syncs RAM to disk every 2 minutes while servers are running. Additionally, `msm stop` automatically syncs all RAM worlds to disk before shutting down, ensuring data safety during planned reboots and system updates.

**Warning:** In case of unexpected shutdown (power loss, kernel panic), you may lose up to 2 minutes of changes.

### Why worldstorage/?

MSM stores worlds in a `worldstorage/` subdirectory and creates symlinks on server start. This indirection enables RAM disk support—when you enable RAM for a world, MSM symlinks it to `/dev/shm` instead of the local directory. Without this architecture, RAM disk wouldn't be possible.

See [IMPORTING.md](IMPORTING.md) for details on how MSM handles world storage.

## Fabric Mod Loader

MSM natively supports [Fabric](https://fabricmc.net/):

```bash
# Enable Fabric
msm fabric set survival on

# Disable Fabric
msm fabric set survival off

# Check status
msm fabric status survival

# Start server (auto-downloads Fabric launcher if enabled)
msm start survival

# List available versions
msm fabric versions 1.21.4
```

When Fabric is enabled, MSM prevents upgrading to Minecraft versions that Fabric doesn't support:

```bash
$ msm jar survival minecraft-26.1
Error: fabric does not yet support minecraft 26.1 - upgrade blocked
  Hint: Use --force to override
```

### Fabric Configuration

The `msm fabric set <server> on` command automatically updates the server's configuration. You can also manually edit `<server>/server.conf` to pin specific versions:

```bash
FABRIC_ENABLED="true"                # Set by 'msm fabric set <server> on'
FABRIC_LOADER_VERSION="0.16.10"      # Optional: pin specific loader version
FABRIC_INSTALLER_VERSION="1.1.0"     # Optional: pin specific installer version
```

## Full Command Reference

See [COMMANDS.md](COMMANDS.md) for the complete command reference, including all flags, environment variables, and exit codes.

## Shell Completion

MSM can generate shell completion scripts that enable Tab-completion for commands, subcommands, and server names. After installing, restart your shell or open a new terminal.

**Bash:**
```bash
sudo msm completion bash > /etc/bash_completion.d/msm
```

**Zsh:**
```bash
msm completion zsh > "${fpath[1]}/_msm"
```

**Fish:**
```bash
msm completion fish > ~/.config/fish/completions/msm.fish
```

Once installed, typing `msm ser<Tab>` completes to `msm server`, and `msm server li<Tab>` completes to `msm server list`.

## Upgrading from Bash MSM

See [UPGRADING.md](UPGRADING.md) for migration instructions from the original bash MSM.

The Go rewrite maintains full compatibility with existing configurations, server directories, and world storage, while adding new features:

- **Fabric mod loader** - Native support with version compatibility checks
- **Shell completion** - Auto-generated Tab-completion now includes Zsh and Fish
- **Safer RAM disk** - 2-minute background sync limits data loss on unexpected shutdown
- **Better visibility** - Tabular output showing status, ports, versions at a glance

## Troubleshooting

See [TROUBLESHOOTING.md](TROUBLESHOOTING.md) for common issues including:
- Permission denied errors
- Server stops immediately
- World not loading
- Fabric issues

Quick fix for most permission issues: `sudo msm setup`

## Development

```bash
# Run tests
make test

# Run linters
make lint

# Build for all platforms
make build-all
```

## License

MSM is released under the [GPLv3 license](LICENSE.md).

## Acknowledgements

This project grew out of [Ahtenus' Minecraft Init Script](https://github.com/Ahtenus/minecraft-init). The Go rewrite maintains compatibility with the original bash MSM configuration format.

## Support

- [Submit an issue](https://github.com/msmhq/msm/issues) for bugs or feature requests
- See the [changelog](CHANGELOG.md) for version history
