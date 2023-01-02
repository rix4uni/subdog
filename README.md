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
Requirements: `anew`
```
go install -v github.com/tomnomnom/anew@latest
go install github.com/tomnomnom/unfurl@latest
```

```
mkdir -p ~/tools && cd ~/tools && git clone https://github.com/rix4uni/SubDog.git && cd SubDog && chmod +x subdog
```

### Create shortcut
```
echo -e "\nalias subdog='~/tools/SubDog/subdog'" >> ~/.bashrc
source ~/.bashrc
```

### Usage

**Scan a single domain**
```
subdog -d example.com
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
