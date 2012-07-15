---
layout: docs
title: Setup Commands · Docs · MSM
tab: docs
docsitem: commands-setup
---

Setup Commands
==============

These commands help you create, delete and organise Minecraft servers.

    server list

Outputs a list of the Minecraft server directories managed by the Minecraft Server Management script.

    server create <name>

Creates a new Minecraft server directory, and downloads the latest Minecraft jar. Servers are created in the directory specified by SERVER_STORAGE_PATH in the configuration file.

    server delete <name>

Deletes a server directory, asking for confirmation first. Remember, this will delete anything stored within that server directory also, including worlds, and any plugins.

    server rename <name> <new-name>

Renames a server directory. This can only be done when a server is not running.
