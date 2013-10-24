#!/bin/bash

USERNAME="minecraft"
DIR="${MSM_DIR:-$(pwd)}"
DEFUALT_CONF="${MSM_DEFAULT_CONF:-${DIR}/msm.conf}"
TESTS_DIR="${MSM_TESTS_DIR:-${DIR}/tests}"
TMP_DIR="/tmp/msmtest"
TEST_RAM="256"

# Exit codes
EX_OK=0
EX_INVALID_USER=64
EX_INVALID_COMMAND=65
EX_INVALID_ARGUMENT=66
EX_SERVER_STOPPED=67
EX_SERVER_RUNNING=68
EX_NAME_NOT_FOUND=69
EX_FILE_NOT_FOUND=70
EX_DUPLICATE_NAME=71
EX_LOGS_NOT_ROLLED=72
EX_CONF_ERROR=73

oneTimeSetUp() {
	# Variables used in tests
	TEST_SCRIPT="${MSM_TEST_SCRIPT:-${DIR}/init/msm}"
	export MSM_CONF="${TMP_DIR}/msm.conf"	
}


setUp() {
	# Create the testing conf from the default one
	mkdir -p "$TMP_DIR" && chown "$USERNAME" "$TMP_DIR"
	cp "$DEFUALT_CONF" "$MSM_CONF" && chown "$USERNAME" "$MSM_CONF"

	# Coppy versioning files
	cp -r "${DIR}/versioning" "${TMP_DIR}/versioning"
	
	# Overwrite variables to use for testing purposes
	echo "" >> "$MSM_CONF"
	echo "# Auto appended by test script:" >> "$MSM_CONF"
	echo "SERVER_STORAGE_PATH=\"${TMP_DIR}/servers\"" >> "$MSM_CONF"
	echo "JAR_STORAGE_PATH=\"${TMP_DIR}/jars\"" >> "$MSM_CONF"
	echo "RAMDISK_STORAGE_PATH=\"${TMP_DIR}/fakeramdisk\"" >> "$MSM_CONF"
	echo "VERSIONING_STORAGE_PATH=\"${TMP_DIR}/versioning\"" >> "$MSM_CONF"
	echo "WORLD_ARCHIVE_PATH=\"${TMP_DIR}/archives/worlds\"" >> "$MSM_CONF"
	echo "LOG_ARCHIVE_PATH=\"${TMP_DIR}/archives/logs\"" >> "$MSM_CONF"
	echo "BACKUP_ARCHIVE_PATH=\"${TMP_DIR}/archives/backups\"" >> "$MSM_CONF"
	echo "DEBUG=\"true\"" >> "$MSM_CONF"
	echo "DEFAULT_USERNAME=\"${USERNAME}\"" >> "$MSM_CONF"
	echo "DEFAULT_SCREEN_NAME=\"msmtest-{SERVER_NAME}\"" >> "$MSM_CONF"
	echo "DEFAULT_INVOCATION=\"java -Xmx${TEST_RAM}M -XX:+UseConcMarkSweepGC -XX:+CMSIncrementalPacing -XX:+AggressiveOpts -jar {JAR} nogui\"" >> "$MSM_CONF"

	source $TEST_SCRIPT
	
	# Variables accessible by all tests, which are set by the stdall, stderr,
	# stdout and quiet utility functions.
	declare OUTPUT
	declare EXIT_CODE
}

tearDown() {
	# Clear the variables used in test functions
	unset EXIT_CODE
	unset OUTPUT

	manager_dirty_all
	
	# Remove the temporary testing directory
	if [ -d "${TMP_DIR}" ]; then
		rm -rf "${TMP_DIR}"
	fi
}


# Utils
# -----

reset_vars() {
	unset OUTPUT
	unset EXIT_CODE
}

# Ensure there is output in stderr
expect_stderr() {
	reset_vars
	OUTPUT=$("$@" 2>&1 1>/dev/null)
	EXIT_CODE=$?
	assertNotNull "No stderr given when expected." "$OUTPUT"
}

# Ensure there is NO output in stderr
expect_stderr_empty() {
	reset_vars
	OUTPUT=$("$@" 2>&1 1>/dev/null)
	EXIT_CODE=$?
	assertNull "$OUTPUT" "$OUTPUT"
}

# Ensure there is output in stdout
expect_stdout() {
	reset_vars
	OUTPUT=$("$@")
	EXIT_CODE=$?
	assertNotNull "No stdout given when expected." "$OUTPUT"
}

# Ensure there is NO output in stdout
expect_stdout_empty() {
	reset_vars
	OUTPUT=$("$@")
	EXIT_CODE=$?
	assertNull "$OUTPUT" "$OUTPUT"
}

# Execute a command silently (ignore stdout and stderr)
quiet() {
	reset_vars
	"$@" >& /dev/null
	EXIT_CODE=$?
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

test_listing_no_servers() {
	expect_stderr_empty $TEST_SCRIPT server list
	assertEquals "Incorrect exit code." $EX_OK $EXIT_CODE
}

### "msm server create" tests

test_reserved_server_names() {
	manager_property SERVER_STORAGE_PATH
	for name in "start" "stop" "restart" "server" "version" "jargroup" "all"; do
		expect_stderr $TEST_SCRIPT server create $name
		assertEquals "Incorrect exit code when creating server name \"$name\"." $EX_INVALID_ARGUMENT $EXIT_CODE
		assertFalse "Server \"$name\" directory was created when it should not have been." "[ -d \"$SETTINGS_SERVER_STORAGE_PATH/$name\" ]"
	done
}

test_common_invalid_server_names() {
	manager_property SERVER_STORAGE_PATH
	expect_stderr $TEST_SCRIPT server create "name with spaces"
	assertEquals "Incorrect exit code when creating server name \"name with spaces\"." $EX_INVALID_ARGUMENT $EXIT_CODE
	assertFalse "Server \"name with spaces\" directory was created when it should not have been." "[ -d \"$SETTINGS_SERVER_STORAGE_PATH/name with spaces\" ]"
}

test_valid_edge_case_server_names() {
	manager_property SERVER_STORAGE_PATH
	for name in "serverstart" "CapitalLetters" "0987654321" "name-with-dashes" "name_with_underscores" "Combination-of_different1Things2"; do
		expect_stderr_empty $TEST_SCRIPT server create $name
		assertEquals "Incorrect exit code when creating server name \"$name\"." $EX_OK $EXIT_CODE
		assertTrue "Server \"$name\" directory was NOT created when it should have been." "[ -d \"$SETTINGS_SERVER_STORAGE_PATH/$name\" ]"
	done
}

test_creating_server_without_any_jargroups() {
	manager_property SERVER_STORAGE_PATH
	expect_stderr_empty $TEST_SCRIPT server create example
	assertEquals "Incorrect exit code." $EX_OK $EXIT_CODE
	assertTrue "Server was not created." "[ -d \"$SETTINGS_SERVER_STORAGE_PATH/example\" ]"
}

# Assumes: test_create_server_without_any_jargroups
test_creating_server_when_that_name_already_exists() {
	# Create server "example"
	quiet $TEST_SCRIPT server create example
	# Create another server called "example", should be prevented
	expect_stderr $TEST_SCRIPT server create example
	
	assertEquals "Incorrect exit code." $EX_DUPLICATE_NAME $EXIT_CODE
}

# Assumes: test_creating_jargroup
test_creating_server_with_jar_groups() {
	# Create the "minecraft" jar group, which is used by default when creating
	# new servers.
	quiet $TEST_SCRIPT jargroup create minecraft "https://s3.amazonaws.com/MinecraftDownload/launcher/minecraft_server.jar"
	# Create a new server that will use the "minecraft" jar group.
	expect_stderr_empty $TEST_SCRIPT server create example

	# Script must be sourced again to recognise the new server
	source $TEST_SCRIPT

	manager_property SERVER_STORAGE_PATH
	server_get_id "example"
	local sid="$RETURN"
	server_property "$sid" JAR_PATH
	server_property "$sid" PATH

	assertEquals "Incorrect exit code." $EX_OK $EXIT_CODE
	assertTrue "Server direcotry was not created." "[ -d \"${SERVER_PATH[$sid]}\" ]"
	assertTrue "Server jar (${SERVER_JAR_PATH[$sid]}) was not linked." "[ -f \"${SERVER_JAR_PATH[$sid]}\" ]"
}

### "msm server delete" tests

test_deleting_server_that_does_not_exist() {
	expect_stderr $TEST_SCRIPT server delete example
	assertEquals "Incorrect exit code." $EX_NAME_NOT_FOUND $EXIT_CODE
}

# Assumes: test_creating_server_without_any_jargroups
test_deleting_server_that_exists_and_is_stopped() {
	quiet $TEST_SCRIPT server create example
	expect_stderr_empty $TEST_SCRIPT server delete example <<< "y"
	
	assertEquals "Incorrect exit code." $EX_OK $EXIT_CODE
	assertFalse "Server directory was not removed." "[ -d \"$SERVER_STORAGE_PATH/example\" ]"
}

### "msm server rename" tests

test_renaming_server_that_does_not_exist() {
	expect_stderr $TEST_SCRIPT server rename example example_new_name
	assertEquals "Incorrect exit code." $EX_NAME_NOT_FOUND $EXIT_CODE
}

# Assumes: test_creating_server_without_any_jargroups
test_renaming_server_that_exists_and_is_stopped() {
	quiet $TEST_SCRIPT server create example
	expect_stderr_empty $TEST_SCRIPT server rename example example_new_name

	manager_property SERVER_STORAGE_PATH
	
	assertEquals "Incorrect exit code." $EX_OK $EXIT_CODE
	assertFalse "Original server name directory still exists." "[ -d \"$SETTINGS_SERVER_STORAGE_PATH/example\" ]"
	assertTrue "New server name directory was not created." "[ -d \"$SETTINGS_SERVER_STORAGE_PATH/example_new_name\" ]"
}


# Stopped Server Tests
# --------------------

### "msm <server> stop" tests

test_stopped_server_stop() {
	quiet $TEST_SCRIPT server create stoppedserver
	expect_stdout $TEST_SCRIPT stoppedserver stop
	assertEquals "Incorrect exit code." $EX_OK $EXIT_CODE
}

### "msm <server> status" tests

test_stopped_server_status() {
	quiet $TEST_SCRIPT server create stoppedserver
	expect_stdout $TEST_SCRIPT stoppedserver status
	assertEquals "Incorrect exit code." $EX_OK $EXIT_CODE
}

### "msm <server> connected" tests

test_stopped_server_connected() {
	quiet $TEST_SCRIPT server create stoppedserver
	expect_stdout $TEST_SCRIPT stoppedserver connected
	assertEquals "Incorrect exit code." $EX_OK $EXIT_CODE
}

### "msm <server> worlds list" tests

test_stopped_server_worlds_list_none() {
	quiet $TEST_SCRIPT server create stoppedserver
	expect_stdout $TEST_SCRIPT stoppedserver worlds list
	assertEquals "Incorrect exit code." $EX_OK $EXIT_CODE
}

### "msm <server> worlds load" tests

test_stopped_server_worlds_load_none() {
	quiet $TEST_SCRIPT server create stoppedserver
	expect_stdout $TEST_SCRIPT stoppedserver worlds load
	assertEquals "Incorrect exit code." $EX_OK $EXIT_CODE
}

### "msm <server> worlds ram" tests

test_stopped_server_worlds_ram_name_not_found() {
	quiet $TEST_SCRIPT server create stoppedserver
	expect_stderr $TEST_SCRIPT stoppedserver worlds ram non_existant_world_name
	assertEquals "Incorrect exit code." $EX_NAME_NOT_FOUND $EXIT_CODE
}

### "msm <server> worlds todisk" tests

test_stopped_server_worlds_todisk_none() {
	quiet $TEST_SCRIPT server create stoppedserver
	expect_stdout $TEST_SCRIPT stoppedserver worlds todisk
	assertEquals "Incorrect exit code." $EX_OK $EXIT_CODE
}

### "msm <server> worlds backup" tests

test_stopped_server_worlds_todisk_none() {
	quiet $TEST_SCRIPT server create stoppedserver
	expect_stdout $TEST_SCRIPT stoppedserver worlds backup
	assertEquals "Incorrect exit code." $EX_OK $EXIT_CODE
}

### "msm <server> worlds on" tests

### "msm <server> worlds off" tests

### "msm <server> logroll" tests

test_stopped_server_logroll_empty() {
	quiet $TEST_SCRIPT server create stoppedserver
	expect_stdout $TEST_SCRIPT stoppedserver logroll
	assertEquals "Incorrect exit code." $EX_OK $EXIT_CODE
}

### "msm <server> backup" tests

test_stopped_server_backup() {
	manager_property BACKUP_ARCHIVE_PATH

	quiet $TEST_SCRIPT server create stoppedserver
	expect_stdout $TEST_SCRIPT stoppedserver backup
	assertEquals "Incorrect exit code." $EX_OK $EXIT_CODE
	assertNotNull "Server backup was not created." "find '$SETTINGS_BACKUP_ARCHIVE_PATH/stoppedserver' -mindepth 1 -maxdepth 1 -type f -name '*.zip'"
}

### "msm <server> jar" tests

test_stopped_server_jar() {
	quiet $TEST_SCRIPT server create stoppedserver

	source $TEST_SCRIPT

	server_get_id "stoppedserver"
	local sid="$RETURN"
	
	server_property "$sid" JAR_PATH

	# Remove the existing jar file
	rm -f "$JAR_PATH"

	# Create new jargroup
	quiet $TEST_SCRIPT jargroup create newgroup "https://s3.amazonaws.com/MinecraftDownload/launcher/minecraft_server.jar"

	# Assign jargroup's latest jar to server
	expect_stdout $TEST_SCRIPT stoppedserver jar newgroup

	assertEquals "Incorrect exit code." $EX_OK $EXIT_CODE
	assertTrue "Server jar (${SERVER_JAR_PATH[$sid]}) was not linked." "[ -f \"${SERVER_JAR_PATH[$sid]}\" ]"
}

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

### "msm jargroup create" test

test_reserved_jargroup_names() {
	for name in "start" "stop" "restart" "server" "version" "jargroup" "all"; do
		expect_stderr $TEST_SCRIPT jargroup create "$name" https://s3.amazonaws.com/MinecraftDownload/launcher/minecraft_server.jar
		assertEquals "Incorrect exit code when creating jar group name \"$name\"." $EX_INVALID_ARGUMENT $EXIT_CODE
		assertFalse "Jar group \"$name\" directory was created when it should not have been." "[ -d \"$JAR_STORAGE_PATH/$name\" ]"
	done
}

test_common_invalid_jargroups_names() {
	expect_stderr $TEST_SCRIPT jargroup create "name with spaces" https://s3.amazonaws.com/MinecraftDownload/launcher/minecraft_server.jar
	assertEquals "Incorrect exit code when creating jar group name \"name with spaces\"." $EX_INVALID_ARGUMENT $EXIT_CODE
	assertFalse "Jar group \"name with spaces\" directory was created when it should not have been." "[ -d \"$JAR_STORAGE_PATH/name with spaces\" ]"
}

test_valid_edge_case_jargroup_names() {
	manager_property JAR_STORAGE_PATH
	for name in "serverstart" "CapitalLetters" "0987654321" "name-with-dashes" "name_with_underscores" "Combination-of_different1Things2"; do
		expect_stderr_empty $TEST_SCRIPT jargroup create $name https://s3.amazonaws.com/MinecraftDownload/launcher/minecraft_server.jar
		assertEquals "Incorrect exit code when creating jar group name \"$name\"." $EX_OK $EXIT_CODE
		assertTrue "Jar group \"$name\" directory was NOT created when it should have been." "[ -d \"$SETTINGS_JAR_STORAGE_PATH/$name\" ]"
	done
}

test_creating_jargroup() {
	expect_stderr_empty $TEST_SCRIPT jargroup create minecraft https://s3.amazonaws.com/MinecraftDownload/launcher/minecraft_server.jar
	assertEquals "Incorrect exit code." $EX_OK $EXIT_CODE
}

# Assumes: test_creating_jargroup
test_creating_jargroup_when_that_name_already_exists() {
	quiet $TEST_SCRIPT jargroup create minecraft https://s3.amazonaws.com/MinecraftDownload/launcher/minecraft_server.jar
	expect_stderr $TEST_SCRIPT jargroup create minecraft https://s3.amazonaws.com/MinecraftDownload/launcher/minecraft_server.jar
	assertEquals "Incorrect exit code." $EX_DUPLICATE_NAME $EXIT_CODE
}

### "msm jargroup list" test

test_listing_no_jargroups() {
	expect_stderr_empty $TEST_SCRIPT jargroup list
	assertEquals "Incorrect exit code." $EX_OK $EXIT_CODE
}

# Assumes: test_creating_jargroup
test_listing_one_jargroup() {
	quiet $TEST_SCRIPT jargroup create minecraft https://s3.amazonaws.com/MinecraftDownload/launcher/minecraft_server.jar
	expect_stderr_empty $TEST_SCRIPT jargroup list
	assertEquals "Incorrect exit code." $EX_OK $EXIT_CODE
}

# Assumes: test_creating_jargroup
test_listing_multiple_jargroups() {
	quiet $TEST_SCRIPT jargroup create minecraft https://s3.amazonaws.com/MinecraftDownload/launcher/minecraft_server.jar
	quiet $TEST_SCRIPT jargroup create craftbukkit http://dl.bukkit.org/latest-rb/craftbukkit.jar
	expect_stderr_empty $TEST_SCRIPT jargroup list
	assertEquals "Incorrect exit code." $EX_OK $EXIT_CODE
}

### "msm jargroup delete" test

### "msm jargroup rename" test

### "msm jargroup changetarget" test

### "msm jargroup getlatest" test


# Perform tests
# local copy of shunit2
source shunit2 2>/dev/null
if [ $? -ne 0 ]; then
	# default install location for shunit2
	source "/usr/share/shunit2/shunit2" 2>/dev/null
	if [ $? -ne 0 ]; then
		echo "Could not find shUnit2 install."
	fi
fi
