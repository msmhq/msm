---
layout: docs
title: Server Commands · Docs · MSM
tab: docs
docsitem: commands-server
---

Server Commands
===============

Once you have created a server, you can control it using "server commands", they all have the same syntax: `<server> <command>` the server name followed by the command to invoke. These commands help you with all the behind the scenes stuff that goes into being a good administrator.

    <server> start

Starts a server, replace `<server>` with the name of the server you wish to start for example: `creativeserver start`.

Starting a server using this command puts it into the "active" state, which is a recognition that the server was started intentionally. This is the only way to show you intentionally started the server.

*See the [Active vs Inactive Servers Page][active-vs-inactive] for a more detailed explanation on the concept.*

    <server> stop [now]

`creativeserver stop` would stop a running server after warning in-game players, and then after (by default) a 10 second delay. `creativeserver stop now` issues no such warning and incurs no delay.

Stopping a server using this command puts it into the "inactive" state, recognising that it has been stopped intentionally.

*Again, see the [Active vs Inactive Servers Page][active-vs-inactive] for a more detailed explanation on the concept.*

    <server> restart [now]

In the same way as the previous "stop" command, the restart command reboots a server after giving a warning to players, and then after a delay stopping the server, and then starting it again. Adding the word "now" on the end–as before–removes the warning and the delay.

Restarting a server using this command puts it into the "active" state (if it was not already), recognising that is has been intentionally started.

*Again, see the [Active vs Inactive Servers Page][active-vs-inactive] for a more detailed explanation on the concept.*

    <server> status

Outputs whether or not a server is running or stopped.

    <server> connected

Lists the players which are currently in-game on a server.

    <server> worlds list

Lists the Minecraft worlds used by a server. Worlds can be found in a server's "world storage" directory.

    <server> worlds load

Enables worlds in a servers "world storage" directory to be used by the Minecraft server. This is achieved by creating symbolic links in the server directory that point to the locations of each world directory. **This is done automatically for you on server start up.**

    <server> worlds ram <world>

Toggles the status of a server's world between "in RAM" and "on disk". On disk is the typical/default state from which Minecraft worlds are ran, in other words the world is located somewhere on the systems hard drive disk. A world which is loaded into RAM can be accessed at higher speeds, which means reduction in lag for servers with many players all exploring the same world.

Use `<server> worlds list` to check which worlds are currently in RAM and which are on disk.

    <server> worlds todisk

RAM is a volatile storage type, meaning when a system shuts down anything stored in RAM is lost. This is a danger to any Minecraft worlds stored in RAM for the speed increase. This command synchronises all changes made to a world in RAM back to the persistent version stored on disk. **This is done automatically for you on server shut down**, but should also be scheduled in a cron job to run every (for example) half hour to mitigate losses caused by a computer crash.

    <server> worlds backup

Backs up all worlds a server has by copying the world folder, compressing it's size into a zip file and storing it in the world archive (config: WORLD_ARCHIVE_PATH). These backups/archives are compatible for use with [WorldEdit snapshots][we-snapshots], which allow you to restore regions of you world in-game to a previous saved state.

Backups need to be scheduled in a cron job, and can be run as often or as occasionally as one likes.

    <server> logroll

Moves the contents of the Minecraft server log to a timestamped file archive, and compresses it. Running this command periodically helps to keep the server log file size small. A large log size may cause a server to slow down when it writes to log, and partitioning logs into files in this way can help with house-keeping.

    <server> backup

Backing up a server's worlds is covered by the `<server> worlds backup` command, however there is more to a server that just worlds, think: configuration, plugins, whitelists, blacklists. All of that would be lost to a hard drive failure. This command zip's up the entire server directory and everything in it, making starting again after a storage loss simple by removing the need to configure everything afresh.

    <server> jar <jargroup> [<file>]

Sets a the jar file a server will use when starting up. If only a jar group is specified the latest jar in that group will be used. However if a file name is also given, it will be looked for in the jar group and used instead.

    <server> console

Resumes and connects the user to the underlying screen session which runs the server. "screen" is a tool which allows a terminal to run in the background and be returned to when necessary. The `<server> console` command "resumes" the screen session allowing direct access to the server console.

**Warning: Pressing Control+C in a screen session will terminate/stop/kill the Minecraft server.** To leave a screen session (and the server) running use the "Control+A then D" key combination.

    <server> config <setting> <value>

Each vanilla Minecraft server has a `server.properties` file which controls aspects from which map to load, to the game difficulty, to special settings understandable only to Minecraft Server Manager. These values can be change by editing the file directly, or by using this command.

If a setting already exists in the file, it will be updated with the new value, if it does not yet exist a new entry will be appended to the end of the file.

This command may be a good entry point for admins not familiar with command line file editors such as vim or emacs.



[we-snapshots]: https://wiki.sk89q.com/wiki/WorldEdit/Snapshots
[active-vs-inactive]: {{ site.baseurl }}/docs/concepts/active-vs-inactive.html
