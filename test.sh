#!/bin/bash

DIR="${MSM_DIR:-$(pwd)}"
DEFUALT_CONF="${MSM_DEFAULT_CONF:-${DIR}/msm.conf}"
TESTS_DIR="${MSM_TESTS_DIR:-${DIR}/tests}"
TMP_DIR="/tmp/msmtest"

# Exit codes
declare -r EX_OK=0
declare -r EX_INVALID_USER=64
declare -r EX_INVALID_COMMAND=65
declare -r EX_INVALID_ARGUMENT=66
declare -r EX_SERVER_STOPPED=67
declare -r EX_SERVER_RUNNING=68
declare -r EX_NAME_NOT_FOUND=69
declare -r EX_FILE_NOT_FOUND=70
declare -r EX_DUPLICATE_NAME=71
declare -r EX_LOGS_NOT_ROLLED=72
declare -r EX_CONF_ERROR=73

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
	
	source "$MSM_CONF"
	
	# Variables accessible by all tests, which are set by the stdall, stderr,
	# stdout and quiet utility functions.
	declare OUTPUT
	declare EXIT_CODE
}

tearDown() {
	# Clear the variables used in test functions
	unset EXIT_CODE
	unset OUTPUT
	
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
	expect_stderr_empty $SCRIPT server list
	assertEquals "Incorrect exit code." $EX_OK $EXIT_CODE
}

### "msm server create" tests

test_reserved_server_names() {
	for name in "start" "stop" "restart" "server" "version" "jargroup" "all"; do
		expect_stderr $SCRIPT server create $name
		assertEquals "Incorrect exit code when creating server name \"$name\"." $EX_INVALID_ARGUMENT $EXIT_CODE
		assertFalse "Server \"$name\" directory was created when it should not have been." "[ -d \"$SERVER_STORAGE_PATH/$name\" ]"
	done
}

test_common_invalid_server_names() {
	expect_stderr $SCRIPT server create "name with spaces"
	assertEquals "Incorrect exit code when creating server name \"name with spaces\"." $EX_INVALID_ARGUMENT $EXIT_CODE
	assertFalse "Server \"name with spaces\" directory was created when it should not have been." "[ -d \"$SERVER_STORAGE_PATH/$name\" ]"
}

test_valid_edge_case_server_names() {
	for name in "serverstart" "CapitalLetters" "0987654321" "name-with-dashes" "name_with_underscores" "Combination-of_different1Things2"; do
		expect_stderr_empty $SCRIPT server create $name
		assertEquals "Incorrect exit code when creating server name \"$name\"." $EX_OK $EXIT_CODE
		assertTrue "Server \"$name\" directory was NOT created when it should not have been." "[ -d \"$SERVER_STORAGE_PATH/$name\" ]"
	done
}

test_create_server_without_any_jargroups() {
	expect_stderr_empty $SCRIPT server create example
	assertEquals "Incorrect exit code." $EX_OK $EXIT_CODE
	assertTrue "Server was not created." "[ -d \"$SERVER_STORAGE_PATH/example\" ]"
}

# Assumes: test_create_server_without_any_jargroups
test_creating_server_when_that_name_already_exists() {
	# Create server "example"
	quiet $SCRIPT server create example
	# Create another server called "example", should be prevented
	expect_stderr $SCRIPT server create example
	
	assertEquals "Incorrect exit code." $EX_DUPLICATE_NAME $EXIT_CODE
}

# Assumes: test_creating_jargroup
test_create_server_with_jar_groups() {
	# Create the "minecraft" jar group, which is used by default when creating
	# new servers.
	quiet $SCRIPT jargroup create minecraft "https://s3.amazonaws.com/MinecraftDownload/launcher/minecraft_server.jar"
	# Create a new server that will use the "minecraft" jar group.
	expect_stderr_empty $SCRIPT server create example
	
	assertEquals "Incorrect exit code." $EX_OK $EXIT_CODE
	assertTrue "Server direcotry was not created." "[ -d \"$SERVER_STORAGE_PATH/example\" ]"
	assertTrue "Server jar was not linked." "[ -f \"$SERVER_STORAGE_PATH/example/$DEFAULT_JAR\" ]"
}

### "msm server delete" tests

test_deleting_server_that_does_not_exist() {
	expect_stderr $SCRIPT server delete example
	
	assertEquals "Incorrect exit code." $EX_NAME_NOT_FOUND $EXIT_CODE
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

# test_creating_jargroup() {
# 	
# }

### "msm jargroup delete" test

### "msm jargroup rename" test

### "msm jargroup changetarget" test

### "msm jargroup getlatest" test


# Perform tests
source shunit2