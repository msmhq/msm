# Importing Existing Worlds

If you have an existing Minecraft world (from single-player or another server), you can import it into MSM.

## Contents

- [From Single-Player](#from-single-player)
- [From Another Server](#from-another-server)
- [Custom level-name](#custom-level-name)
- [Modpacks](#modpacks)
- [World Folder Structure](#world-folder-structure)
- [How MSM Handles Worlds](#how-msm-handles-worlds)

## From Single-Player

Single-player worlds are stored in your Minecraft saves folder:
- **Linux**: `~/.minecraft/saves/<world_name>/`
- **macOS**: `~/Library/Application Support/minecraft/saves/<world_name>/`
- **Windows**: `%APPDATA%\.minecraft\saves\<world_name>\`

To import:

```bash
# 1. Create the server
msm server create survival

# 2. Copy your world into the server's worldstorage directory
#    Target folder name must match level-name in server.properties (default: "world")
sudo rsync -r --chown=minecraft:minecraft ~/.minecraft/saves/MyWorld/ \
    /opt/msm/servers/survival/worldstorage/world/

# 3. Download Minecraft and start
msm jar download survival
msm start survival
```

## From Another Server

If migrating from an existing Minecraft server:

```bash
# 1. Create the MSM server directory
sudo mkdir -p /opt/msm/servers/migrated

# 2. Copy server files (server.properties, whitelist, ops, etc.)
sudo rsync -r --chown=minecraft:minecraft /path/to/old-server/ \
    /opt/msm/servers/migrated/ \
    --exclude='world' --exclude='logs'

# 3. Copy the world folder to worldstorage
#    Use the same target folder name as level-name in server.properties
sudo rsync -r --chown=minecraft:minecraft /path/to/old-server/world/ \
    /opt/msm/servers/migrated/worldstorage/world/

# 4. Initialize missing config and accept EULA
sudo msm server init migrated

# 5. Download Minecraft (or use existing jar) and start
msm jar download migrated
msm start migrated
```

## Custom level-name

If your old server used a custom `level-name` (e.g., `level-name=MyWorld`), name the worldstorage folder to match:

```
/opt/msm/servers/migrated/worldstorage/MyWorld/
```

MSM automatically creates a symlink from the server directory to worldstorage on start.

## Modpacks

Some modpacks include their own server jar. If your world folder contains a `.jar` file, configure `JAR_PATH` in `server.conf` to use it:

```bash
# Edit server.conf
echo 'JAR_PATH="worldstorage/world/server.jar"' >> /opt/msm/servers/migrated/server.conf
```

## World Folder Structure

A valid Minecraft world folder contains:

```
world/
├── level.dat          # World metadata (required)
├── region/            # Overworld chunk data
├── DIM-1/             # Nether (if visited)
│   └── region/
├── DIM1/              # The End (if visited)
│   └── region/
├── data/              # Map data, raids, etc.
└── playerdata/        # Player inventories
```

The folder name must match `level-name` in `server.properties` (default: `world`).

## How MSM Handles Worlds

MSM uses a `worldstorage/` subdirectory to organize worlds:

```
/opt/msm/servers/survival/
├── server.conf
├── server.properties
├── worldstorage/
│   └── world/           # Your actual world data
└── world -> worldstorage/world  # Symlink created on start
```

When you run `msm start`:
1. MSM reads `level-name` from `server.properties`
2. Checks if that folder exists in `worldstorage/`
3. Creates a symlink from the server directory to worldstorage
4. Minecraft follows the symlink and loads your world

### Why this architecture?

**This design is required for RAM disk support.** When you enable RAM for a world (`msm worlds ram on <server> <world>`), MSM copies the world to `/dev/shm` and symlinks to it there. This dramatically reduces I/O latency for busy servers.

The `worldstorage/` indirection makes this possible—without it, there would be no clean way to swap between disk and RAM storage.

Additional benefits:
- Multiple worlds per server (swap by changing `level-name`)
- Clean separation of config and world data
