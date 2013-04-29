# MSM version file for BungeeCord 1.0 and above

console_command STOP "end" \
  "Proxy Stopped."

console_command CONNECTED "glist" \
  "([^\ ]*)?(, [^\ ]+)*$" \
  "$"

console_command SAY "alert <message>"
