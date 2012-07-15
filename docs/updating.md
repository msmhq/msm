---
layout: docs
title: Updating · Docs · MSM
tab: docs
docsitem: updating
---

Updating
========

New versions of MSM contain new features, more reliable code, bug fixes, and contributions from the community. Since MSM is a single script replacing it with the new version is enough to get up to date.

The MSM code repository is managed by git source control, and has tags which mark versions as well as a special tag called `latest` which always points to the recommended latest version.


Downloading the Latest Version
------------------------------

    sudo wget https://raw.github.com/marcuswhybrow/minecraft-server-manager/latest/init/msm -O /etc/init.d/msm
    sudo wget https://raw.github.com/marcuswhybrow/minecraft-server-manager/latest/bash_completion/msm -O /etc/bash_completion.d/msm

Those two commands will retrive the MSM script, and the accompanying bash completion script, putting them in the recommended locations.


Check Your Version
------------------

    sudo msm version

The version command will return the current MSM version you have installed, useful for checking the upgrade was successful.