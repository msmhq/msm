Minecraft Server Manager
========================

Visit https://github.com/marcuswhybrow/minecraft-server-manager for more information.

Maintainers:
------------

* [Marcus Whybrow](https://github.com/marcuswybrow)

Change Log
----------

### [0.2.2](https://github.com/marcuswhybrow/minecraft-server-manager/compare/0.2.1...0.2.2)

* Added `msm version` command, to check your installation version.

### [0.2.1](https://github.com/marcuswhybrow/minecraft-server-manager/compare/0.2...0.2.1)

* Added `msm <server> worlds on|off <world>` command, which activates and deactives worlds. Inactive worlds don't get backed up, cannot be used by the server, but are included in complete server backups.
* Added `CHANGELOG.markdown`.
* Added error messages preventing script use when a required `msm.conf` settings is not specified.
* Removed fallbacks for settigns not specified in `msm.conf`.

### [0.2](https://github.com/marcuswhybrow/minecraft-server-manager/compare/0.1.3...0.2)

* Added bash completion.
* Altered project layout to feature directories for bash_completion, cron and init files. This is useful for dealing with three files all named `msm`.
* Removed `msm <server> worlds toram` as an accessible command.
* Fixed a bug preventing `WORLD_STORAGE_PATH` from being overriden by a server.
* Fixed error in `msm.conf`

### [0.1.3](https://github.com/marcuswhybrow/minecraft-server-manager/compare/0.1.2...0.1.3)

* Added error output for the case when creating a symbolic link for a server's world it is found that a file already exists with that name.
* Added `msm <server> operator list` command.
* Added help output for the `msm <server> jar` command.
* Improved world backups to be [WorldEdit snapshot](wiki.sk89q.com/wiki/WorldEdit/Snapshots) compatible.
* Altered "complete server backups" to zip the directory instead of the contents.
* Removed requiring a server to be running to list whitelisted and blacklisted (banned) players.
* Fixed a bug preventing worlds from getting backed up.
* Fixed a bug preventing actions performed upon server worlds.
* Fixed a bug preventing all server worlds from being moved into or back to disk RAM.
* Fixed a bug preventing server logs from being rolled correctly.
* Fixed a bug preventing `msm <server> connected` from working.
* Fixed a bug preventing `msm <server> jar` from working.

### [0.1.2](https://github.com/marcuswhybrow/minecraft-server-manager/compare/0.1.1...0.1.2)

* Added in-game messages, and correct pausing and resuming of Minecraft world saving around world backups.
* Added "Backup took n seconds." message to world and complete backups.
* Fixed a bug listing a server's worlds.
* Fixed a bug when stopping all servers in parallel, where a server might not stop at the correct time.

### [0.1.1](https://github.com/marcuswhybrow/minecraft-server-manager/compare/0.1...0.1.1)

* Fixed a bug preventing world links from being created.
* Fixed a bug where a server stopped as the result of a command would not always be made inactive.

### [0.1](https://github.com/marcuswhybrow/minecraft-server-manager/compare/aaaa69dbd6c45ff0eb7dd69b756b4103abadfbc4...0.1)

* Refactored and improved protype code.
* Transitioned from managing a single to multiple servers.