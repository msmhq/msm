BASE_URL="https://raw.githubusercontent.com/msmhq/msm/master/installers/"
TMP_FILE="/tmp/msminst.sh"

function check_os() {
    if [[ -f /etc/debian_version ]]; then
        INSTALL_SCRIPT="debian.sh"
    elif [[ -f /etc/redhat-release ]]; then
        INSTALL_SCRIPT="redhat.sh"
    elif [[ -f /etc/arch-release ]]; then
        INSTALL_SCRIPT="arch.sh"
    elif [[ -f /bin/freebsd-version ]]; then
        INSTALL_SCRIPT="freebsd.sh"
    else
        echo "Error, unsupported distribution. Please install manually."
        exit 1
    fi
}

function download_installer() {
    local _dl_uri="${BASE_URL}/${INSTALL_SCRIPT}"

    case "$1" in
        curl)
            curl -L "$_dl_uri" -o "$TMP_FILE"
        ;;
        wget)
            wget -q "$_dl_uri" -O "$TMP_FILE"
        ;;
        fetch)
            fetch -m -o "$TMP_FILE" "$_dl_uri"
        ;;
    esac
}

function get_installer() {

    for util in curl wget fetch; do
        type "$util" 1>/dev/null 2>&1
        if [[ $? -eq 0 ]]; then
            download_installer $util
            return  # we suceeded, exit this function
        fi
    done

    echo "neither curl, wget, nor fetch was found. Install one and try again."
    exit 1
}

function do_install() {
    if [[ -f "$TMP_FILE" ]]; then
        chmod u+x "$TMP_FILE"
        "$TMP_FILE" && rm -f "$TMP_FILE"
    else
        echo "Error, failed to download install script."
        exit 1
    fi
}

check_os
get_installer
do_install
