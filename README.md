# Description
subdog collect number of different sources to create a list of root subdomains (i.e.: corp.example.com)

# Prerequisite
> **Note**: before installing SubDog, make sure to install <a href="https://go.dev/dl/">`Go Language`</a>

# Installation

```
mkdir -p ~/tools
cd ~/tools
git clone https://github.com/rix4uni/SubDog.git
cd SubDog
chmod +x install.sh subdog
./install.sh
```

# Example Usages

Single URL:
```
subdog -d example.com
```

Multiple URLs:
```
subdog -l wildcards.txt
```

# Usage
```
                 ___  __  __  ____  ____  _____  ___
                / __)(  )(  )(  _ \(  _ \(  _  )/ __)
                \__ \ )(__)(  ) _ < )(_) ))(_)(( (_-.
                (___/(______)(____/(____/(_____)\___/
                SubDog v1.1                    coded by @rix4uni in INDIA
TARGET OPTIONS
   -d, --domain            Single Target domain (domain.com)
   -l, --list              Multiple Target domain (wildcards.txt)

USAGE EXAMPLES
   subdog -d example.com
   subdog -l wildcards.txt
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
