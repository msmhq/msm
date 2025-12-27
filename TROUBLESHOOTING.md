# Troubleshooting

Common issues and solutions for MSM.

## Contents

- [Permission denied errors](#permission-denied-errors)
- [Server stops immediately after starting](#server-stops-immediately-after-starting)
- [World not loading (new random world generated)](#world-not-loading-new-random-world-generated)
- [Status shows different results for different users](#status-shows-different-results-for-different-users)
- [Fabric-related issues](#fabric-related-issues)
- [RAM disk issues](#ram-disk-issues)
- [Screen session issues](#screen-session-issues)
- [Getting more information](#getting-more-information)

## Permission denied errors

If you see "permission denied" errors when running MSM commands:

```bash
# Fix directory ownership and permissions
sudo msm setup
```

This ensures all MSM directories are owned by `minecraft:minecraft` with proper group-write permissions (2775).

### Common permission scenarios

| Error | Cause | Fix |
|-------|-------|-----|
| `failed to save fabric cache` | Root created cache files | `sudo msm setup` |
| `cannot create directory` | Directory owned by wrong user | `sudo msm setup` |
| `permission denied: server owned by user X` | Running as wrong user | Use `sudo` or switch to server owner |

## Server stops immediately after starting

Check the server logs:

```bash
cat /opt/msm/servers/<name>/logs/latest.log
```

### Common causes

**EULA not accepted**
```bash
sudo msm server init <name>
```
This auto-accepts the EULA and creates any missing config files.

**Java not installed**

| Minecraft Version | Minimum Java |
|-------------------|--------------|
| 1.21+ | Java 21 |
| 1.18 - 1.20.x | Java 17 |
| 1.17 | Java 16 |
| 1.12 - 1.16.x | Java 8 |

Install Java:
```bash
# Debian/Ubuntu
sudo apt install openjdk-21-jre-headless

# RHEL/Fedora
sudo dnf install java-21-openjdk-headless

# macOS
brew install openjdk@21
```

**Wrong Minecraft version**

Ensure the jar version matches your world. Older worlds may need older Minecraft versions.

## World not loading (new random world generated)

MSM stores worlds in `worldstorage/` and creates symlinks on start. This architecture is required for [RAM disk support](README.md#ram-disk-support)—see [IMPORTING.md](IMPORTING.md#how-msm-handles-worlds) for details.

### Checklist

1. **World folder exists:**
   ```bash
   ls /opt/msm/servers/<name>/worldstorage/
   ```

2. **Folder name matches `level-name`:**
   ```bash
   grep level-name /opt/msm/servers/<name>/server.properties
   ```
   The worldstorage folder must have the same name (default: `world`).

3. **Symlink created on start:**
   ```bash
   ls -la /opt/msm/servers/<name>/world
   ```
   Should show: `world -> /opt/msm/servers/<name>/worldstorage/world`

4. **No conflicting directory:**
   If a real `world/` directory exists (not a symlink) with data, MSM won't overwrite it. Move or delete it first.

## Status shows different results for different users

Screen sessions are user-specific. When running as root via `sudo`, MSM checks the configured user's screen sessions.

```bash
# These should show the same result
sudo msm status <server>
su - minecraft -c "msm status <server>"
```

If they differ, the server may have been started by a different user than configured.

## Fabric-related issues

### Fabric version not supported

```
Error: fabric does not yet support minecraft X.Y.Z
```

Fabric may not support the latest Minecraft snapshot/release yet. Check [fabricmc.net](https://fabricmc.net/) for supported versions.

### Fabric jar download fails

```bash
# Check network connectivity
curl -I https://meta.fabricmc.net/v2/versions/game

# Clear fabric cache and retry
sudo rm /opt/msm/fabric/cache.json
sudo msm setup
msm start <server>
```

## RAM disk issues

### Forbidden symlinks error

```
Found forbidden symlinks: /opt/msm/servers/survival/worldstorage/world -> /dev/shm/msm/survival/world
```

Minecraft 1.20+ blocks symlinks pointing outside the server directory by default. MSM automatically manages the `allowed_symlinks.txt` file when you enable RAM disk, but if you're seeing this error:

1. **Re-enable RAM** to regenerate the allowlist:
   ```bash
   msm worlds ram off <server> <world>
   msm worlds ram on <server> <world>
   ```

2. **Or manually create** `/opt/msm/servers/<name>/allowed_symlinks.txt`:
   ```
   prefix/dev/shm
   ```

## Screen session issues

### Can't attach to console

```bash
# Check if session exists
screen -ls

# Try attaching as the server owner
sudo -u minecraft screen -r msm-<servername>
```

### Session exists but server not running

The screen session may be orphaned. Kill it and restart:

```bash
msm stop <server> --now
screen -wipe
msm start <server>
```

## Getting more information

Enable verbose logging:

```bash
msm -v start <server>
msm --verbose status <server>
```

Check system logs:

```bash
journalctl -u msm  # If using systemd
```
