UPDATE_URL="https://raw.github.com/marcuswhybrow/minecraft-server-manager/master"
curl -L "${UPDATE_URL}/installers/common.sh" -o /tmp/msmcommon.sh  #wget isn't installed on Arch by default
source /tmp/msmcommon.sh && rm -f /tmp/msmcommon.sh

function update_system_packages() {
    install_log "Updating sources"
    sudo pacman -Syy || install_error "Couldn't update packages"
}

function install_dependencies() {
    install_log "Installing required packages"
    sudo pacman --noconfirm -S screen rsync zip wget || install_error "Couldn't install dependencies"
}

function enable_init() {
    install_log "Installing systemd service unit"
    sudo wget ${UPDATE_URL}/init/msm.service \
        -O /etc/systemd/system/msm.service
    
    install_log "Enabling automatic startup and shutdown"
    sudo systemctl enable msm.service
}

# Verifies existence and permissions of msm server directory (default /opt/msm)
function create_msm_directories() {
    install_log "Creating MSM directories"
    if [ ! -d "$msm_dir" ]; then
        sudo mkdir -p "$msm_dir" || install_error "Couldn't create directory '$msm_dir'"
    fi
    sudo chown -R $msm_user:$msm_user "$msm_dir" || install_error "Couldn't change file ownership for '$msm_dir'"
    
    if [ ! -d "/etc/init.d" ]; then
        sudo mkdir -p "/etc/init.d/" || install_error "Couldn't create directory '/etc/init.d'"
    fi
}

install_msm
