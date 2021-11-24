# MSM version file for Minecraft 1.3.0 and above

extends "minecraft/1.2.0"

console_command GAMEMODE "gamemode <mode> <player>" \
  "That player cannot be found" \
  "The number you have entered (<mode>) is too big, it must be at most 2" \
  "The number you have entered (<mode>) is too small, it must be at least 0" \
  "'<mode>' is not a valid number"

console_command KICK "kick <player>" \
  "Kicked <player> from the game" \
  "That player cannot be found"

console_command TIME_SET "time set <time>" \
  "Set the time to <time>" \
  "The number you have entered (<time>) is too small, it must be at least 0" \
  "'<time>' is not a valid number"
console_command TIME_ADD "time add <time>" \
  "Added <time> to the time" \
  "The number you have entered (<time>) is too small, it must be at least 0" \
  "'<time>' is not a valid number"

console_command TOGGLEDOWNFALL "toggledownfall" \
  "Toggled downfall"

console_command GIVE "give <player> <item> <amount> <damage>" \
  "Given .+ (ID <item>) \* <amount> to <player>" \
  "That player cannot be found" \
  "'<item>' is not a valid number"

console_command XP "xp <amount> <player>" \
  "Given <amount> experience to <player>" \
  "That player cannot be found" \
  "'<amount>' is not a valid number" \
  "The number you have entered (<amount>) is too small, it must be at least 0"

console_command CONNECTED "list" \
  "([^\ ]*)?(, [^\ ]+)*$" \
  "$"

console_command SAVE_ON "save-on" \
  "Turned on world auto-saving"
console_command SAVE_OFF "save-off" \
  "Turned off world auto-saving"
console_command SAVE_ALL:10 "save-all" \
  "Saved the world"
