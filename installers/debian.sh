source <(wget -qO- https://raw.github.com/ayeso/minecraft-server-manager/patch-1/installers/common.sh)

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

echo "This script requires superuser access to install files to /etc."
echo "You will be prompted for your password by sudo."

# Clear existing sudo credentials
sudo -k

# Run install functions as sudo
# These functions are defined in another script
# after which this script is sourced
sudo sh <<SCRIPT

    config_installation
    update_system_packages
    install_dependencies
    add_minecraft_user
    create_msm_directories
    download_latest_files
    patch_latest_files
    install_config
    install_cron
    reload_cron
    install_init
    enable_init
    update_msm
    install_complete

SCRIPT
