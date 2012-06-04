#!/bin/bash

DIR="${MSM_DIR:-$(pwd)}"
DEFUALT_CONF="${MSM_DEFAULT_CONF:-${DIR}/msm.conf}"
TESTS_DIR="${MSM_TESTS_DIR:-${DIR}/tests}"
TMP_DIR="/tmp/msmtest"

oneTimeSetUp() {
	# Variables used in tests
	SCRIPT="${MSM_SCRIPT:-${DIR}/init/msm}"
	export MSM_CONF="${TMP_DIR}/msm.conf"
	
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

# Manager Tests
# -------------

test_create_server_without_any_jargroups() {
	$SCRIPT server create example
	
	source "$MSM_CONF"
	assertTrue "Server was not created." '[ -d "$SERVER_STORAGE_PATH/example" ]'
}


# Server Tests
# ------------




# Jargroup Tests
# --------------




# Perform tests
source shunit2