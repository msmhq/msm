#!/bin/sh
#
# NOTES: 
# - uses the standard FreeBSD Bourne shell, meaning +POSIX and -bash
# - does not attempt to load common.sh, due to heavy Linux assumptions
# - user creation is done automatically by the FreeBSD port of minecraft-server

install_error() {
	echo -e "\n\033[1;37;41mMSM INSTALL ERROR: $*\033[m"
	exit 1
}

update_system_packages() {
	install_log "Updating sources"
	sudo pkg update || install_error "Couldn't update package list"
}

install_minecraft()
{
	echo "installing java"
	pkg install openjdk8 || exit

	echo "installing minecraft deps"
	pkg install tmux dialog4ports screen rsync || exit

	echo "installing minecraft"
	if grep -qs minecraft /etc/make.conf;
	then
		# already configured
	else
		echo "games_minecraft-server_SET=DAEMON" | tee -a /etc/make.conf
	fi
	export BATCH=${BATCH:="1"}
	make -C /usr/ports/games/minecraft-server install clean
}

configure_minecraft()
{
	echo "configuring minecraft"

	service minecraft onestart
	local _eula="$STAGE_MNT/usr/local/etc/minecraft-server/eula.txt"
	until [ -f "$_eula" ]; do	
		echo "waiting for $_eula to appear"
		sleep 1
	done
	echo "accepting EULA"
	sed -i .bak -e '/^eula/ s/false/true/' "$_eula"
	echo "done"
}

start_minecraft()
{
	echo "starting minecraft"
	sysrc minecraft_enable=YES
	service minecraft start
	sleep 3
}

test_minecraft()
{
	echo "testing minecraft"
	sockstat -l -4 | grep :25565 || exit
	echo "it worked"
}

install_minecraft
configure_minecraft
start_minecraft
test_minecraft
