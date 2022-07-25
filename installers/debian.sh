#!/bin/bash

UPDATE_URL="https://raw.githubusercontent.com/msmhq/msm/master"
wget -q ${UPDATE_URL}/installers/common.sh -O /tmp/msmcommon.sh
source /tmp/msmcommon.sh && rm -f /tmp/msmcommon.sh

# Installs sudo if it is not installed
function install_sudo() {
    install_log "Installing sudo"
    DEBIAN_FRONTEND=noninteractive apt-get -yqq install sudo
}

function update_system_packages() {
    install_log "Updating sources"
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        if [ "$NAME" == 'Ubuntu' ]; then
            if ! command -v add-apt-repository > /dev/null 2>&1; then
                sudo apt-get -yqq install software-properties-common
            fi

            sudo add-apt-repository universe || install_error "Couldn't enable universe repository"
        fi
    fi
    sudo apt-get -yqq update || install_error "Couldn't update package list"
    DEBIAN_FRONTEND=noninteractive sudo apt-get -yqq upgrade || install_error "Couldn't upgrade packages"
}

function install_dependencies() {
    install_log "Installing required packages"
    DEBIAN_FRONTEND=noninteractive sudo apt-get -yqq install jq rsync screen wget zip || install_error "Couldn't install dependencies"
    if ! command -v java > /dev/null 2>&1; then
        install_log "Installing Java (OpenJDK 17)"
        DEBIAN_FRONTEND=noninteractive sudo apt-get -yqq install openjdk-17-jre-headless || install_error "Couldn't install Java"
    fi
}

function enable_init() {
    install_log "Enabling automatic startup and shutdown"
    hash insserv 2>/dev/null
    if [[ $? == 0 ]]; then
        sudo insserv msm
    else
        sudo update-rc.d msm defaults 99 10
    fi
}

install_msm
