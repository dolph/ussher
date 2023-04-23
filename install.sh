#!/bin/bash
set -e

# Build a new binary if one does not exist
if ! [ -f ussher ]; then
    build.sh
fi

# Create a user if one does not exist
if ! id -u ussher &>/dev/null; then
    sudo adduser --system --user-group ussher
fi

# Install ussher
sudo install -o root -g ussher -m 0750 ussher /usr/local/bin/ussher

# Create a configuration directory
sudo mkdir --parents /opt/ussher
sudo chown -R ussher:ussher /opt/ussher

# Create a cache directory if one does not exist
sudo mkdir --parents /var/cache
sudo mkdir --parents --mode=0700 /var/cache/ussher
sudo chown ussher:ussher /var/cache/ussher

# Create a log directory if one does not exist
sudo mkdir --parents --mode=0700 /var/log/ussher
sudo chown ussher:ussher /var/log/ussher

# Update sshd configuration
sudo sed -i -E "s~^#?AuthorizedKeysCommand .*~AuthorizedKeysCommand /usr/local/bin/ussher~" /etc/ssh/sshd_config
sudo sed -i -E "s~^#?AuthorizedKeysCommandUser .*~AuthorizedKeysCommandUser ussher~" /etc/ssh/sshd_config
