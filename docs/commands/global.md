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

*See the [Active vs Inactive Servers Page][active-vs-inactive] for a more information on the differences between an active and an inactive server.*


[active-vs-inactive]: {{ site.baseurl }}docs/concepts/active-vs-inactive.html