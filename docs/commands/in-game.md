---
layout: docs
title: In Game Commands · Docs · MSM
tab: docs
docsitem: commands-in-game
---

In Game Commands
================

These commands can also be called "server commands" but effect who can play, and the experience they have. The **server must be running** for these commands to be used.

    <server> whitelist|wl on|off

A whitelist is a set of names defining who is allowed onto a server. If a player attempts to logon to a server with a whitelist enabled, and their name is not on that whitelist, they cannot join and are told why. This command turns a whitelist on or off.

Note the use of the `|` character means "or", for example `creativeserver whitelist on` is acceptable or `creativeserver wl off`, but not `creativeserver wl on off` or `creativeserver whitelist wl on` or `creativeserver whitelist on|off`.

    <server> whitelist|wl list

Lists the players who are on the whitelist.

    <server> whitelist|wl add|remove <player>

Adds a player to, or removes them from the whitelist.

    <server> blacklist|bl list

Lists the players who are blacklisted. A blacklisted player cannot join a server, even if that person is also whitelisted. A difference from whitelisting is that in addition to players, an IP can be blacklisted. This command outputs both blacklisted players and IP address.

    <server> blacklist|bl player add|remove <player>

Adds or removes a player from the blacklist.

    <server> blacklist|bl ip add|remove <ip-address>

Adds or removes an IP address from the blacklist.

    <server> operator|op add|remove <player>

Adds or removes [operator status][op] for a player. Operators are also considered whitelisted.

    <server> operator|op list

Lists the players who are operators.

    <server> gamemode|gm survival|creative <player>

Sets a player's game mode to either [survival][survival] or [creative][creative].

    <server> kick <player>

Forcibly disconnects a player.

    <server> say <message>

Broadcast a message to all players which appears in bright pink, the invocation `creativeserver say "Hello all"` would produce: `[SERVER]: Hello all`.

    <server> time set|add <number>

Set or increment the in-game time of day. The number must be between 0 and 24000.

    <server> toggledownfall

Toggles any rain or snow downfall.

    <server> give <player> <item> [amount] [data]

Gives an item to a player who is online right now. See [this wiki page][data-values] for an exhaustive list of item codes. The `amount` argument defaults to `1` if omitted, and the `data` argument is (as far as in known) used for applying damage to items which degrade with use. This damage value equates to the number of "uses" a tool has undergone. Remember a diamond pickaxe can take more damage than a stone pickaxe.

    <server> xp <player> <amount>

Gives experience to a player who is online right now. The amount can be a positive or negative integer. Negative integers detract from the amount of experience a player has.

    <server> save on|off

Enables or disables "saving" of world changes to file. In-game players can still make changes to the world whilst saving is off, but rebooting the server would mean loosing those changes. Any changes made whilst saving is off are queued, meaning they are saved to file when saving is enabled again.

    <server> save all

Forces the writing of all pending unsaved changes to file.

    <server> cmd <command>

Sends a command to the server, for example: `creativeserver cmd "tp craftysaurus MrWizardBeef"`.

    <server> cmdlog <command>

Like the previous command, sends a string to the server to execute, but in addition outputs the server log to the screen afterwards. You can stop watching the log output by pressing Ctrl+C.


[op]: http://www.minecraftwiki.net/wiki/Operator
[survival]: http://www.minecraftwiki.net/wiki/Survival
[creative]: http://www.minecraftwiki.net/wiki/Creative
[data-values]: http://www.minecraftwiki.net/wiki/Data_values