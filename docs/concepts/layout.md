---
layout: docs
title: Server Layout · Docs · MSM
tab: docs
docsitem: concepts-layout
---

Server Layout
=============

MSM is designed to get out of your way as much as possible. A server directory should look nearly exactly the same except for a few key changes.

1. Create a `worldstorage` folder and move all your worlds into it.
2. Optionally create the `worldstorage_inactive` folder for inactive worlds.

So a typical Bukkit server managed by MSM would look like this:

	/opt/msm
	`- servers
	   `- server_name
	      |- plugins/
	      `- worldstorage/
	         |- world/
	         |- world_nether/
	         `- world_the_end/
	      |- worldstorage_inactive/
	      |- server.jar
	      |- server.properties
	      |- server.log
	      |- ops.txt
	      |- white-list.txt
	      |- banned-ips.txt
	      `- banned-players.txt

MSM expects to find your worlds in the `worldstorage` folder. This allows MSM to reliably determine which folders are worlds for making backups, moving worlds into RAM and other behind the scenes stuff.

**Advanced:** *You can change the location MSM expects to find world using the `msm-world-storage-path` setting in your server's `server.properties` file. See [Configuration &rarr; Individual Server][config-server] for more info.*

[config-server]: {{ site.baseurl }}/docs/configuration/server.html