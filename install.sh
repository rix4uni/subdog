#!/usr/bin/env bash

pip3 install -r requirements.txt

# install some tools
apt install jq -y
go install -v github.com/tomnomnom/anew@latest
go install github.com/tomnomnom/unfurl@latest

# setup
mkdir -p ~/bin
if ! grep -qxF 'export PATH="$HOME/bin/SubDog:$PATH"' ~/.bashrc ; then echo -e '\nexport PATH="$HOME/bin/SubDog:$PATH"' >> ~/.bashrc ; fi
