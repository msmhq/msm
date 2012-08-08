---
layout: docs
title: Installation · Docs · MSM
tab: docs
docsitem: installation
---

Installation
============

This script uses the following tools: **screen**, **rsync**, **zip**

*There may be a script to do this in the future, but its quite easy to do yourself, and I prefer to place things manually whenever possible so I know what's going on myself.*

## Install The Configuration File and Create Directories

1. Download the configuration file and place it in the `/etc` directory:

        sudo wget https://raw.github.com/marcuswhybrow/minecraft-server-manager/latest/msm.conf -O /etc/msm.conf

2. Edit `/etc/msm.conf`, have a read through and make any adjustments, the comments will help you out.

3. Create and give permissions to directories used by MSM, use the values specified in `/etc/msm.conf` if you change them:

    With the default configuration, all files are stored within the `/opt/msm` directory.

        sudo mkdir /opt/msm
        sudo chown minecraft /opt/msm
        sudo chmod -R 775 /opt/msm

    Also be sure to create the directory for ram storage if your system supports ram disk:

        sudo mkdir /dev/shm/msm
        sudo chown minecraft /dev/shm/msm
        sudo chmod -R 775 /dev/shm/msm

## Install The Script

1. Download the Minecraft Server Manager script (called msm) placing it in the `/etc/init.d` directory:

        sudo wget https://raw.github.com/marcuswhybrow/minecraft-server-manager/latest/init/msm -O /etc/init.d/msm

2. Give the script the required permissions:

        sudo chmod 755 /etc/init.d/msm

3. Update the rc.d links, this allows–for example–servers to be started at system startup automatically:

        sudo update-rc.d msm defaults


## Configure Cron (Scheduled Actions)

1. Download the cron file and place it in `/etc/cron.d` with the name `msm`:

        sudo wget https://raw.github.com/marcuswhybrow/minecraft-server-manager/latest/cron/msm -O /etc/cron.d/msm

2. If you changed `SERVER_USER` in `/etc/msm.conf` to something other than "minecraft", edit `/etc/cron.d/msm` and change the occurrences of the username "minecraft" on each line, to the username you chose.

*See the [Using Cron Page][cron] if you are unfamiliar with cron, or for an explanation of how cron is being used in this case.*


## Enable Bash Tab Completion Support

1. Download the bash_completion file and place it in `/etc/bash_completion.d` with the name `msm`:

        sudo wget https://raw.github.com/marcuswhybrow/minecraft-server-manager/latest/bash_completion/msm -O /etc/bash_completion.d/msm

2. Ensure the following lines are uncommented in your `~/.bashrc` file or if you don't have that, find it in `/etc/bash.bashrc` instead:

        if [ -f /etc/bash_completion ]; then
            . /etc/bash_completion
        fi

3. Bash completion will now be active when you login or if you want it right away do:

        source /etc/bash_completion.d/msm

Now whenever you are stuck for what comes next in a command, you can press tab to auto-complete options.

## Simplify Command

1. Place a symbolic link to the script in `/usr/local/bin`

        sudo ln -s /etc/init.d/msm /usr/local/bin/msm

This allows you to do `msm bla bla` instead of the slightly longer `/etc/init.d/msm bla bla`.

### Setup A Server

1. If you have an existing server, making sure it's stopped you should move it to the server store (**WARNING Not yet supported!**). This performs a straight move, but also organises (if present) the jar file, and moves existing worlds to the "world storage" directory:

        sudo msm server transition <new-name> /path/to/existing/server/dir

    or if you are starting from scratch, first create a jar group which also downloads the latest version:

        sudo msm jargroup create minecraft https://s3.amazonaws.com/MinecraftDownload/launcher/minecraft_server.jar

    or if you want to run a Bukkit server:

        sudo msm jargroup create craftbukkit http://dl.bukkit.org/latest-rb/craftbukkit.jar

    Remember, you can have as many jar groups as you like, and call them whatever you want.

    Next create a new server, and choose the jargroup you just made:

        sudo msm server create <server-name>
        sudo msm <server-name> jar minecraft

2. Start your server:

        sudo msm <server-name> start

    Currently in version 0.2.x, servers will create new worlds in the server directory instead of directly in the `worldstorage` directory. Once a new world is created you unfortunately have to stop the server and move the world into (and probably create) the `worldstorage` directory yourself. The `worldstorage` directory should be created in the server directory. Follow issue [#26][#26] for changes to this area.


## Check Out The Commands

If stuck, you can always type `/etc/init.d/msm help` for a quick list of all commands, or if you want to know more about a command read the [command documentation][commands].

[config]: {{ site.baseurl }}docs/configuration/
[cron]: {{ site.baseurl }}docs/concepts/cron.html
[commands]: {{ site.baseurl }}docs/commands/
[#26]: https://github.com/marcuswhybrow/minecraft-server-manager/issues/26