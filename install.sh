#!/usr/bin/env bash

eval "pip3 install -r requirements.txt"

# install some tools
if [[ ! -x "$(command -v jq)" ]]; then
	printf "${BBLUE}Installing jq\n"
	apt install jq -y &>/dev/null;
elif [[ ! -x "$(command -v anew)" ]]; then
	printf "${BBLUE}Installing anew\n"
	go install -v github.com/tomnomnom/anew@latest &>/dev/null;
elif [[ ! -x "$(command -v unfurl)" ]]; then
	printf "${BBLUE}Installing unfurl\n"
	go install github.com/tomnomnom/unfurl@latest &>/dev/null;
fi
