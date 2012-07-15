---
layout: docs
title: Active vs Inactive · Docs · MSM
tab: docs
docsitem: concepts-active-vs-inactive
---

Active vs Inactive
==================

Whether a server is active or inactive is a reflection of whether you (the administrator) intend the server to be running or not. If the server should be running it is active, if it should be off or not in use, then it is inactive.

## Why Is This Useful?

If a server is found to be not running, it would otherwise not be possible to deduce whether the administrator means it to be stopped, or whether the server has crashed. In this case the Minecraft Server Manager script looks to the active/inactive status of the server to determine what to do.

On the other hand, and just as useful is the ability to discern whether a running server should be running to should be stopped. This prevents servers from starting again by whatever means and utilising resources or providing a service you–as the administrator–did not want to provide.

## Behind The Scenes

The active/inactive status of a server is defined by the presence or absence of the file `active` located in the server directory.

Note: Although there will never be a command to directly effect the active/inactive status of a server outside the start/stop/restart commands, an administrator may alter a servers status, if they had a reason to, at any time by creating or removing the aforementioned `active` file, which is just an empty file created simply using `touch active` (be sure to give it the same permissions as other server files so MSM can manage it.)

## Commands That Effect The Active/Inactive Status

    <server> start
    <server> restart [now]

The server start, and restart commands always moves a server into the "active" state.

    <server> stop [now]

The server stop command always moves a server into the "inactive" state.

*See the [Commands Page][commands] for more information on commands.*

### Notable Exceptions

Stopping the server in any other way–not through the aforementioned `<server> stop` command–does not change the active/inactive state of a server. For example the in-game `/stop` command which can be called by players given operator status, will be regarded as a non-intentional stop, and therefor leave the server in the active state.

In the same way, starting the server using any method other than the `<server> start` command will not effect the servers active/inactive status in any way.

## Commands That Recognise The Active/Inactive Status

    start

The global start command starts all *active* servers, leaving any inactive servers (intentionally stopped) not running.

    stop [now]

The global stop command stops all servers, leaving the active/inactive state of a server alone.

    restart [now]

The global restart command stops all severs, preserving the active/inactive state, and only restarts those which were active. This means that—for example—any crashed servers will be found stopped but active, and will be started after this command.

*Again, see the [Commands Page][commands] for more information on commands.*


[commands]: {{ site.baseurl }}docs/commands/