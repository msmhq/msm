---
layout: docs
title: Server Configuration · Docs · MSM
tab: docs
docsitem: configuration-server
---

Server Configuration
====================

Default values which apply to all servers can be found in `/etc/msm.conf` under the heading "Server Defaults". See [Configuration &rarr; Global Settings][global-config] intead for setitngs which apply to all servers.

If a server is not configured individually–and it does not have to be–global default values are used instead. You may however override these defaults for a particular server by editing its `server.properties` file found in the respective server directory, and adding any of the settings from the following sections.

*Note: MSM's default values are used for every unspecified option, so `server.properties` can be left completely unedited.*

## A Note on Paths

A setting ending `-path` accepts a relative or absolute path. Relative paths are appended to the location of the server directory. Neither absolute nor relative paths should end in a trailing slash, otherwise MSM will interpret them incorrectly.

Remember an absolute path is one that starts with a forwards slash, whereas a relative path does not. 

## Invocation

These options cover starting the server in the way you want:

    msm-username=                    # The linux user used to start the server
    msm-version=                     # Set the Minecraft verion so MSM knows how to talk it's language **IMPORTANT**
    msm-screen-name=                 # The name used to reference the screen session for this server
    msm-world-storage-path=          # The directory (relative to the server directory) to store worlds in
    msm-world-storage-inactive-path= # The directory to store unused worlds

    msm-log-path=                    # The location of the Minecraft server log file
    msm-whitelist-path=              # The location of the whitelist file
    msm-banned-players-path=         # The location of banned players file
    msm-banned-ips-path=             # The location of the banned ups file

    msm-ops-path=                    # The location of the ops file
    msm-jar-path=                    # The jar file (relative to the server directory) to launch
    msm-ram=                         # The amount of RAM (in MBs) to dedicate to this server
    msm-invocation=                  # The line called to start the server: you may use the tags "{JAR}" and
                                     # "{RAM}" in this string which will be replaced with the values from
                                     # the previous settings

### Versioning Options

The `msm-version` setting is a code comprised of two parts: the server type, currently "minecraft" or "craftbukkit", and then the dot separated version number (e.g. `1.3.1`). For example `msm-version=minecraft/1.3.1` or `msm-version=craftbukkit/1.2.5`.

MSM selects, using this code, the most appropiate included versioning file, so the previous examples might resolve to "minecraft/1.3.0" and "craftubukkit/1.2.0" respectively. These **versioning files help MSM interact with differing Minecraft updates**.

If you are running a unique server type, maybe Tekkit for instance, pick the closest available version. Or let me know, by [creating an issue][issues], that you think MSM should accomodate it.

[issues]: https://github.com/marcuswhybrow/minecraft-server-manager/issues/
    
## Delays

When a server is going to stop, it's nice to give users some time to finish up:

    msm-stop-delay=                  # The delay between warning players and stopping the server
    msm-restart-delay=               # The delay between warning players and restarting the server

## Messages

When interruptions to the normal experience occur, it's nice to give players some explanation:
    
    msm-message-stop=                      # The message shown to players before stopping the server
    msm-message-stop-abort=                # The message shown to players when a shut-down is aborted
    msm-message-restart=                   # The message shown to players before restarting the server
    msm-message-restart-abort=             # The message shown to players when a restart is aborted
    msm-message-world-backup-started=      # The message shown to players when starting to backup server worlds
    msm-message-world-backup-finished=     # The message shown to players when finished backing up server worlds
    msm-message-complete-backup-started=   # The message shown to players when starting to backup the entire server
    msm-message-complete-backup-finished=  # The message shown to players when finished backing up the entire server

[global-config]: {{ site.baseurl }}/docs/configuration/msm.html