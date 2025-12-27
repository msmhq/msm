# MSM Command Reference

Complete reference for all MSM commands.

## Contents

- [Global Commands](#global-commands)
- [Server Management](#server-management)
- [World Management](#world-management)
- [Backups](#backups)
- [Player Management](#player-management)
- [Jar Group Management](#jar-group-management)
- [Fabric Mod Loader](#fabric-mod-loader)
- [Global Flags](#global-flags)
- [Environment Variables](#environment-variables)
- [Exit Codes](#exit-codes)

## Global Commands

Commands that operate on all servers or the MSM system itself.

| Command | Description |
|---------|-------------|
| `msm start` | Start all servers |
| `msm stop [--now]` | Stop all servers (--now skips player warning) |
| `msm restart [--now]` | Restart all servers |
| `msm config` | Display global configuration |
| `msm version` | Print version number |
| `msm cron generate` | Output cron file to stdout |
| `msm cron install` | Install cron file to /etc/cron.d/msm |
| `msm logroll [--all]` | Archive and compress server logs |

## Server Management

### Listing and Creating

| Command | Description |
|---------|-------------|
| `msm server list` | List all servers with status |
| `msm server create <name>` | Create a new server |
| `msm server delete <name>` | Delete a server (must be stopped) |
| `msm server rename <old> <new>` | Rename a server |
| `msm server init <name>` | Initialize missing config files |
| `msm server config <server>` | Show all config values |
| `msm server config <server> <key>` | Show specific config value |
| `msm server config set <server> <key> <value>` | Set a config value |

### Lifecycle

| Command | Description |
|---------|-------------|
| `msm start <server>` | Start a server |
| `msm stop <server> [--now]` | Stop a server |
| `msm restart <server> [--now]` | Restart a server |
| `msm status <server>` | Show server status |

### Console and Commands

| Command | Description |
|---------|-------------|
| `msm console <server>` | Attach to server console (Ctrl+A, D to detach) |
| `msm cmd <server> <command>` | Send command to server console |
| `msm say <server> <message>` | Broadcast message to all players |
| `msm kick <server> <player> [reason]` | Kick a player |
| `msm connected <server>` | List connected players |

## World Management

| Command | Description |
|---------|-------------|
| `msm worlds list <server>` | List all worlds for a server |
| `msm worlds set <server> <world> <on\|off>` | Activate or deactivate a world |
| `msm worlds ram <server> <world>` | Show RAM disk status for a world |
| `msm worlds ram set <server> <world> <on\|off>` | Enable or disable RAM disk for a world |
| `msm worlds todisk [server]` | Sync RAM worlds to disk |
| `msm worlds todisk --all` | Sync all RAM worlds (all servers) |
| `msm worlds backup [server]` | Backup all worlds for a server |
| `msm worlds backup --all` | Backup all worlds (all servers) |

## Backups

| Command | Description |
|---------|-------------|
| `msm worlds backup <server>` | Backup all worlds for a server |
| `msm worlds backup --all` | Backup all worlds (all servers) |
| `msm backup <server>` | Full server backup (worlds + configs) |

## Player Management

### Operators

| Command | Description |
|---------|-------------|
| `msm op add <server> <player>` | Make player an operator |
| `msm op remove <server> <player>` | Remove operator status |
| `msm op list <server>` | List all operators |

### Allowlist (Whitelist)

| Command | Description |
|---------|-------------|
| `msm allowlist set <server> <on\|off>` | Enable or disable allowlist enforcement |
| `msm allowlist add <server> <player>...` | Add player(s) to allowlist |
| `msm allowlist remove <server> <player>...` | Remove player(s) from allowlist |
| `msm allowlist list <server>` | List all allowlisted players |

### Blocklist (Bans)

| Command | Description |
|---------|-------------|
| `msm blocklist player add <server> <player>... [--reason "..."]` | Ban player(s) |
| `msm blocklist player remove <server> <player>...` | Unban player(s) |
| `msm blocklist player list <server>` | List banned players |
| `msm blocklist ip add <server> <ip>... [--reason "..."]` | Ban IP address(es) |
| `msm blocklist ip remove <server> <ip>...` | Unban IP address(es) |
| `msm blocklist ip list <server>` | List banned IPs |

## Jar Group Management

Jar groups organize server JAR files for easy version management.

| Command | Description |
|---------|-------------|
| `msm jargroup list` | List all jar groups |
| `msm jargroup create <name> <url>` | Create a jar group with download URL |
| `msm jargroup delete <name>` | Delete a jar group |
| `msm jargroup rename <old> <new>` | Rename a jar group |
| `msm jargroup changeurl <name> <url>` | Change the download URL |
| `msm jargroup getlatest <name>` | Download latest JAR from URL |

### Server Jar Management

| Command | Description |
|---------|-------------|
| `msm jar <server>` | Show current jar for server |
| `msm jar link <server> <jargroup>` | Link server to latest jar in group |
| `msm jar link <server> <jargroup> <file>` | Link server to specific jar file |
| `msm jar link <server> <jargroup> --force` | Force link (bypasses Fabric check) |
| `msm jar download <server>` | Download latest vanilla Minecraft |
| `msm jar download <server> <version>` | Download specific Minecraft version |

## Fabric Mod Loader

Commands for managing Fabric mod loader integration.

| Command | Description |
|---------|-------------|
| `msm fabric status <server>` | Show Fabric status and versions |
| `msm fabric set <server> <on\|off>` | Enable or disable Fabric for server |
| `msm fabric versions` | List supported Minecraft versions |
| `msm fabric versions <mc-version>` | List loader versions for MC version |
| `msm fabric update <server>` | Check for and download newer loader |
| `msm fabric set-loader <server> <version>` | Pin specific loader version |
| `msm fabric set-installer <server> <version>` | Pin specific installer version |

## Global Flags

These flags work with any command:

| Flag | Description |
|------|-------------|
| `--config <path>` | Use alternate config file |
| `-v, --verbose` | Enable debug logging |
| `-h, --help` | Show help for command |

## Environment Variables

| Variable | Description |
|----------|-------------|
| `MSM_CONF` | Path to config file (default: /etc/msm.conf) |

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 64 | Invalid user |
| 65 | Invalid command |
| 66 | Invalid argument |
| 67 | Server stopped (when running command expected) |
| 68 | Server running (when stopped expected) |
| 69 | Name not found |
| 70 | File not found |
| 71 | Duplicate name |
