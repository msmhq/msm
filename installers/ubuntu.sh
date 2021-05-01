wget -q https://raw.github.com/marcuswhybrow/minecraft-server-manager/master/installers/common.sh -O /tmp/msmcommon.sh
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
    sudo service cron reload
}

function enable_init() {
    install_log "Downloading latest upstart config"
    sudo wget https://raw.github.com/marcuswhybrow/minecraft-server-manager/latest/upstart/msm.conf \
        -O "$dl_dir/msm.conf.upstart" || install_error "Couldn't download upstart config file"
    install_log "Enabling automatic startup and shutdown"
    sudo install -b -m0644  "$dl_dir/msm.conf.upstart" /etc/init/msm.conf
    install_log "MSM can now be globally started and stopped via 'start msm' and 'stop msm'."
}

install_msm
