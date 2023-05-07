#!/usr/bin/env bash

eval "pip3 install -r requirements.txt"

# install some tools
apt install jq -y
go install -v github.com/tomnomnom/anew@latest
go install github.com/tomnomnom/unfurl@latest

# setup
mkdir ~/bin
mv SubDog ~/bin
chmod +x ~/bin/SubDog/*
if ! grep -qxF 'export PATH="$HOME/bin/SubDog:$PATH"' ~/.bashrc ; then echo -e '\nexport PATH="$HOME/bin/SubDog:$PATH"' >> ~/.bashrc ; fi