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


# Manager Tests
# -------------

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
	assertTrue "Server was not created." '[ -d "$SERVER_STORAGE_PATH/example" ]'
}

test_create_server_with_jar_groups() {
	$SCRIPT jargroup create minecraft "https://s3.amazonaws.com/MinecraftDownload/launcher/minecraft_server.jar" > /dev/null
	$SCRIPT server create example > /dev/null
	
	assertTrue "Server was not created." "[ -d \"$SERVER_STORAGE_PATH/example\" ]"
	assertTrue "Server jar was not linked." "[ -f \"$SERVER_STORAGE_PATH/example/$DEFAULT_JAR\" ]"
}


# Server Tests
# ------------




# Jargroup Tests
# --------------




# Perform tests
source shunit2