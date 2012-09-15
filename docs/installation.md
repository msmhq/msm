---
layout: docs
title: Installation · Docs · MSM
tab: docs
docsitem: installation
---

Installation
============

Install MSM on a Debian server (Ubuntu etc.) by pasting this line in your terminal:

    wget -qO- http://git.io/YhclQA | sh

Have a read of the script [http://git.io/KxE3HQ](http://git.io/KxE3HQ) it it calms your nerves. Or, install MSM manually:

Manual Installation
-------------------

1. This script uses the following tools: **screen**, **rsync**, **zip**. Let's install them:

        sudo apt-get update
        sudo apt-get install screen rsycn zip

2. Download MSM's default configuraiton file, then **read through it** and change anything you want. Be sure to carry any changes you make forwards through this guide.

        sudo wget http://git.io/6eiCSg -O /etc/msm.conf

3. Create the directories MSM needs:

    First create the folder where MSM will store server, jars, and other files:

        sudo mkdir /opt/msm
        sudo chown minecraft /opt/msm
        sudo chmod -R 775 /opt/msm

    If you want to store your world folders in RAM for a performance boost, create this directory in ramdisk (`/dev/shm` is the default in Ubuntu):

        sudo mkdir /dev/shm/msm
        sudo chown minecraft /dev/shm/msm
        sudo chmod -R 775 /dev/shm/msm

4. Download the MSM script and place it in `/etc/init.d`:

        sudo wget http://git.io/J1GAxA -O /etc/init.d/msm

5. Set script permissions, and integrate script with startup/shutdown:

        sudo chmod 755 /etc/init.d/msm
        sudo update-rc.d msm defaults

6. Create a shortcut so we can use just type `msm`:

        sudo ln -s /etc/init.d/msm /usr/local/bin/msm

7. Ask MSM to update, getting the latest files:

        sudo msm update

8. Setup MSM's included cron script for sheduled tasks and force cron to load script:

        sudo wget http://git.io/pczolg -O /etc/cron.d/msm
        sudo service cron reload

    If you understand cron, give it a read and alter the timings to your liking.

9. Create a jar group to manage current and future Minecraft versions:

        sudo msm jargroup create minecraft https://s3.amazonaws.com/MinecraftDownload/launcher/minecraft_server.jar

10. Create a new server, and tell it to use the latest of those Minecraft jars:

        sudo msm server create <server-name>
        sudo msm <server-name> jar minecraft

11. Start your new server:

    Before you start make sure to let MSM know which version of Minecraft you are running by adding the following line to `/opt/msm/servers/<server-name>/server.properties`

        msm-version=minecraft/1.3.1

    The available strings are currently:

    * `minecraft/1.2.0` and above.
    * `craftbukkit/1.2.0` and above.

    See [versioning/versions.txt][versions] for the up to the minute list.

    Once that clunky step it out the way (I'm working on making it more transparent) let's start the server:

        sudo msm <server-name> start

[versions]: https://github.com/marcuswhybrow/minecraft-server-manager/blob/latest/versioning/versions.txt


Check Out The Commands
----------------------

If stuck, you can always type `msm help` for a quick list of all commands, or if you want to know more about a command read the [command documentation][commands].

[config]: {{ site.baseurl }}/docs/configuration/
[cron]: {{ site.baseurl }}/docs/concepts/cron.html
[commands]: {{ site.baseurl }}/docs/commands/
[#26]: https://github.com/marcuswhybrow/minecraft-server-manager/issues/26