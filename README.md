## subdog

subdog collects subdomains from number of different sources to create a list of root subdomains (i.e.: corp.example.com)
- Note: subdog tool gives unfiltered data it can be appex domain or subdomains. you have to filter manually for e.g, `grep -aE "^(.*\.)?domain\.com$"`

## Installation
```
go install github.com/rix4uni/subdog@latest
```

## Download prebuilt binaries
```
wget https://github.com/rix4uni/subdog/releases/download/v0.0.3/subdog-linux-amd64-0.0.3.tgz
tar -xvzf subdog-linux-amd64-0.0.3.tgz
rm -rf subdog-linux-amd64-0.0.3.tgz
mv subdog ~/go/bin/subdog
```
Or download [binary release](https://github.com/rix4uni/subdog/releases) for your platform.

## Compile from source
```
git clone --depth 1 github.com/rix4uni/subdog.git
cd subdog; go install
```

## Usage
```
Usage of subdog:
  -silent
        silent mode.
  -tools string
        Choose tools: subdomaincenter, jldc, virustotal, alienvault, urlscan, certspotter, hackertarget, crtsh, trickest, subdomainfinder, chaos, merklemap, shodan, reverseipdomain, or all (default "all")
  -verbose
        enable verbose mode
  -version
        Print the version of the tool and exit.
```

## Examples Usages

### Single Domain:
1. Run with `subdomaincenter` only:
```
echo "dell.com" | subdog -tools subdomaincenter
```

2. Run with multiple tools comma-seprated `subdomaincenter` and `jldc`:
```
echo "dell.com" | subdog -tools subdomaincenter,virustotal
```

3. Run with all tools, (default):
```
echo "dell.com" | subdog -tools all
```

### Multiple Domains:
1. Run with `subdomaincenter` only:
```
cat wildcards.txt | subdog -tools subdomaincenter
```

2. Run with multiple tools comma-seprated `subdomaincenter` and `jldc`:
```
cat wildcards.txt | subdog -tools subdomaincenter,virustotal
```

3. Run with all tools, (default):
```
cat wildcards.txt | subdog -tools all
```

## Sources 
- [subdomainfinder](https://subdomainfinder.c99.nl)
- [trickest](https://github.com/trickest/inventory)
- [crt.sh](https://crt.sh)
- [hackertarget](https://api.hackertarget.com)
- [dnsdumpster](https://dnsdumpster.com)
- [certspotter](https://api.certspotter.com)
- [urlscan](https://urlscan.io)
- [alienvault](https://otx.alienvault.com)
- [virustotal](https://www.virustotal.com)
- [jldc](https://jldc.me)
- [subdomaincenter](https://api.subdomain.center)
- [chaos](https://chaos.projectdiscovery.io)
- [merklemap](https://api.merklemap.com)
- [shodan](https://api.shodan.io)
- [reverseipdomain](https://sub-scan-api.reverseipdomain.com)