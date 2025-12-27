# Upgrading from Bash MSM

The Go rewrite of MSM maintains full compatibility with existing MSM configurations.

## Contents

- [What stays the same](#what-stays-the-same)
- [What's new](#whats-new)
- [Migration steps](#migration-steps)
- [Terminology changes](#terminology-changes)
- [New commands](#new-commands)
- [Systemd integration](#systemd-integration)
- [Troubleshooting migration](#troubleshooting-migration)

## What stays the same

- Configuration file format (`/etc/msm.conf`)
- Server directory structure
- Per-server `server.conf` files
- World storage layout
- JSON player files (whitelist, ops, bans)

## What's new

- Single static binary (no bash dependencies)
- Structured logging with `-v/--verbose` flag
- Faster startup and execution
- Automatic RAM sync daemon (no cron required)
- `msm setup` command for permission management

## Migration steps

### 1. Stop all servers

```bash
msm stop --now
```

### 2. Backup your configuration

```bash
cp /etc/msm.conf /etc/msm.conf.backup
```

### 3. Clone and build

```bash
git clone https://github.com/msmhq/msm.git
cd msm
make build
```

### 4. Remove old bash MSM

```bash
sudo make migrate
```

This removes:
- `/etc/init.d/msm` (old bash script)
- Init.d symlinks in `/etc/rc*.d/`
- Old cron entries (replaced by `msm cron install`)

### 5. Install

```bash
sudo make setup    # Creates user and directories
sudo make install  # Installs binary, config, and cron
```

### 6. Verify and start

```bash
msm config
msm server list
msm start
```

## Terminology changes

| Old (bash MSM) | New (Go MSM) |
|----------------|--------------|
| `msm whitelist` | `msm allowlist` |
| `msm blacklist` | `msm blocklist` |

The underlying Minecraft files (`whitelist.json`, `banned-players.json`) remain unchanged.

## New commands

| Command | Description |
|---------|-------------|
| `msm setup` | Fix directory ownership and permissions |
| `msm server init <name>` | Initialize config for imported worlds |
| `msm jar download <server> [version]` | Download vanilla Minecraft jar |
| `msm fabric on/off <server>` | Enable/disable Fabric mod loader |

## Systemd integration

The new MSM includes optional systemd service files:

```bash
sudo make systemd-install
sudo systemctl enable msm
```

This replaces the old init.d scripts with proper systemd units.

## Troubleshooting migration

### Servers not starting after upgrade

```bash
# Check permissions
sudo msm setup

# Verify server config
msm server config <name>
```

### Screen sessions from old MSM

Old screen sessions may still exist:

```bash
screen -ls
screen -wipe  # Clean up dead sessions
```

### Cron jobs

The new MSM manages its own cron job:

```bash
msm cron install   # Install/update cron job
msm cron remove    # Remove cron job
```

Check `/etc/cron.d/msm` for the installed cron configuration.
