---
layout: docs
title: Installation · Docs · MSM
tab: docs
docsitem: installation
---

Installation
============

**Note:** You must have both `sudo` and `cron` install on your system for MSM to function properly.

Install MSM on a Debian server (Ubuntu etc.) by pasting this line in your terminal:

    wget -q http://git.io/Sxpr9g -O /tmp/msm && bash /tmp/msm

Or on RedHat using:

    wget -q http://git.io/lu0ULA -O /tmp/msm && bash /tmp/msm

Or [suggest a new platform][issues].

You can paste the links in your browser to read the script that will be executed. Or, if you don't want to trust a script, install MSM manually:

Manual Installation
-------------------

1. This script uses the following tools: **screen**, **rsync**, **zip**. Let's install them:

        sudo apt-get update
        sudo apt-get install screen rsync zip

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

12. Move generated worlds to the world storage folder:

    When a Minecraft server starts for the first time, it generates a new world according to the `level-name` entry in `server.properties`, which is `world` by default.

    MSM requires that you move this newly generated world (and any accompanying folders such as `world_nether` or `world_the_end` in the case of CraftBukkit) to a new directory which you must create inside of the server root folder called `worldstorage`. 
    
    After moving them, run `sudo msm <server-name> worlds load` to create the symlinks needed for the server to recognize  new worlds.
    
    See [Concepts &rarr; Server Layout][layout] for further explanation.


Check Out The Commands
----------------------

If stuck, you can always type `msm help` for a quick list of all commands, or if you want to know more about a command read the [command documentation][commands].

[config]: {{ site.baseurl }}/docs/configuration/
[cron]: {{ site.baseurl }}/docs/concepts/cron.html
[commands]: {{ site.baseurl }}/docs/commands/
[#26]: https://github.com/marcuswhybrow/minecraft-server-manager/issues/26
[versions]: https://github.com/marcuswhybrow/minecraft-server-manager/blob/latest/versioning/versions.txt
[layout]: {{ site.baseurl }}/docs/concepts/layout.html
[issues]: https://github.com/marcuswhybrow/minecraft-server-manager/issues
