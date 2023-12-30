#!/bin/bash

NODE_VERSION="16";
NODE_MAJOR=20

if ! command -V node >/dev/null 2>&1; then
    # curl -sL https://deb.nodesource.com/setup_${NODE_VERSION}.x | sudo -E bash -
    # sudo apt install -q -y  nodejs
    sudo apt-get update
    sudo apt-get install -y ca-certificates curl gnupg
    sudo mkdir -p /etc/apt/keyrings
    curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | sudo gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg
    echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_$NODE_MAJOR.x nodistro main" | sudo tee /etc/apt/sources.list.d/nodesource.list

    sudo apt-get update
    sudo apt-get install nodejs -y

    # Alias
    source "$(dirname $0)/helpers/alias.sh";

    addAlias ~/.profile "node_modules" "export PATH=\$PATH:node_modules/.bin"
fi

# update npm
echo "[v] Update npm to latest...";
npm install -g npm@10.2.5

# dev packages
echo "[v] Installing npm global packages for development...";
sudo npm install -g  yarn npm-check-updates