wget -q https://raw.github.com/marcuswhybrow/minecraft-server-manager/master/installers/common.sh -O /tmp/msmcommon.sh
source /tmp/msmcommon.sh
rm -f /tmp/msmcommon.sh

function update_system_packages() {
    install_log "Updating sources"
    yum update || install_error "Couldn't update packages"
}

function install_dependencies() {
    install_log "Installing required packages"
    yum install screen rsync zip || install_error "Couldn't install dependencies"
}

function reload_cron() {
    install_log "Reloading cron service"
    service crond reload
}

function enable_init() {
    install_log "Enabling automatic startup and shutdown"
    chkconfig --add msm
}

install_msm
