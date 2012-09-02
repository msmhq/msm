source common.sh

function update_system_packages() {
    install_log "Updating sources"
    apt-get update || install_error "Couldn't update package list"
    apt-get upgrade || install_error "Couldn't upgrade packages"
}

function install_dependencies() {
    install_log "Installing required packages"
    apt-get install screen rsync zip || install_error "Couldn't install dependencies"
}

function reload_cron() {
    install_log "Reloading cron service"
    if [ -x $(which service) ]; then
        service cron reload
    else
        /etc/init.d/cron reload
    fi
}

function enable_init() {
    install_log "Enabling automatic startup and shutdown"
    update-rc.d msm defaults
}
