---
layout: docs
title: Global Configuration · Docs · MSM
tab: docs
docsitem: configuration-msm
---

Gloal Configuration
===================

Configuration for MSM is handled by a single file: `/etc/msm.conf` by default, and the actual concrete values used by MSM may be listed using the `msm config` command. Each item listed by `msm config` may be overriden in `/etc/msm.conf`

Changing the Config Location
----------------------------

If you want MSM to use a file in a different location for determining its configuration, set the `MSM_CONF` evironment variable with the path to the new file:

    export MSM_CONF="/path/to/new/msm.conf"

With this approach you will have to export that variable for each of the bash profiles MSM might be ran from. This is, as you can imagine, a less than satisfactory workflow. Comming soon MSM will check the standard location for environment setup: `/etc/profile.d/msm`. This will allow a one time override of the `MSM_CONF` variable.

Default Values
--------------

MSM has default values for all possible settings. This means that `/etc/msm.conf` could be left blank and MSM would still work with sensible defaults.

It is however advised that each value be explicitly defined nonetheless, allowing admins to be sure of what MSM is doing. The <a href="{{ site.baseurl }}/docs/installation.html">installation documentation</a> encourages this approach by downloading the <a href="https://github.com/msmhq/msm/blob/latest/msm.conf">msm.conf file</a> ready with default values.

Server Defaults
---------------

Variables beginning witht `DEFAULT_` are default values for unconfigured servers. Variables of this kind that also end in `_PATH` are linux paths interpreted relative to the arbitrary server directory when evaluated.

Path Values
-----------

Path values can be absolute (starts with a forward slash) or relative (does not start with a forward slash). Paths for server defaults, as above mentioned, when relative are relative to the root directory for each server.

All other path settings must be absolute, these include the general MSM settings such as where to store all servers, backups and logs.

Paths must **never** have a trailing slash. In other words `/opt/msm/servers` is OK but `/opt/msm/servers/` is not.
