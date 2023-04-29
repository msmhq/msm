UPDATE_URL="https://raw.githubusercontent.com/msmhq/msm/master"
wget -q ${UPDATE_URL}/installers/common.sh -O /tmp/msmcommon.sh
source /tmp/msmcommon.sh && rm -f /tmp/msmcommon.sh

# Installs sudo if it is not installed
function install_sudo() {
    yum -yq install sudo
}

function update_system_packages() {
    install_log "Updating sources"
    sudo yum -yq update --skip-broken || install_error "Couldn't update packages"
}

function install_dependencies() {
    install_log "Installing required packages"
    sudo yum -yq install jq rsync screen wget zip || install_error "Couldn't install dependencies"
    if ! command -v java > /dev/null 2>&1; then
        install_log "Installing Java (OpenJDK 17)"
        sudo yum -yq install java-17-openjdk || install_error "Couldn't install Java"
    fi
}

function enable_init() {
    install_log "Enabling automatic startup and shutdown"
    sudo chkconfig --add msm
}

install_msm
