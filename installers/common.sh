msm_dir="/opt/msm"
msm_user="minecraft"
dl_dir="$(mktemp -t msm.XXXX)"

# Outputs an MSM INSTALL log line
function install_log() {
    echo -e "\n\033[1;32mMSM INSTALL: $*\033[m"
}

# Outputs an MSM INSTALL ERROR log line and exits with status code 1
function install_error() {
    echo -e "\n\033[1;37;41mMSM INSTALL ERROR: $*\033[m"
    exit 1
}

### NOTE: all the below functions are overloadable for system-specific installs
### NOTE: some of the below functions MUST be overloaded due to system-specific installs

function config_installation() {
    install_log "Configure installation"
    echo "Install directory ${msm_dir}"
    read msm_dir

    echo "New server user to be created ${msm_user}"
    read msm_user

    read -p "Do you wish to install with these values y/n" answer
   if [ "$answer" == "y" ]; then
    echo "Proceeding with install."
   else
	echo "Cancelling installation"
    exit 0
   fi
}

# Runs a system software update to make sure we're using all fresh packages
function update_system_packages() {
    # OVERLOAD THIS
    install_error "No function definition for update_system_packages"
}

# Installs additional dependencies (screen, rsync, zip, wget) using system package manager
function install_dependencies() {
    # OVERLOAD THIS
    install_error "No function definition for install_dependencies"
}

# Verifies existence of or adds user for Minecraft server (default "minecraft")
function add_minecraft_user() {
    install_log "Creating default user '$msm_user'"
    useradd --system --user-group $msm_user
}

# Verifies existence and permissions of msm server directory (default /opt/msm)
function create_msm_directories() {
    install_log "Creating MSM directories"
    if [ ! -d "$msm_dir" ]; then
        mkdir -p "$msm_dir" || install_error "Couldn't create directory '$msm_dir'"
    fi
    chown -R $msm_user:$msm_user "$msm_dir" || install_error "Couldn't change file ownership for '$msm_dir'"
}

# Fetches latest msm.conf, cron job, and init script
function download_latest_files() {
    if [ ! -d "$dl_dir" ]; then
        install_error "Temporary download directory was not created properly"
    fi

    install_log "Downloading latest MSM configuration file"
    wget https://raw.github.com/marcuswhybrow/minecraft-server-manager/latest/msm.conf \
        -O "$dl_dir/msm.conf.orig" || install_error "Couldn't download configuration file"

    install_log "Downloading latest MSM cron file"
    wget https://raw.github.com/marcuswhybrow/minecraft-server-manager/latest/cron/msm \
        -O "$dl_dir/msm.cron.orig" || install_error "Couldn't download cron file"

    install_log "Downloading latest MSM version"
    wget https://raw.github.com/marcuswhybrow/minecraft-server-manager/latest/init/msm \
        -O "$dl_dir/msm.init.orig" || install_error "Couldn't download init file"
}

# Patches msm.conf and cron job to use specified username and directory
function patch_latest_files() {
    # patch config file
    install_log "Patching MSM configuration file"
    sed 's#USERNAME="minecraft"#USERNAME="'$msm_user'"#g' "$dl_dir/msm.conf.orig" | \
        sed "s#/opt/msm#$msm_dir#g" >"$dl_dir/msm.conf"

    # patch cron file
    install_log "Patching MSM cron file"
    awk '{ if ($0 !~ /^#/) sub(/minecraft/, "'$msm_user'"); print }' \
        "$dl_dir/msm.cron.orig" >"$dl_dir/msm.cron"

    # patch init file
    install_log "Patching MSM init file"
    cp "$dl_dir/msm.init.orig" "$dl_dir/msm.init"
}

# Installs msm.conf into /etc
function install_config() {
    install_log "Installing MSM configuration file"
    install -b -m0644 "$dl_dir/msm.conf" /etc/msm.conf
    if [ ! -e /etc/msm.conf ]; then
        install_error "Couldn't install configuration file"
    fi
}

# Installs msm.cron into /etc/cron.d
function install_cron() {
    install_log "Installing MSM cron file"
    install "$dl_dir/msm.cron" /etc/cron.d/msm || install_error "Couldn't install cron file"
}

# Reloads cron service (if necessary)
function reload_cron() {
    # OVERLOAD THIS
    install_error "No function defined for reload_cron"
}

# Installs init script into /etc/init.d
function install_init() {
    install_log "Installing MSM init file"
    install -b "$dl_dir/msm.init" /etc/init.d/msm || install_error "Couldn't install init file"

    install_log "Making MSM accessible as the command 'msm'"
    ln -s /etc/init.d/msm /usr/local/bin/msm
}

# Enables init script in default runlevels
function enable_init() {
    # OVERLOAD THIS
    install_error "No function defined for enable_init"
}

# Updates rest of MSM using init script updater
function update_msm() {
    install_log "Asking MSM to update itself"
    /etc/init.d/msm update --noinput
}

function install_complete() {
    install_log "Done. Type 'msm help' to get started. Have fun!"
}
