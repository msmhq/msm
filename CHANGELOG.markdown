Minecraft Server Manager
========================

Visit https://github.com/marcuswhybrow/minecraft-server-manager for more information.

Maintainers:
------------

* [Marcus Whybrow](https://github.com/marcuswhybrow)

Change Log
----------

### [0.4.2](https://github.com/marcuswhybrow/minecraft-server-manager/compare/0.4.1...0.4.2)

* Fixed the `msm <server> jar <jar-group> <jar-name>` command, which would state the `<jar-name>` was invalid, when it was not.


### [0.4.1](https://github.com/marcuswhybrow/minecraft-server-manager/compare/0.4.0...0.4.1)

* Fixed issue [#37][#37] where the `msm stop` command would hang. Thanks to merging pull request [#38][#38] from [@escortkeel][@escortkeel]. 
* Fixed issue [#39][#39] where server jar symlinks were not created. Thanks to merging pull request [#40][#40] from [@jives][@jives]. 

[@escortkeel]: https://github.com/escortkeel
[#37]: https://github.com/marcuswhybrow/minecraft-server-manager/issues/37
[#38]: https://github.com/marcuswhybrow/minecraft-server-manager/pull/38
[@jives]: https://github.com/jives
[#39]: https://github.com/marcuswhybrow/minecraft-server-manager/issues/39
[#40]: https://github.com/marcuswhybrow/minecraft-server-manager/pull/40

### [0.4.0](https://github.com/marcuswhybrow/minecraft-server-manager/compare/0.3.4...0.4.0)

* Added the `<server> give <player> <item> [amount] [data]` command. Closes [#34][#34]. Thanks to [@Vi3GameHkr][@Vi3GameHkr] for the suggestion.
* Added the `<server> xp <player> <amount>` command.
* Fixed the `gamemode`, `cmd` and `cmdlog` server commands *again*. They now accept multiple arguments as intented.

[@Vi3GameHkr]: https://github.com/Vi3GameHkr
[#34]: https://github.com/marcuswhybrow/minecraft-server-manager/issues/34

### [0.3.4](https://github.com/marcuswhybrow/minecraft-server-manager/compare/0.3.3...0.3.4)

* Fixed the `gamemode`, `cmd` and `cmdlog` server commands. They were not previously accessable. 

### [0.3.3](https://github.com/marcuswhybrow/minecraft-server-manager/compare/0.3.2...0.3.3)

* Fixed [#33][#33], by not assuming there are three invisible characters at the end of every log line.
* Fixed function names which were accidently converted to uppercase.

[#32]: https://github.com/marcuswhybrow/minecraft-server-manager/issues/32


### [0.3.2](https://github.com/marcuswhybrow/minecraft-server-manager/compare/0.3.1...0.3.2)

* Removes an erroneous debug statement. Fixes [#32][#32]

[#32]: https://github.com/marcuswhybrow/minecraft-server-manager/issues/32


### [0.3.1](https://github.com/marcuswhybrow/minecraft-server-manager/compare/0.3.0...0.3.1)

* Prevented `DEFAULT_SERVER_CONF` from being a required setting.


### [0.3.0](https://github.com/marcuswhybrow/minecraft-server-manager/compare/0.2.16...0.3.0)

* Closes [#22][#22], by merging `server.conf` into `server.properties`.
* Removed the `DEFAULT_SERVER_CONF` setting from `/etc/msm.conf`.

[#22]: https://github.com/marcuswhybrow/minecraft-server-manager/issues/22


### [0.2.16](https://github.com/marcuswhybrow/minecraft-server-manager/compare/0.2.15...0.2.16)

* Improved the `msm <server> worlds ram <world>` command to insist the server is stopped.
* Refactored global variables to be uppercase. This avoids name clashes.


### [0.2.15](https://github.com/marcuswhybrow/minecraft-server-manager/compare/0.2.14...0.2.15)

* Improved the removal of new line characters when reading Minecraft log files. Closes [#19][#19]

[#19]: https://github.com/marcuswhybrow/minecraft-server-manager/issues/19


### [0.2.14](https://github.com/marcuswhybrow/minecraft-server-manager/compare/0.2.13...0.2.14)

* Improved checks done before reading server files. Closes [#20][#20]. 

[#20]: https://github.com/marcuswhybrow/minecraft-server-manager/issues/20


### [0.2.13](https://github.com/marcuswhybrow/minecraft-server-manager/compare/0.2.12...0.2.13)

* Fixed the `msm <server> restart`, `msm <server> operator remove`, `msm <server> time set` and `msm <server> time add` commands.


### [0.2.12](https://github.com/marcuswhybrow/minecraft-server-manager/compare/0.2.11...0.2.12)

* Fixed `msm <server> stop` command. Which would fail for server's not first in the list.


### [0.2.11](https://github.com/marcuswhybrow/minecraft-server-manager/compare/0.2.10...0.2.11)

* Improved substring approach, which failed on some systems.


### [0.2.10](https://github.com/marcuswhybrow/minecraft-server-manager/compare/0.2.9...0.2.10)

* Fixed issue [#28][#28]. Which prevented the `msm jargroup delete` command working.
* Fixed a typo in the `msm version` help text.

[#28]: https://github.com/marcuswhybrow/minecraft-server-manager/issues/28


### [0.2.9](https://github.com/marcuswhybrow/minecraft-server-manager/compare/0.2.8...0.2.9)

* Fixed issue [#23][#23] which caused commands to fail registration on some systems.

[#23]: https://github.com/marcuswhybrow/minecraft-server-manager/issues/23


### [0.2.8](https://github.com/marcuswhybrow/minecraft-server-manager/compare/0.2.7...0.2.8)

* Improved command registration to provided more accurate error messages if a command cannot be registered.
* Improved Travis CI to only test the master branch.


### [0.2.7](https://github.com/marcuswhybrow/minecraft-server-manager/compare/0.2.6...0.2.7)

* Refactored commands such that they are now registered using a function, rather than one huge swith/case statement.


### [0.2.6](https://github.com/marcuswhybrow/minecraft-server-manager/compare/0.2.5...0.2.6)

* Added the "all" keyword which can be used to execute a command for all servers. Closes [#16][#16]
* Added [Travis CI integration][travis] for per commit unit testing.
* Improved command error messages to give specific command errors. Closes issue [#10][#10].
* Fixed issue [#15][#15] which left a "tail" command process running after certain commands.

[#15]: https://github.com/marcuswhybrow/minecraft-server-manager/issues/15
[#16]: https://github.com/marcuswhybrow/minecraft-server-manager/issues/16
[#10]: https://github.com/marcuswhybrow/minecraft-server-manager/issues/10

[travis]: http://travis-ci.org/#!/marcuswhybrow/minecraft-server-manager


### [0.2.5](https://github.com/marcuswhybrow/minecraft-server-manager/compare/0.2.4...0.2.5)

* Added `test.sh`, a script which runs `shunit2` unit tests.
* Added test coverage for all `msm server` commands, except coverage involving running servers.
* Added option to override the configuration file used by the MSM script by specifying its path in the `MSM_CONF` environment variable.
* Added Explanatory exit codes for the MSM script (currently 0 and 64-73).
* Improved server creation to not require the "minecraft" jar group to be present.
* Fixed issue [#14][#14] which prevented world names including spaces.
* Fixed a bug listing servers when there were none.
* Fixed a bug where server names containing reserved names anywhere where considered invalid.

[#11]: https://github.com/marcuswhybrow/minecraft-server-manager/issues/14


### [0.2.4](https://github.com/marcuswhybrow/minecraft-server-manager/compare/0.2.3...0.2.4)

* Fixes issue [#11][#11]. Now running server's cannot be renamed, preventing not being able to stop a server renamed whilst running.

[#11]: https://github.com/marcuswhybrow/minecraft-server-manager/issues/11


### [0.2.3](https://github.com/marcuswhybrow/minecraft-server-manager/compare/0.2.2...0.2.3)

* Fixed issues [#12][#12] and [#13][#13] thanks to [@Jonek01][@Jonek01], which prevented the first loaded server from having the correct path, and prevented the same for the first loaded world.

[#12]: https://github.com/marcuswhybrow/minecraft-server-manager/issues/12
[#13]: https://github.com/marcuswhybrow/minecraft-server-manager/issues/13
[@Jonek01]: https://github.com/Jonek01


### [0.2.2](https://github.com/marcuswhybrow/minecraft-server-manager/compare/0.2.1...0.2.2)

* Added `msm version` command, to check your installation version.


### [0.2.1](https://github.com/marcuswhybrow/minecraft-server-manager/compare/0.2.0...0.2.1)

* Added `msm <server> worlds on|off <world>` command, which activates and deactives worlds. Inactive worlds don't get backed up, cannot be used by the server, but are included in complete server backups.
* Added `CHANGELOG.markdown`.
* Added error messages preventing script use when a required `msm.conf` settings is not specified.
* Removed fallbacks for settigns not specified in `msm.conf`.


### [0.2.0](https://github.com/marcuswhybrow/minecraft-server-manager/compare/0.1.3...0.2.0)

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


### [0.1.1](https://github.com/marcuswhybrow/minecraft-server-manager/compare/0.1.0...0.1.1)

* Fixed a bug preventing world links from being created.
* Fixed a bug where a server stopped as the result of a command would not always be made inactive.


### [0.1.0](https://github.com/marcuswhybrow/minecraft-server-manager/compare/aaaa69dbd6c45ff0eb7dd69b756b4103abadfbc4...0.1.0)

* Refactored and improved protype code.
* Transitioned from managing a single to multiple servers.