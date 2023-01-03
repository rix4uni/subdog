#!/usr/bin/env bash

echo -e "Checking Golang latest version"
go_latest_version=$(curl -L -s https://golang.org/VERSION?m=text)
go_system_version=$(go version | awk {'print $3'})

if [[ "$go_latest_version" == "$go_system_version" ]]; then
    echo -e "Golang is already installed and updated"
elif [[ "$go_latest_version" != "$go_system_version" ]]; then
    echo -e "Installing Golang latest version"
    eval wget https://dl.google.com/go/${go_latest_version}.linux-amd64.tar.gz
    sudo tar -C /usr/local -xzf ${go_latest_version}.linux-amd64.tar.gz && rm -rf ${go_latest_version}.linux-amd64.tar.gz
    export GOROOT=/usr/local/go
    export GOPATH=$HOME/go
    export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
    if [[ -f ~/.zshrc ]]; then
        echo -e "\n#Golang Variable" >> ~/.zshrc
        echo 'export GOROOT=/usr/local/go' >> ~/.zshrc
        echo 'export GOPATH=$HOME/go' >> ~/.zshrc
        echo 'export PATH=$GOPATH/bin:$GOROOT/bin:$PATH' >> ~/.zshrc
        eval source ~/.zshrc
    elif [[ -f ~/.bashrc ]]; then
        echo -e "\n#Golang Variable" >> ~/.bashrc
        echo 'export GOROOT=/usr/local/go' >> ~/.bashrc
        echo 'export GOPATH=$HOME/go' >> ~/.bashrc
        echo 'export PATH=$GOPATH/bin:$GOROOT/bin:$PATH' >> ~/.bashrc
        eval source ~/.bashrc
    fi
fi

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

# Create shortcut
echo -e "\nalias subdog='~/tools/SubDog/subdog'" >> ~/.bashrc
eval source ~/.bashrc
