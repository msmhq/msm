#/bin/sh
PATH=/usr/bin
UPDATE_URL="https://raw.githubusercontent.com/msmhq/msm/master"
if [ ! -f "/usr/bin/wget" ]; then
  echo "/usr/bin/wget not found, installing via RPM"
  sudo yum install wget --skip-broken
fi
wget -q ${UPDATE_URL}/installers/common.sh -O /tmp/msmcommon.sh
source /tmp/msmcommon.sh && rm -f /tmp/msmcommon.sh

function update_system_packages() {
    install_log "Updating sources"
    sudo yum update --skip-broken || install_error "Couldn't update packages"
}

function install_dependencies() {
    install_log "Installing required packages"
    sudo yum install screen rsync zip java || install_error "Couldn't install dependencies"
}

function enable_init() {
    install_log "Enabling automatic startup and shutdown"
    sudo chkconfig --add msm
}

install_msm
