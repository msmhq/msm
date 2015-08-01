---
layout: docs
title: Jar Group Commands · Docs · MSM
tab: docs
docsitem: commands-jar-group
---

Jar Group Commands
==================

Jar groups help you to manage the different jar files you might have. Each jar group is a directory and URL, the URL specifies where the latest version of this logical group of jars can be downloaded. The "minecraft" and "craftbukkit" jar groups are automatically created for you.

Storing jar files in a single location helps to organise things, each server has a symbolic link to the jar file it currently uses, all of this meaning you only have to download the new version of a server jar once.

    jargroup list

Lists the stored jar files grouped by jar group.

    jargroup create <name> <download-url>

You can create your own jar groups giving them a name and a URL specifying where the latest version can always be found and downloaded. New jar groups are stored along side the existing "minecraft" and "craftbukkit" groups.

    jargroup create <name> minecraft
    jargroup create <name> minecraft-snapshot

If you'd like your jar group to download the latest release or snapshot version of minecraft-server, the download-url parameter supports minecraft (for release) or minecraft-snapshot (for snapshot).

    jargroup delete <name>

If you don't need a group of jars anymore you can delete all of them using this command. You won't be able to delete a jar group if any server is currently using a jar in that group.

    jargroup rename <name> <new-name>

For your convenience this command allows you to rename an existing jar group to a name you choose.

    jargroup changeurl <name> <download-url>

If a download URL for a jar group should change, or you made a mistake, you can change it to the correct URL with this command.

    jargroup getlatest <name>

When a new version of a jar has been released, this command allows you to download (based upon the download URL specified for the group) and save it within it's appropriate jar group directory.

The latest version is always downloaded when this command is invoked, then checked for differences against the existing latest jar in that group. If there are differences the new version is saved as the latest version, if not the downloaded version is removed.
