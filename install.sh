#!/bin/bash
set -e

# Download the latest binary if one does not exist locally
if ! [ -f ussher ]; then
    curl -L -o https://github.com/dolph/ussher/releases/latest/download/ussher
fi

# Create a user if one does not exist
if ! id -u ussher &>/dev/null; then
    sudo adduser --system --user-group ussher
fi

# Install ussher
sudo install -o root -g ussher -m 0750 ussher /usr/local/bin/ussher

# Create a configuration directory
sudo mkdir --parents /etc/ussher

# Create a cache directory if one does not exist
sudo mkdir --parents /var/cache
sudo mkdir --parents --mode=0700 /var/cache/ussher
sudo chown ussher:ussher /var/cache/ussher

# Create a log directory if one does not exist
sudo mkdir --parents --mode=0700 /var/log/ussher
sudo chown ussher:ussher /var/log/ussher

# Update sshd configuration, validate, and apply
sudo sed -i -E "s~^#?AuthorizedKeysCommand .*~AuthorizedKeysCommand /usr/local/bin/ussher~" /etc/ssh/sshd_config
sudo sed -i -E "s~^#?AuthorizedKeysCommandUser .*~AuthorizedKeysCommandUser ussher~" /etc/ssh/sshd_config
sudo sshd -t
sudo systemctl restart sshd
