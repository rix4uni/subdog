```

███████╗██╗   ██╗██████╗ ██████╗  ██████╗  ██████╗ 
██╔════╝██║   ██║██╔══██╗██╔══██╗██╔═══██╗██╔════╝ 
███████╗██║   ██║██████╔╝██║  ██║██║   ██║██║  ███╗
╚════██║██║   ██║██╔══██╗██║  ██║██║   ██║██║   ██║
███████║╚██████╔╝██████╔╝██████╔╝╚██████╔╝╚██████╔╝
╚══════╝ ╚═════╝ ╚═════╝ ╚═════╝  ╚═════╝  ╚═════╝  

```
       
### Description
subdog collect number of different sources to create a list of root subdomains (i.e.: corp.example.com)                                         

### Installation

```
mkdir -p ~/tools
cd ~/tools
git clone https://github.com/rix4uni/SubDog.git
cd SubDog
chmod +x install.sh subdog
pip3 install -r requirements.txt
./install.sh
```

### Usage

Single URL:
```
subdog -d example.com
```

Multiple URLs:
```
for target in $(cat wildcards.txt);do subdog -d $target;done
```

### Sources 
- [rapiddns](https://rapiddns.io)
- [threatminer](https://api.threatminer.org) 
- [riddler](https://riddler.io)
- [alienvault](https://otx.alienvault.com)
- [WayBackMachine](http://web.archive.org)
- [hackertarget](https://api.hackertarget.com)
- [crt.sh](https://crt.sh)
- [jldc.me](https://jldc.me)
- [urlscan](https://urlscan.io)
- [subdomainfinder](https://subdomainfinder.c99.nl)
