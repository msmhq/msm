---
layout: docs
title: Commands Overview · Docs · MSM
tab: docs
docsitem: commands-overview
---

Commands
========

The Minecraft Server Manager script has different options or commands you can invoke by adding more words after it on the command line.

For this command documentation, required phrases are written as is (i.e. `start`), required variables are denoted by angle bracket placeholders (i.e. `<server>`), and optional phrases are denoted by square brackets (i.e `[now]`.) A variable placeholder should be **replaced** with a string as required (and described) by each individual command.

You can always get a quick list of commands by typing `msm help`, which outputs the following info:

	$ msm help
	Usage: msm command:

	--Setup Commands------------------------------------------------
	  server list                                   List servers
	  server create <name>                          Creates a new Minecraft server
	  server delete <name>                          Deletes an existing Minecraft server
	  server rename <name> <new-name>               Renames an existing Minecraft server

	--Server Mangement Commands-------------------------------------
	  <server> start                                Starts a server
	  <server> stop [now]                           Stops a server after warning players, or right now
	  <server> restart [now]                        Restarts a server after warning players, or right now
	  <server> status                               Show the running/stopped status of a server
	  <server> connected                            List a servers connected players
	  <server> worlds list                          Lists the worlds a server has
	  <server> worlds load                          Creates links to worlds in storage for a server
	  <server> worlds ram <world>                   Toggles a world's "in RAM" status
	  <server> worlds todisk                        Synchronises any "in RAM" worlds to disk a server has
	  <server> worlds backup                        Makes a backup of all worlds a server has
	  <server> worlds on|off <world>                Activate or deactivate a world, inactive worlds are not backed up
	  <server> logroll                              Move a server log to a gziped archive, to reduce lag
	  <server> backup                               Makes a backup of an entire server directory
	  <server> jar <jargroup> [<file>]              Sets a server's jar file
	  <server> console                              Connects to the interactive console. Access may be limited
	  <server> config [<setting> <value>]           Lists server settings, or sets a specific setting.

	--Server Pass Through Commands----------------------------------
	  <server> wl on|off                            Enables/disables server whitelist checking
	  <server> wl add|remove <player>               Add/remove a player to/from a server's whitelist
	  <server> wl list                              List the players whitelisted for a server
	  <server> bl player add|remove <player>        Ban/pardon a player from/for a server
	  <server> bl ip add|remove <ip address>        Ban/pardon an IP address from/for a server
	  <server> bl list                              Lists the banned players and IP address for a server
	  <server> op add|remove <player>               Add/remove operator status for a player on a server
	  <server> op list                              Lists the operator players for a server
	  <server> gm survival|creative <player>        Change the game mode for a player on a server
	  <server> kick <player>                        Forcibly disconnect a player from a server
	  <server> say <message>                        Broadcast a (pink) message to all players on a server
	  <server> time set|add <number>                Set/increment time on a server (0-24000)
	  <server> toggledownfall                       Toggles rain and snow on a server
	  <server> give <player> <item> [amount] [data] Gives an entity to a player
	  <server> xp <player> <amount>                 Gives XP to, or takes away (when negative) XP from, a player
	  <server> save on|off                          Enable/disable writing world changes to file
	  <server> save all                             Force the writing of all non-saved world changes to file
	  <server> cmd <command>                        Send a command string to the server and return
	  <server> cmdlog <command>                     Same as 'cmd' but shows log output afterwards (Ctrl+C to exit)

	--Jar Commands--------------------------------------------------
	  jargroup list                                 List the stored jar files.
	  jargroup create <name> <download-url>         Create a new jar group, with a URL for new downloads
	  jargroup delete <name>                        Delete a jar group
	  jargroup rename <name> <new-name>             Rename a jar group
	  jargroup changeurl <name> <download-url>      Change the download URL for a jar group
	  jargroup getlatest <name>                     Download the latest jar file for a jar group

	--Global Commands-----------------------------------------------
	  start                                         Starts all active servers
	  stop [now]                                    Stops all running servers
	  restart [now]                                 Restarts all active servers
	  version                                       Prints the Minecraft Server Manager version installed
	  config                                        Displays a list of the config values used by MSM



[op]: http://www.minecraftwiki.net/wiki/Operator
[survival]: http://www.minecraftwiki.net/wiki/Survival
[creative]: http://www.minecraftwiki.net/wiki/Creative
[we-snapshots]: http://wiki.sk89q.com/wiki/WorldEdit/Snapshots
[active-vs-inactive]: https://github.com/marcuswhybrow/minecraft-server-manager/wiki/Active-vs-Inactive-Servers