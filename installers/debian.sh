UPDATE_URL="https://raw.github.com/marcuswhybrow/minecraft-server-manager/latest"
wget -q ${UPDATE_URL}/installers/common.sh -O /tmp/msmcommon.sh
source /tmp/msmcommon.sh && rm -f /tmp/msmcommon.sh

function update_system_packages() {
    install_log "Updating sources"
    sudo apt-get update || install_error "Couldn't update package list"
    sudo apt-get upgrade || install_error "Couldn't upgrade packages"
}

function install_dependencies() {
    install_log "Installing required packages"
    sudo apt-get install screen rsync zip || install_error "Couldn't install dependencies"
}

function reload_cron() {
    install_log "Reloading cron service"
    hash service 2>/dev/null
    if [[ $? == 0 ]]; then
        sudo service cron reload
    else
        sudo /etc/init.d/cron reload
    fi
}

function enable_init() {
    install_log "Enabling automatic startup and shutdown"
    hash insserv 2>/dev/null
    if [[ $? == 0 ]]; then
        sudo insserv msm
    else
        sudo update-rc.d msm defaults
    fi
}

install_msm
