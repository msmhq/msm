# Minecraft Server Manager [![Build Status](https://secure.travis-ci.org/marcuswhybrow/minecraft-server-manager.png?branch=master)](http://travis-ci.org/marcuswhybrow/minecraft-server-manager) [![Donate to MSM](https://www.paypalobjects.com/en_GB/i/btn/btn_donate_SM.gif)](https://www.paypal.com/cgi-bin/webscr?cmd=_s-xclick&hosted_button_id=Z7XQDNF7U5GLL)

A single init script which makes running multiple Minecraft/Bukkit servers easier for us admins.

## Quick Installers

Install MSM on a **debian box**:

    wget -q http://git.io/Sxpr9g -O /tmp/msm && bash /tmp/msm

Or install MSM under **RedHat**:

	wget -q http://git.io/lu0ULA -O /tmp/msm && bash /tmp/msm

Or [suggest a new platform][issues].

## Getting Started

* [Install][install] MSM on your box.
* [Visit the docs][docs] for help with getting started.
* Fork and contribute your own modifications.
* Not a coder? [Post an issue][issues] with your idea for MSM instead.
* Read the [changelog][changelog] to get a picture of how MSM has evolved over time.

MSM is released under the GPLv3 licence, which is included in the repository [here][licence]. I'm open to suggestions where licencing is concerned.

## Features

As well as starting, stopping and restarting MSM has the following features:

* One script handles multiple servers, run two or more servers on one machine.
* Can create and start new servers with a single command, downloads the jars for you.
* Periodically makes [WorldEdit snapshot][we-snapshot] compatible backups of your worlds.
* Backup the entire server directory for complete protection.
* Load world's into RAM for faster access (reduces lag).
* Easily configurable global defaults, with per server overrides if needed.
* Apply server commands to one, multiple, or all servers in one go (useful for whitelisting a player on all servers.)
* Tab completion for all commands, makes everything faster and getting started a breeze.
* Keep server logs organsied by periodically "rolling" them.
* Organises jar files into groups (such as minecraft and craftbukkit) and links each server to a single jar. Includes automated download of new versions.
* Plethora of in-game commands (whitelist, blacklist, operator, gamemode, kick, say, time, toggledownfall, save)
* Send commands straight to the server via the command line.

## Support

1. If you find a problem with MSM and you think the problem is one that requires changing code [submit an issue][issues] via GitHub.
2. If you otherwise have a problem, question or suggestion you can email me directly at msm@marcuswhybrow.net.


## Upcomming features

* **QuickBackup:** If you store your backups non-locally (maybe on a NAS), QuickBackup optionally creates a backup locally for speed, and then moves it after your players are building again! My initial testing shows a 54 second network backup confaltes to 23 seconds of in-game time.
* **Restore:** Roll-back to an old world or whole server backup automatically.


## Versioning

MSM uses semantic version numbers to better describe what code one might have installed, and indicate backwards incompatible changes.

Releases will be numbered in the following format:

`<major>.<minor>.<patch>`

And constructed with the following guidelines:

* Breaking backward compatibility bumps the major (and resets the minor and patch)
* New additions without breaking backward compatibility bumps the minor (and resets the patch)
* Bug fixes and misc changes bumps the patch

For more information on SemVer, visit http://semver.org/.


## Acknowledgements

This code grew out of an old version of [Ahtenus' Minecraft Init Script][ahtenus-minecraft-init].

[we-snapshot]: http://wiki.sk89q.com/wiki/WorldEdit/Snapshots
[ahtenus-minecraft-init]: https://github.com/Ahtenus/minecraft-init
[docs]: http://msmhq.com/docs/
[install]: http://msmhq.com/docs/installation.html
[changelog]: https://github.com/marcuswhybrow/minecraft-server-manager/blob/master/CHANGELOG.markdown
[licence]: https://github.com/marcuswhybrow/minecraft-server-manager/blob/master/LICENSE.markdown
[issues]: https://github.com/marcuswhybrow/minecraft-server-manager/issues
