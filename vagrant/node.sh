#!/bin/bash

NODE_MAJOR=16;

if ! command -V node >/dev/null 2>&1; then
   
   sudo snap install node --classic --channel=${NODE_MAJOR}
    # curl -sL https://deb.nodesource.com/setup_${NODE_VERSION}.x | sudo -E bash -
    # sudo apt install -q -y  nodejs

    # Alias
    source "$(dirname $0)/helpers/alias.sh";

    addAlias ~/.profile "node_modules" "export PATH=\$PATH:node_modules/.bin"
fi

node --version
npm --version

# dev packages
echo "[v] Installing npm global packages for development...";
sudo npm install -g yarn 

# update npm
echo "[v] Update npm to latest...";
sudo npm install -g npm@latest

node --version
npm --version
yarn --version