source <(wget -qO- https://raw.github.com/ayeso/minecraft-server-manager/patch-1/installers/common.sh)
source <(wget -qO- https://raw.github.com/ayeso/minecraft-server-manager/patch-1/installers/$1.sh)

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
