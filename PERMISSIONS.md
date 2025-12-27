# MSM Permissions Guide

This guide explains how MSM handles permissions and how to set up multi-user access to Minecraft servers.

## Contents

- [How It Works](#how-it-works)
- [Quick Reference](#quick-reference)
- [Default Setup (Single User)](#default-setup-single-user)
- [Multi-User Setup](#multi-user-setup)
- [Permission Errors](#permission-errors)
- [Directory Permissions Reference](#directory-permissions-reference)
- [Unix Permissions Primer](#unix-permissions-primer)
- [Troubleshooting](#troubleshooting)

## How It Works

MSM uses Unix user ownership to control who can manage each server:

- Each server has an **owner** (the `USERNAME` in its `server.conf`)
- You can only manage servers you own (or if you're root)
- The systemd service runs as root, so it can start all servers on boot

## Quick Reference

| Running as | Can manage own servers | Can manage others' servers |
|------------|----------------------|---------------------------|
| `root` (sudo) | Yes | Yes |
| Server owner | Yes | No |
| Other users | No | No |

## Default Setup (Single User)

Out of the box, MSM uses a single `minecraft` system user:

```bash
# All servers owned by minecraft user
ls -la /opt/msm/servers/
drwxr-xr-x minecraft minecraft survival/
drwxr-xr-x minecraft minecraft creative/

# Manage servers as root
sudo msm start survival

# Or as the minecraft user
sudo -u minecraft msm start survival
```

This is the simplest setup and works well for most installations.

## Multi-User Setup

Want multiple people to manage their own servers on shared hardware? Here's how:

### Step 1: Create the minecraft group

```bash
sudo groupadd minecraft
```

### Step 2: Set up shared directories

```bash
# Ensure directories exist with correct group
sudo chgrp -R minecraft /opt/msm
sudo chmod -R g+rwX /opt/msm
sudo chmod g+s /opt/msm/servers  # setgid for new directories

# Jar groups are shared (read-only for users)
sudo chmod 755 /opt/msm/jars
```

### Step 3: Add users to the group

```bash
# Add existing user to minecraft group
sudo usermod -aG minecraft alice
sudo usermod -aG minecraft bob

# User must log out and back in for group to take effect
```

### Step 4: Users create their servers

```bash
# Alice creates her server (automatically owned by alice)
alice$ msm server create alice-survival
Created server "alice-survival" at /opt/msm/servers/alice-survival

# Server is owned by alice
ls -la /opt/msm/servers/alice-survival/
drwxr-xr-x alice minecraft .
-rw-r--r-- alice minecraft server.conf
```

### Step 5: Users manage their own servers

```bash
# Alice can manage her server
alice$ msm start alice-survival
Started server "alice-survival"

# Alice cannot manage Bob's server
alice$ msm stop bob-creative
Error: permission denied: server "bob-creative" is owned by user "bob"
  You are running as "alice"
  Hint: Run with sudo, or as user "bob"

# Bob manages his own
bob$ msm start bob-creative
Started server "bob-creative"
```

### Step 6: Systemd starts everything on boot

The systemd service runs as root, so it starts all servers regardless of ownership:

```bash
sudo systemctl enable msm
# On boot: starts alice-survival, bob-creative, and any other servers
```

## Permission Errors

### "permission denied: server X is owned by user Y"

You're trying to manage a server you don't own.

**Solutions:**
1. Run as root: `sudo msm <command>`
2. Run as the owner: `sudo -u <owner> msm <command>`
3. If you should own it, have an admin change ownership

### "failed to set ownership"

When creating a server, MSM couldn't set the directory ownership.

**Causes:**
- Running as non-root user who doesn't exist in `/etc/passwd`
- Target user doesn't exist

**Solution:** Create servers as root or ensure your user account exists.

### Server won't start on boot

The systemd service might not be enabled or the server config has issues.

**Check:**
```bash
sudo systemctl status msm
sudo journalctl -u msm
```

## Directory Permissions Reference

```
/opt/msm/                    root:minecraft 755
├── servers/                 root:minecraft 2775 (setgid)
│   ├── alice-survival/      alice:minecraft 755
│   └── bob-creative/        bob:minecraft 755
├── jars/                    root:minecraft 755 (shared, read-only)
├── archives/                root:minecraft 2775 (setgid)
│   ├── worlds/
│   ├── logs/
│   └── backups/
└── fabric/                  root:minecraft 755 (shared cache)
```

The `setgid` bit (2xxx) on directories ensures new files inherit the `minecraft` group.

## Unix Permissions Primer

If you're new to Unix permissions, here's a quick overview:

### Users and Groups

- Every file has an **owner** (a user) and a **group**
- Users can belong to multiple groups
- Check your groups: `groups` or `id`

### Permission Bits

```
drwxr-xr-x alice minecraft servers/
│├─┤├─┤├─┤
│ │  │  └── others: r-x (read + execute)
│ │  └───── group:  r-x (read + execute)
│ └──────── owner:  rwx (read + write + execute)
└────────── d = directory
```

- `r` (read): List directory contents / read file
- `w` (write): Create/delete files in directory / modify file
- `x` (execute): Enter directory / run file

### Common Commands

```bash
# View permissions
ls -la /opt/msm/servers/

# Change owner
sudo chown alice:minecraft /opt/msm/servers/alice-survival

# Change permissions
sudo chmod 755 /opt/msm/servers/alice-survival

# Add user to group
sudo usermod -aG minecraft alice
```

## Troubleshooting

### Check who owns a server

```bash
ls -la /opt/msm/servers/myserver/server.conf
# Or check the config
grep USERNAME /opt/msm/servers/myserver/server.conf
```

### Check what user you're running as

```bash
whoami
id
```

### Check group membership

```bash
groups              # Your groups
groups alice        # Alice's groups
getent group minecraft  # Who's in minecraft group
```

### Test permissions

```bash
# Can you read the server directory?
ls /opt/msm/servers/myserver/

# Can you write to it?
touch /opt/msm/servers/myserver/test && rm /opt/msm/servers/myserver/test
```
