# MSM version file for Minecraft 1.2.0 and above

set_property LOG_PATH "server.log"

console_event REGEX "^[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} \[.*\]"
console_event START:30 "Done"

console_command WHITELIST_ON "whitelist on"
console_command WHITELIST_OFF "whitelist off"
console_command WHITELIST_ADD "whitelist add <player>"
console_command WHITELIST_REMOVE "whitelist remove <player>"
console_command WHITELIST_LIST "whitelist list"
console_command WHITELIST_RELOAD "whitelist reload"

console_command BLACKLIST_PLAYER_ADD "ban <player>"
console_command BLACKLIST_PLAYER_REMOVE "pardon <player>"
console_command BLACKLIST_PLAYER_LIST "banlist"
console_command BLACKLIST_IP_ADD "ban-ip <address>"
console_command BLACKLIST_IP_REMOVE "pardon-ip <address>"
console_command BLACKLIST_IP_LIST "banlist ips"

console_command OP_ADD "op <player>"
console_command OP_REMOVE "deop <player>"

console_command GAMEMODE "gamemode <player> <mode>" \
  "CONSOLE: Setting <player> to game mode (1|0)" \
  "Can't find user <player>" \
  "<player> already has game mode (1|0)"
console_command KICK "kick <player>" \
  "CONSOLE: Kicking <player>" \
  "Can't find user <player>"
console_command SAY "say <message>"
console_command TOGGLEDOWNFALL "toggledownfall" \
  "CONSOLE: Toggling downfall (on|off) for world" \
  "No world exists with the name"
console_command GIVE "give <player> <item> <amount> <damage>" \
  "CONSOLE: Giving <player> some .+ (.+)" \
  "Can't find user <player>" \
  "There's no item called <item>"
console_command XP "xp <player> <amount>" \
  "CONSOLE: Giving <amount> exp to <player>" \
  "Can't find user <player>" \
  "Invalid exp count: <amount>"
console_command CONNECTED "list" \
  "Connected players:"

console_command TIME_SET "time set <time>" \
  "CONSOLE: Set time to <time>" \
  "Unable to convert time value"
console_command TIME_ADD "time add <time>" \
  "CONSOLE: Added <time> to time" \
  "Unable to convert time value"

console_command SAVE_ON "save-on" \
  "CONSOLE: Enabling level saving.."
console_command SAVE_OFF "save-off" \
  "CONSOLE: Disabling level saving.."
console_command SAVE_ALL:10 "save-all" \
  "CONSOLE: Save complete."
