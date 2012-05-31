# Minecraft Server Manager

A single init script which makes running multiple Minecraft/Bukkit servers easier for us admins.

[Visit the Wiki][wiki] for detailed guides.

## Features

* One script handles multiple servers, run two or more servers on one machine.
* Can create and start new servers with a single command, downloads the jars for you.
* Start, stop and restart your servers.
* Periodically makes [WorldEdit snapshot][we-snapshot] compatible backups of your worlds.
* Backup the entire server folder for complete protection.
* Load world's into RAM for faster access (reduces lag).
* Easily configurable, through a single file
* Per server configuration of invocation, in-game messages and timings.
* Apply server commands to one, multiple, or all servers in one go (useful for whitelisting a player on all servers.)

## Upcomming features

* **Restore:** Roll-back to an old world or whole server backup automatically.
* **QuickBackup:** If you store your backups non-locally (maybe on a NAS), QuickBackup optionally creates a backup locally for speed, and then moves it after your players are building again! My initial testing shows a 54 second network backup confaltes to 23 seconds of ingame time.

## Acknowledgements

This code grew out of an old version of [Ahtenus' Minecraft Init Script][ahtenus-minecraft-init].

[we-snapshot]: http://wiki.sk89q.com/wiki/WorldEdit/Snapshots
[ahtenus-minecraft-init]: https://github.com/Ahtenus/minecraft-init
[wiki]: https://github.com/marcuswhybrow/minecraft-server-manager/wiki