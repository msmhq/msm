Minecraft Server Manager
========================

Visit https://github.com/msmhq/msm for more information.

Change Log
----------

### [0.9.4](https://github.com/msmhq/msm/compare/0.9.4...0.10.0)

First minor version in 4 years.

* Allow snapshots to be downloaded [#369][#369]
* Add dollar sign for interpolation [#389][#389]
* Add new instantiated Systemd unit file [#352][#352]
* Add configurable archive cleanup [#394][#394]
* Fix bash variable expansion for Debian installer [#401][#401]
* Remove invalid java invocation arguments [#403][#403]
* Convert server txt files to json [#405][#405]

### [0.9.4](https://github.com/msmhq/msm/compare/0.9.3...0.9.4)

* Added versioning support for the following properties. [#310][#310]
  LOG_PATH;WHITELIST_PATH;BANNED_PLAYERS_PATH;BANNED_IPS_PATH;OPS_PATH;OPS_LIST;
* Moved LOG_PATH to the version files.

[#310]: https://github.com/msmhq/msm/pull/310

### [0.9.3](https://github.com/msmhq/msm/compare/0.9.2...0.9.3)

* Added pre and post 1.7.0 log support [#308][#308]

[#308]: https://github.com/msmhq/msm/pull/308

### [0.9.2](https://github.com/msmhq/msm/compare/0.9.1...0.9.2)

* Added check and warning if eula isn't set to true [#271][#271]

[#271]: https://github.com/msmhq/msm/issues/271

### [0.9.1](https://github.com/msmhq/msm/compare/0.9.0...0.9.1)

* Fixed version target and default jargroup for backwards compatibility [#298][#298]

[#298]: https://github.com/msmhq/msm/pull/298

### [0.9.0](https://github.com/msmhq/msm/compare/0.8.16...0.9.0)

* Project has been moved to it's own organization
* Added "minecraft" default jar group on MSM setup. [a39e19][a39e19]
* Added downloading latest Minecraft versions (via versions.json) [#169][#169]
* Added default DEFAULT_OPS_LIST to msm.conf [#195][#195]
* Added release and snapshot download support for minecraft jars [#214][#214]
* Added Arch Linux installer (including systemd support) [#217][#217]
* Added testing using default shunit2 location if command not found [#218][#218]
* Added RDIFF backup support [#221][#221]
* Added RSync backup support [#243][#243]
* Added unified distro-detecting installer [#227][#227]
* Added option to create server user as "system" on MSM install [#232][#232]
* Fixed hidden world storage path error message [#282][#282]
* Fixed update versioning comparisons. [#132][#132]
* Fixed annoying "No directory, logging in with HOME=/" message [#182][#182]
* Fixed deleting running servers [#185][#185]
* Fixed broken error messages when deactivating worlds and calling commands [#188][#188]
* Fixed new worlds created outside of world storage directory [#191][#191]
* Fixed directory nesting when toggling world active state [#204][#204]
* Fixed `msm restart` broadcasting wrong message to players [#196][#196]
* Fixed RedHat installing `/etc/cron.d/msm` with incorrect permissions [#203][#203]
* Fixed RedHat installer to require java package [#236][#236]
* Fixed cron not reloading after installation [#245][#245]
* Fixed not following redirects in install scripts [#267][#267]
* Fixed minor message and comment typos [#192][#192], [#226][#226], [#246][#246]
* Refactored install scripts allowing alternate github repos [#193][#193]

[#132]: https://github.com/msmhq/msm/pull/132
[#169]: https://github.com/msmhq/msm/pull/169
[#182]: https://github.com/msmhq/msm/pull/182
[#185]: https://github.com/msmhq/msm/pull/185
[#186]: https://github.com/msmhq/msm/pull/186
[#188]: https://github.com/msmhq/msm/pull/188
[#191]: https://github.com/msmhq/msm/pull/191
[#192]: https://github.com/msmhq/msm/pull/192
[#193]: https://github.com/msmhq/msm/pull/193
[#195]: https://github.com/msmhq/msm/pull/195
[#196]: https://github.com/msmhq/msm/pull/196
[#203]: https://github.com/msmhq/msm/pull/203
[#204]: https://github.com/msmhq/msm/pull/204
[#214]: https://github.com/msmhq/msm/pull/214
[#217]: https://github.com/msmhq/msm/pull/217
[#218]: https://github.com/msmhq/msm/pull/218
[#221]: https://github.com/msmhq/msm/pull/221
[#226]: https://github.com/msmhq/msm/pull/226
[#227]: https://github.com/msmhq/msm/pull/227
[#232]: https://github.com/msmhq/msm/pull/232
[#236]: https://github.com/msmhq/msm/pull/236
[#243]: https://github.com/msmhq/msm/pull/243
[#245]: https://github.com/msmhq/msm/pull/245
[#246]: https://github.com/msmhq/msm/pull/246
[#267]: https://github.com/msmhq/msm/pull/267
[#282]: https://github.com/msmhq/msm/issues/282
[a39e19]: https://github.com/msmhq/msm/commit/a39e1948abbe01190d6e9b5472676713557e1910

### [0.8.16](https://github.com/msmhq/msm/compare/0.8.15...0.8.16)

* Fix various typos in MSM and the documentation


### [0.8.15](https://github.com/msmhq/msm/compare/0.8.14...0.8.15)

* Fix incorrect reference to `$latest_minecraft_version`


### [0.8.14](https://github.com/msmhq/msm/compare/0.8.13...0.8.14)

* Merged pull request [#110][#110], which fixes a typo in the cron file
* Made RAMDisk enabled by default in conf file

[#110]: https://github.com/msmhq/msm/pull/110


### [0.8.13](https://github.com/msmhq/msm/compare/0.8.12...0.8.13)

* Fixed issue [#105][#105], which showed it was impossible to disable RAMDisk as suggested by MSM.

[#105]: https://github.com/msmhq/msm/issues/105


### [0.8.12](https://github.com/msmhq/msm/compare/0.8.11...0.8.12)

* Fixed issue [#96][#96], which noted that the `msm jargroup changeurl` command was not implemented.

[#96]: https://github.com/msmhq/msm/issues/96


### [0.8.11](https://github.com/msmhq/msm/compare/0.8.10...0.8.11)

* Fixed issue [#76][#76], where tab completing "help" or "update" would exit the terminal.
* Made "help" and "update" reserved names to prevent clashing with server names.

[#76]: https://github.com/msmhq/msm/issues/76


### [0.8.10](https://github.com/msmhq/msm/compare/0.8.9...0.8.10)

* Improved versioning to use the latest Minecraft version by default, and tell the user what's happening and how to control it themselves. Fixes issue [#79][#79].
* Improved install script to also create the default minecraft user, and the default install location directories.

[#79]: https://github.com/msmhq/msm/issues/79


### [0.8.9](https://github.com/msmhq/msm/compare/0.8.8...0.8.9)

* Improved install script to update MSM without input


### [0.8.8](https://github.com/msmhq/msm/compare/0.8.7...0.8.8)

* Added an MSM install script.
* Fixed bug in update code which prevent cleaning download directory.


### [0.8.7](https://github.com/msmhq/msm/compare/0.8.6...0.8.7)

* Improved upon 0.8.6 to to prevent MSM complaing when not being invoked by either the "root" or MSM user.


### [0.8.6](https://github.com/msmhq/msm/compare/0.8.5...0.8.6)

* Prevent script complaing that `This command must be executed as the user "root"` after every command. Fixes [#77][#77] and [#78][#78].

[#77]: https://github.com/msmhq/msm/issues/77
[#78]: https://github.com/msmhq/msm/issues/78


### [0.8.5](https://github.com/msmhq/msm/compare/0.8.4...0.8.5)

* Fixed `msm update` output showing the "created files" section only when files were being *updated*.


### [0.8.4](https://github.com/msmhq/msm/compare/0.8.3...0.8.4)

* Fixed logic in `msm update` which allowed MSM to think everything was already updated.


### [0.8.3](https://github.com/msmhq/msm/compare/0.8.2...0.8.3)

* Fixed a bug in `msm update` checks, which caused a `command not found` error.


### [0.8.2](https://github.com/msmhq/msm/compare/0.8.1...0.8.2)

* Improved `msm update` output to make more sense.


### [0.8.1](https://github.com/msmhq/msm/compare/0.8.0...0.8.1)

* Added bash completion for `msm update`
* Improved `msm update` to work when upgrading from older versions, and handle updating the update code.


### [0.8.0](https://github.com/msmhq/msm/compare/0.7.5...0.8.0)

* Added the `msm update` command, which downloads and installs new versions of MSM.
* Added the `msm-version` server config setting. Add this to `server.properties` to allow MSM to interact correctly with your server (e.g. `msm-version=minecraft/1.3.1`.)
* Removed all confirmation settings, such as `DEFAULT_CONFIRM_SAVE_ON` and `msm-confirm-save-all` etc, this functionality is handled by new [versioning files][versioning-files].
* Fixed issues [#43][#43], [#58][#58], [#59][#59], [#60][#60], [#61][#61], [#63][#63], [#67][#67] and [#68][#68]. By accomodating different Minecraft versions using versioning files.

[#43]: https://github.com/msmhq/msm/issues/43
[#58]: https://github.com/msmhq/msm/issues/58
[#59]: https://github.com/msmhq/msm/issues/59
[#60]: https://github.com/msmhq/msm/issues/60
[#61]: https://github.com/msmhq/msm/issues/61
[#63]: https://github.com/msmhq/msm/issues/63
[#67]: https://github.com/msmhq/msm/issues/67
[#68]: https://github.com/msmhq/msm/issues/68
[versioning-files]: https://github.com/msmhq/msm/tree/master/versioning

### [0.7.5](https://github.com/msmhq/msm/compare/0.7.4...0.7.5)

* Fixed issue [#69][#69] where pressing tab after `msm config` or `msm start` would disconnect the terminal session.

[#69]: https://github.com/msmhq/msm/issues/69


### [0.7.4](https://github.com/msmhq/msm/compare/0.7.3...0.7.4)

* Fixed issue [#64][#64] where using server name `all` and tab completing caused the terminal to exit.

[#64]: https://github.com/msmhq/msm/issues/64


### [0.7.3](https://github.com/msmhq/msm/compare/0.7.2...0.7.3)

* Fixed global stop and start commands from failing since [0.7.0][0.7.0].

[0.7.0]: https://github.com/msmhq/msm/compare/0.6.4...0.7.0


### [0.7.2](https://github.com/msmhq/msm/compare/0.7.1...0.7.2)

* Fixed several bugs where the server username was not being used. Affected the `msm <server> worlds backup`, `msm <server> cmd` and `msm <server> cmdlog` commands.


### [0.7.1](https://github.com/msmhq/msm/compare/0.7.0...0.7.1)

* Fixed issue [#54][#54] where setting a server's jar file would fail to create the symbolic link.

[#54]: https://github.com/msmhq/msm/issues/54


### [0.7.0](https://github.com/msmhq/msm/compare/0.6.4...0.7.0)

* Added `/etc/profile.d/msm.sh` support. Define environment variables there, and both MSM and bash completion scripts will see it. At this point you can add the `MSM_SCRIPT` and `MSM_CONF` environment variables to tell MSM where you have installed your MSM files:
	* Add `export MSM_SCRIPT="/path/to/msm"` if you have changed it from `/etc/init.d/msm`.
	* Add `export MSM_CONF="/path/to/msm.conf"` if you have changed it from `/etc/msm.conf`.
* Added some simple stopped server unit tests for half the commands.
* Refactored config property loading to be "lazy". Speeds up all commands by only loading from file what is needed.
* Improved whitelist, operator, player ban, ip ban, gamemode and kick commands to support multiple space-separated users.
* Improved whitelist, operator, player ban and ip ban commands to also work when the server is stopped.
* Improved MSM script speed for BASH 4 users (loading config properties is now much faster.)
* Improved BASH completion to also use default settings specified in MSM.
* Removed use of subshells in script (using less means faster code.)


### [0.6.4](https://github.com/msmhq/msm/compare/0.6.3...0.6.4)

* Closes issue [#25][#25] by improving the `msm <server> start` command to be more robust. It now alerts the user when the server jar to be started cannot be found, and prints out dots for each log line created during startup. The absence of these dots make it simple to spot a hanging server.
* Fixed a bug with the `msm <server> say` command, where it would only accept 2 or more arguments.

[#25]: https://github.com/msmhq/msm/issues/25


### [0.6.3](https://github.com/msmhq/msm/compare/0.6.2...0.6.3)

* Improved the `msm server list` command to show active/inactive and running/stopped statuses. Thanks to [oldman,j's suggestion][uservoice-2987515].

[uservoice-2987515]: http://minecraftservermanager.uservoice.com/forums/166773-general/suggestions/2987515-the-server-list-command-


### [0.6.2](https://github.com/msmhq/msm/compare/0.6.1...0.6.2)

* Fixed issue [#49][#49] where, probably since 0.6.0, `active` files where not being created upon server start.

[#49]: https://github.com/msmhq/msm/issues/49


### [0.6.1](https://github.com/msmhq/msm/compare/0.6.0...0.6.1)

* Fixed `test.sh` to use new settings names.


### [0.6.0](https://github.com/msmhq/msm/compare/0.5.2...0.6.0)

* Added new command `msm config` which displays a key value list of all MSM settings in use. Useful for transparency.
* Improved `msm <server> config <setting> <value>` to list all MSM settings for a server when arguments are omitted.
* Updated bash completion script for new commands.
* Refactored the way settings are added globally and for individual servers, making development far simpler.
* Changed several setting names in `/etc/msm.conf`:
	* `DEFAULT_SERVER_USER` is now `DEFAULT_USERNAME`
	* `DEFAULT_LOG` is now `DEFAULT_LOG_PATH`
	* `DEFAULT_PROPERTIES` is now `DEFAULT_PROPERTIES_PATH`
	* `DEFAULT_WHITELIST` is now `DEFAULT_WHITELIST_PATH`
	* `DEFAULT_BANNED_PLAYERS` is now `DEFAULT_BANNED_PLAYERS_PATH`
	* `DEFAULT_BANNED_IPS` is now `DEFAULT_BANNED_IPS_PATH`
	* `DEFAULT_OPS` is now `DEFAULT_OPS_PATH`
	* `DEFAULT_JAR` is now `DEFAULT_JAR_PATH`
	* `DEFAULT_STOP_MESSAGE` is now `DEFAULT_MESSAGE_STOP`
	* `DEFAULT_STOP_ABORT` is now `DEFAULT_MESSAGE_STOP_ABORT`
	* `DEFAULT_RESTART_MESSAGE` is now `DEFAULT_MESSAGE_RESTART`
	* `DEFAULT_RESTART_ABORT` is now `DEFAULT_MESSAGE_RESTART_ABORT`
	* `DEFAULT_WORLD_BACKUP_STARTED` is now `DEFAULT_MESSAGE_WORLD_BACKUP_STARTED`
	* `DEFAULT_WORLD_BACKUP_FINISHED` is now `DEFAULT_MESSAGE_WORLD_BACKUP_FINISHED`
	* `DEFAULT_COMPLETE_BACKUP_STARTED` is now `DEFAULT_MESSAGE_COMPLETE_BACKUP_STARTED`
	* `DEFAULT_COMPLETE_BACKUP_FINISHED` is now `DEFAULT_MESSAGE_COMPLETE_BACKUP_STARTED`
* Those same name changes also effect `server.properties` files:
	* `msm-server-user` is now `msm-username`
	* `msm-log` is now `msm-log-path`
	* `msm-properties` is now `msm-properties-path`
	* `msm-whitelist` is now `msm-whitelist-path`
	* `msm-banned-players` is now `msm-banned-players-path`
	* `msm-banned-ips` is now `msm-banned-ips-path`
	* `msm-ops` is now `msm-ops-path`
	* `msm-stop-message` is now `msm-message-stop`
	* `msm-stop-abort` is now `msm-message-stop-abort`
	* `msm-restart-message` is now `msm-message-restart`
	* `msm-restart-abort` is now `msm-message-restart-abort`
	* `msm-world-backup-started` is now `msm-message-world-backup-started`
	* `msm-world-backup-finished` is now `msm-message-world-backup-finished`
	* `msm-complete-backup-started` is now `msm-message-complete-backup-started`
	* `msm-complete-backup-finished` is now `msm-message-complete-backup-finished`
* The [config wiki page][server-config] is now updated with all of these new setting names for reference.

[server-config]: https://github.com/msmhq/msm/wiki/Server-Configuration


### [0.5.2](https://github.com/msmhq/msm/compare/0.5.1...0.5.2)

* Fixed issue [#47][#47] where toggling any world's ram status would alwqays effect the first world, alphabetically, instead.

[#47]: https://github.com/msmhq/msm/issues/47


### [0.5.1](https://github.com/msmhq/msm/compare/0.5.0...0.5.1)

* Added help text line for `msm <server> config` command.


### [0.5.0](https://github.com/msmhq/msm/compare/0.4.3...0.5.0)

* Added the `msm <server> config <setting> <new-value>` command. Thanks to merging pull request [#46][#46] from [@cameroncros][@cameroncros].
* Fixed issue [#44][#44] where bash completion failed to work when accessing server conf information.

[#44]: https://github.com/msmhq/msm/issues/44
[@cameroncros]: https://github.com/cameroncros
[#46]: https://github.com/msmhq/msm/pull/46


### [0.4.3](https://github.com/msmhq/msm/compare/0.4.2...0.4.3)

* Creating a new server now automatically creates a `worldstorage` folder with a `readme.txt` inside.
* Quotations may now be used to delimit `server.properties` msm settings.


### [0.4.2](https://github.com/msmhq/msm/compare/0.4.1...0.4.2)

* Fixed the `msm <server> jar <jar-group> <jar-name>` command, which would state the `<jar-name>` was invalid, when it was not.


### [0.4.1](https://github.com/msmhq/msm/compare/0.4.0...0.4.1)

* Fixed issue [#37][#37] where the `msm stop` command would hang. Thanks to merging pull request [#38][#38] from [@escortkeel][@escortkeel].
* Fixed issue [#39][#39] where server jar symlinks were not created. Thanks to merging pull request [#40][#40] from [@jives][@jives].

[@escortkeel]: https://github.com/escortkeel
[#37]: https://github.com/msmhq/msm/issues/37
[#38]: https://github.com/msmhq/msm/pull/38
[@jives]: https://github.com/jives
[#39]: https://github.com/msmhq/msm/issues/39
[#40]: https://github.com/msmhq/msm/pull/40

### [0.4.0](https://github.com/msmhq/msm/compare/0.3.4...0.4.0)

* Added the `<server> give <player> <item> [amount] [data]` command. Closes [#34][#34]. Thanks to [@Vi3GameHkr][@Vi3GameHkr] for the suggestion.
* Added the `<server> xp <player> <amount>` command.
* Fixed the `gamemode`, `cmd` and `cmdlog` server commands *again*. They now accept multiple arguments as intented.

[@Vi3GameHkr]: https://github.com/Vi3GameHkr
[#34]: https://github.com/msmhq/msm/issues/34

### [0.3.4](https://github.com/msmhq/msm/compare/0.3.3...0.3.4)

* Fixed the `gamemode`, `cmd` and `cmdlog` server commands. They were not previously accessable.

### [0.3.3](https://github.com/msmhq/msm/compare/0.3.2...0.3.3)

* Fixed [#33][#33], by not assuming there are three invisible characters at the end of every log line.
* Fixed function names which were accidently converted to uppercase.

[#32]: https://github.com/msmhq/msm/issues/32


### [0.3.2](https://github.com/msmhq/msm/compare/0.3.1...0.3.2)

* Removes an erroneous debug statement. Fixes [#32][#32]

[#32]: https://github.com/msmhq/msm/issues/32


### [0.3.1](https://github.com/msmhq/msm/compare/0.3.0...0.3.1)

* Prevented `DEFAULT_SERVER_CONF` from being a required setting.


### [0.3.0](https://github.com/msmhq/msm/compare/0.2.16...0.3.0)

* Closes [#22][#22], by merging `server.conf` into `server.properties`.
* Removed the `DEFAULT_SERVER_CONF` setting from `/etc/msm.conf`.

[#22]: https://github.com/msmhq/msm/issues/22


### [0.2.16](https://github.com/msmhq/msm/compare/0.2.15...0.2.16)

* Improved the `msm <server> worlds ram <world>` command to insist the server is stopped.
* Refactored global variables to be uppercase. This avoids name clashes.


### [0.2.15](https://github.com/msmhq/msm/compare/0.2.14...0.2.15)

* Improved the removal of new line characters when reading Minecraft log files. Closes [#19][#19]

[#19]: https://github.com/msmhq/msm/issues/19


### [0.2.14](https://github.com/msmhq/msm/compare/0.2.13...0.2.14)

* Improved checks done before reading server files. Closes [#20][#20].

[#20]: https://github.com/msmhq/msm/issues/20


### [0.2.13](https://github.com/msmhq/msm/compare/0.2.12...0.2.13)

* Fixed the `msm <server> restart`, `msm <server> operator remove`, `msm <server> time set` and `msm <server> time add` commands.


### [0.2.12](https://github.com/msmhq/msm/compare/0.2.11...0.2.12)

* Fixed `msm <server> stop` command. Which would fail for server's not first in the list.


### [0.2.11](https://github.com/msmhq/msm/compare/0.2.10...0.2.11)

* Improved substring approach, which failed on some systems.


### [0.2.10](https://github.com/msmhq/msm/compare/0.2.9...0.2.10)

* Fixed issue [#28][#28]. Which prevented the `msm jargroup delete` command working.
* Fixed a typo in the `msm version` help text.

[#28]: https://github.com/msmhq/msm/issues/28


### [0.2.9](https://github.com/msmhq/msm/compare/0.2.8...0.2.9)

* Fixed issue [#23][#23] which caused commands to fail registration on some systems.

[#23]: https://github.com/msmhq/msm/issues/23


### [0.2.8](https://github.com/msmhq/msm/compare/0.2.7...0.2.8)

* Improved command registration to provided more accurate error messages if a command cannot be registered.
* Improved Travis CI to only test the master branch.


### [0.2.7](https://github.com/msmhq/msm/compare/0.2.6...0.2.7)

* Refactored commands such that they are now registered using a function, rather than one huge swith/case statement.


### [0.2.6](https://github.com/msmhq/msm/compare/0.2.5...0.2.6)

* Added the "all" keyword which can be used to execute a command for all servers. Closes [#16][#16]
* Added [Travis CI integration][travis] for per commit unit testing.
* Improved command error messages to give specific command errors. Closes issue [#10][#10].
* Fixed issue [#15][#15] which left a "tail" command process running after certain commands.

[#15]: https://github.com/msmhq/msm/issues/15
[#16]: https://github.com/msmhq/msm/issues/16
[#10]: https://github.com/msmhq/msm/issues/10

[travis]: http://travis-ci.org/#!/marcuswhybrow/minecraft-server-manager


### [0.2.5](https://github.com/msmhq/msm/compare/0.2.4...0.2.5)

* Added `test.sh`, a script which runs `shunit2` unit tests.
* Added test coverage for all `msm server` commands, except coverage involving running servers.
* Added option to override the configuration file used by the MSM script by specifying its path in the `MSM_CONF` environment variable.
* Added Explanatory exit codes for the MSM script (currently 0 and 64-73).
* Improved server creation to not require the "minecraft" jar group to be present.
* Fixed issue [#14][#14] which prevented world names including spaces.
* Fixed a bug listing servers when there were none.
* Fixed a bug where server names containing reserved names anywhere where considered invalid.

[#11]: https://github.com/msmhq/msm/issues/14


### [0.2.4](https://github.com/msmhq/msm/compare/0.2.3...0.2.4)

* Fixes issue [#11][#11]. Now running server's cannot be renamed, preventing not being able to stop a server renamed whilst running.

[#11]: https://github.com/msmhq/msm/issues/11


### [0.2.3](https://github.com/msmhq/msm/compare/0.2.2...0.2.3)

* Fixed issues [#12][#12] and [#13][#13] thanks to [@Jonek01][@Jonek01], which prevented the first loaded server from having the correct path, and prevented the same for the first loaded world.

[#12]: https://github.com/msmhq/msm/issues/12
[#13]: https://github.com/msmhq/msm/issues/13
[@Jonek01]: https://github.com/Jonek01


### [0.2.2](https://github.com/msmhq/msm/compare/0.2.1...0.2.2)

* Added `msm version` command, to check your installation version.


### [0.2.1](https://github.com/msmhq/msm/compare/0.2.0...0.2.1)

* Added `msm <server> worlds on|off <world>` command, which activates and deactives worlds. Inactive worlds don't get backed up, cannot be used by the server, but are included in complete server backups.
* Added `CHANGELOG.markdown`.
* Added error messages preventing script use when a required `msm.conf` settings is not specified.
* Removed fallbacks for settigns not specified in `msm.conf`.


### [0.2.0](https://github.com/msmhq/msm/compare/0.1.3...0.2.0)

* Added bash completion.
* Altered project layout to feature directories for bash_completion, cron and init files. This is useful for dealing with three files all named `msm`.
* Removed `msm <server> worlds toram` as an accessible command.
* Fixed a bug preventing `WORLD_STORAGE_PATH` from being overriden by a server.
* Fixed error in `msm.conf`


### [0.1.3](https://github.com/msmhq/msm/compare/0.1.2...0.1.3)

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


### [0.1.2](https://github.com/msmhq/msm/compare/0.1.1...0.1.2)

* Added in-game messages, and correct pausing and resuming of Minecraft world saving around world backups.
* Added "Backup took n seconds." message to world and complete backups.
* Fixed a bug listing a server's worlds.
* Fixed a bug when stopping all servers in parallel, where a server might not stop at the correct time.


### [0.1.1](https://github.com/msmhq/msm/compare/0.1.0...0.1.1)

* Fixed a bug preventing world links from being created.
* Fixed a bug where a server stopped as the result of a command would not always be made inactive.


### [0.1.0](https://github.com/msmhq/msm/compare/aaaa69dbd6c45ff0eb7dd69b756b4103abadfbc4...0.1.0)

* Refactored and improved protype code.
* Transitioned from managing a single to multiple servers.
