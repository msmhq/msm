echo "This script requires superuser access to install files to /etc."
echo "You will be prompted for your password by sudo."

# Clear existing sudo credentials
sudo -k

# run script as sudo
sudo sh <<SCRIPT

	echo "\033[1;32mMSM INSTALL: Updating sources\033[m"
	apt-get update

	echo "\n\033[1;32mMSM INSTALL: Installing required packages\033[m"
	apt-get install screen rsync zip

	echo "\n\033[1;32mMSM INSTALL: Creating default user \"minecraft\"\033[m"
	adduser --system minecraft

	echo "\n\033[1;32mMSM INSTALL: Creating MSM directories\033[m"
	mkdir -p /opt/msm
	chown minecraft:minecraft /opt/msm

	echo "\n\033[1;32mMSM INSTALL: Downloading latest MSM configuration file\033[m"
	wget https://raw.github.com/marcuswhybrow/minecraft-server-manager/latest/msm.conf -O /etc/msm.conf

	echo "\n\033[1;32mMSM INSTALL: Downloading latest MSM cron file\033[m"
	wget https://raw.github.com/marcuswhybrow/minecraft-server-manager/latest/cron/msm -O /etc/cron.d/msm

	echo "\n\033[1;32mMSM INSTALL: Downloading latest MSM version\033[m"
	wget https://raw.github.com/marcuswhybrow/minecraft-server-manager/latest/init/msm -O /etc/init.d/msm

	echo "\n\033[1;32mMSM INSTALL: Giving it the necessary permissions\033[m"
	chmod 755 /etc/init.d/msm

	echo "\n\033[1;32mMSM INSTALL: Enabling automatic startup and shutdown\033[m"
	update-rc.d msm defaults

	echo "\n\033[1;32mMSM INSTALL: Making MSM accessible as the command \"msm\"\033[m"
	ln -s /etc/init.d/msm /usr/local/bin/msm

	echo "\n\033[1;32mMSM INSTALL: Asking MSM to update itself\033[m"
	/etc/init.d/msm update --noinput

	echo "\n\033[1;32mMSM INSTALL: Done. Type \"msm help\" to get started. Have fun!\033[m"

SCRIPT