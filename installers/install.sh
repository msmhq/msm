echo "This script requires superuser access to install files to /etc."
echo "You will be prompted for your password by sudo."

# Clear existing sudo credentials
sudo -k

# Get platform-specific install functions from command arguments
install_file="${1:-debian.sh}"

# run script as sudo
sudo sh <<SCRIPT

    source "$install_file"
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
