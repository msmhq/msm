wget -q https://raw.github.com/marcuswhybrow/minecraft-server-manager/master/installers/common.sh -O /tmp/msmcommon.sh
source /tmp/msmcommon.sh && rm -f /tmp/msmcommon.sh

function update_system_packages() {
    install_log "Updating sources"
    sudo emerge --sync || install_error "Couldn't update packages"
}

function install_dependencies() {
    install_log "Installing required packages"
    sudo emerge --noreplace screen rsync zip || install_error "Couldn't install dependencies"
}

function reload_cron() {
    install_log "Reloading cron service(s)"
	for cron in `ls /etc/init.d/*cron*`; do
		sudo $cron stop
		sudo $cron start
	done
}

function enable_init() {
    install_log "Enabling automatic startup and shutdown"
    sudo rc-update add msm
}

install_msm
