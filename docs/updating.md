---
layout: docs
title: Updating · Docs · MSM
tab: docs
docsitem: updating
---

Updating
========

New versions of MSM contain new features, more reliable code, bug fixes, and contributions from the community. Since version `0.8.0` MSM has an inbuilt command to download and install everything for you automatically. Don't worry, **it will ask for confirmation** before overwriting your files.

    sudo msm update


Manually updating
-----------------

If you are upgrading from a version of MSM that does not have the update command, you will have to do one thing by hand. Download and replace your MSM script file with the latest version:

    sudo wget https://raw.github.com/marcuswhybrow/minecraft-server-manager/latest/init/msm -O /etc/init.d/msm
    sudo wget https://raw.github.com/marcuswhybrow/minecraft-server-manager/latest/bash_completion/msm -O /etc/bash_completion.d/msm

Those two commands will retrive the MSM script, and the accompanying bash completion script, putting them in the recommended locations.


Check Your Version
------------------

    sudo msm version

The version command will return the current MSM version you have installed, useful for checking the upgrade was successful.