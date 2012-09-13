source <(wget -qO- https://raw.github.com/marcuswhybrow/minecraft-server-manager/master/installers/common.sh)

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
    hash service 2>/dev/null
    if [[ $? == 0 ]]; then
        service cron reload
    else
        /etc/init.d/cron reload
    fi
}

function enable_init() {
    install_log "Enabling automatic startup and shutdown"
    hash insserv 2>/dev/null
    if [[ $? == 0 ]]; then
        insserv msm
    else
        update-rc.d msm defaults
    fi
}

source <(wget -qO- https://raw.github.com/marcuswhybrow/minecraft-server-manager/master/installers/install.sh)
