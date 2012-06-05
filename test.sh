#!/bin/bash

DIR="${MSM_DIR:-$(pwd)}"
DEFUALT_CONF="${MSM_DEFAULT_CONF:-${DIR}/msm.conf}"
TESTS_DIR="${MSM_TESTS_DIR:-${DIR}/tests}"
TMP_DIR="/tmp/msmtest"

oneTimeSetUp() {
	# Variables used in tests
	SCRIPT="${MSM_SCRIPT:-${DIR}/init/msm}"
	export MSM_CONF="${TMP_DIR}/msm.conf"	
}


setUp() {
	source "$DEFUALT_CONF"
	
	# Create the testing conf from the default one
	mkdir -p "$TMP_DIR" && chown "$USERNAME":"$USERNAME" "$TMP_DIR"
	cp "$DEFUALT_CONF" "$MSM_CONF" && chown "$USERNAME":"$USERNAME" "$MSM_CONF"
	
	# Overwrite the directories to use for testin purposes
	echo "" >> "$MSM_CONF"
	echo "# Auto appended by test script:" >> "$MSM_CONF"
	echo "SERVER_STORAGE_PATH=\"${TMP_DIR}/servers\"" >> "$MSM_CONF"
	echo "JAR_STORAGE_PATH=\"${TMP_DIR}/jars\"" >> "$MSM_CONF"
	echo "RAMDISK_STORAGE_PATH=\"${TMP_DIR}/fakeramdisk\"" >> "$MSM_CONF"
	echo "WORLD_ARCHIVE_PATH=\"${TMP_DIR}/archives/worlds\"" >> "$MSM_CONF"
	echo "LOG_ARCHIVE_PATH=\"${TMP_DIR}/archives/logs\"" >> "$MSM_CONF"
	echo "BACKUP_ARCHIVE_PATH=\"${TMP_DIR}/archives/backups\"" >> "$MSM_CONF"
	echo "DEBUG=\"true\"" >> "$MSM_CONF"
}

tearDown() {
	# Remove the temporary testing directory
	if [ -d "${TMP_DIR}" ]; then
		rm -rf "${TMP_DIR}"
	fi
}


# Utils
# -----

stdall() {
	$1 "${@:2}" 2>&1
}


# Global Command Tests
# -----------------------

### "msm start" tests

### "msm stop" tests

### "msm restart" tests

### "msm version" tests



# Server Management Tests
# -----------------------

### "msm server list" tests

### "msm server create" tests

test_reserved_server_names() {
	local result
	local expected_regex="^Invalid\ name"
	source "$MSM_CONF"
	
	for name in "start" "stop" "restart" "server" "version" "jargroup" "all"; do
		result="$(stdall $SCRIPT server create $name)"
		assertTrue "Server name \"$name\" was accepted but should be invalid." "[[ '$result' =~ $expected_regex ]]"
		assertFalse "Server \"$name\" directory was created when it should not have been." "[ -d \"$SERVER_STORAGE_PATH/$name\" ]"
	done
}

test_common_invalid_server_names() {
	local result
	local expected_regex="^Invalid\ name"
	source "$MSM_CONF"
	
	for name in "name with spaces"; do
		result="$(stdall $SCRIPT server create $name)"
		assertFalse "Server \"$name\" directory was created when it should not have been." "[ -d \"$SERVER_STORAGE_PATH/$name\" ]"
	done
}

test_valid_edge_case_server_names() {
	local result
	local expected_regex="^Invalid\ name"
	source "$MSM_CONF"
	
	for name in "serverstart" "CapitalLetters" "0987654321" "name-with-dashes" "name_with_underscores" "Combination-of_different1Things2"; do
		result="$(stdall $SCRIPT server create $name)"
		assertTrue "Server \"$name\" directory was NOT created when it should not have been." "[ -d \"$SERVER_STORAGE_PATH/$name\" ]"
	done
}

test_create_server_without_any_jargroups() {
	$SCRIPT server create example > /dev/null
	
	source "$MSM_CONF"
	assertTrue "Server was not created." "[ -d \"$SERVER_STORAGE_PATH/example\" ]"
}

# Assumes: test_create_server_without_any_jargroups
test_creating_server_when_that_name_already_exists() {
	$SCRIPT server create example > /dev/null
	local result="$(stdall $SCRIPT server create example)"
	
	source "$MSM_CONF"
	assertTrue "Failed to prevent duplicating an existing server name." "[[ \"$result\" == \"A server with that name already exists.\" ]]"
}

test_create_server_with_jar_groups() {
	$SCRIPT jargroup create minecraft "https://s3.amazonaws.com/MinecraftDownload/launcher/minecraft_server.jar" > /dev/null
	$SCRIPT server create example > /dev/null
	
	assertTrue "Server was not created." "[ -d \"$SERVER_STORAGE_PATH/example\" ]"
	assertTrue "Server jar was not linked." "[ -f \"$SERVER_STORAGE_PATH/example/$DEFAULT_JAR\" ]"
}

### "msm server delete" tests

test_deleting_server_that_does_not_exist() {
	local result="$(stdall $SCRIPT server delete example)"
	local regex="^There\ is\ no\ server\ with\ the\ name"
	
	assertTrue "" "[[ \"$result\" =~ $regex ]]"
}

### "msm server rename" tests


# Individual Server Tests
# -----------------------

### "msm <server> start" tests

### "msm <server> stop" tests

### "msm <server> restart" tests

### "msm <server> status" tests

### "msm <server> connected" tests

### "msm <server> worlds list" tests

### "msm <server> worlds load" tests

### "msm <server> worlds ram" tests

### "msm <server> worlds todisk" tests

### "msm <server> worlds backup" tests

### "msm <server> worlds on" tests

### "msm <server> worlds off" tests

### "msm <server> worlds off" tests

### "msm <server> logroll" tests

### "msm <server> backup" tests

### "msm <server> jar" tests

### "msm <server> whitelist on" tests

### "msm <server> whitelist off" tests

### "msm <server> whitelist add" tests

### "msm <server> whitelist remove" tests

### "msm <server> whitelist list" tests

### "msm <server> blacklist player add" tests

### "msm <server> blacklist player remove" tests

### "msm <server> blacklist ip add" tests

### "msm <server> blacklist ip remove" tests

### "msm <server> blacklist list" tests

### "msm <server> operator add" tests

### "msm <server> operator remove" tests

### "msm <server> operator list" tests

### "msm <server> gamemode" tests

### "msm <server> kick" tests

### "msm <server> say" tests

### "msm <server> time set" tests

### "msm <server> time add" tests

### "msm <server> toggledownfall" tests

### "msm <server> save on" tests

### "msm <server> save off" tests

### "msm <server> save all" tests

### "msm <server> cmd" tests

### "msm <server> cmdlog" tests

### "msm <server> console" tests


# Jargroup Tests
# --------------

### "msm jargroup list" test

### "msm jargroup create" test

### "msm jargroup delete" test

### "msm jargroup rename" test

### "msm jargroup changetarget" test

### "msm jargroup getlatest" test


# Perform tests
source shunit2