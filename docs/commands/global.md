---
layout: docs
title: Global Commands · Docs · MSM
tab: docs
docsitem: commands-global
---

Global Commands
===============

    start

The global start command starts all active servers, leaving inactive servers not running.

    stop [now]

The global stop command stops all servers, but does not effect a servers active/inactive status.

    restart [now]

The global restart command stops all servers, and then restarts only the previously active servers.

*See the [Active vs Inactive Servers Page][active-vs-inactive] for more information on the differences between an active and an inactive server.*

    version

Displays the verison of MSM you are running. Useful when reporting bugs, or looking to upgrade.

    config

Displays a list of values MSM is using that can be configured. Anything in this list can be overriden in `/etc/msm.conf`.

    update [--noinput]

Updates MSM files to the latest versions. Will ask for confirmation before overwriting your files. However one may disable this requirement with the `--noinput` flag.


[active-vs-inactive]: {{ site.baseurl }}/docs/concepts/active-vs-inactive.html