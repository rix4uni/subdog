## subdog

A powerful subdomain enumeration tool that aggregates data from multiple sources to create comprehensive lists of root subdomains.

## Features

- **Multiple Data Sources**: Aggregates subdomains from 16+ different sources
- **Duplicate Removal**: Automatically removes duplicate subdomains across all sources
- **Output to File**: Save results to a file while still displaying output to terminal
- **Verbose Mode**: Detailed summary table showing counts, timing, and status for each source
- **Parallel Processing**: Run multiple sources simultaneously for faster results
- **Normalization**: Automatically filters out wildcard subdomains and email addresses
- **Flexible Source Selection**: Choose specific sources or exclude unwanted ones

## Installation
```
go install github.com/rix4uni/subdog@latest
```

## Download prebuilt binaries
```
wget https://github.com/rix4uni/subdog/releases/download/v0.0.5/subdog-linux-amd64-0.0.5.tgz
tar -xvzf subdog-linux-amd64-0.0.5.tgz
rm -rf subdog-linux-amd64-0.0.5.tgz
mv subdog ~/go/bin/subdog
```
Or download [binary release](https://github.com/rix4uni/subdog/releases) for your platform.

## Compile from source
```
git clone --depth 1 https://github.com/rix4uni/subdog.git
cd subdog; go install
```

## Usage
```yaml
Usage of subdog:
  -e, --exclude-source string   Comma-separated list of sources to exclude when using --source all
  -l, --list-sources            List all available sources and exit
  -o, --output string           Save subdomain results to a file
  -p, --parallel                Run all sources in parallel to speed up scanning
      --silent                  Silent mode.
  -s, --source string           Choose source(s) to use, or 'all' for all sources. Use --list-sources to see available sources (default "all")
      --verbose                 Enable verbose mode (shows detailed summary table with counts and timing)
      --version                 Print the version of the tool and exit.
```

### Available Sources
- `subdomaincenter` - [Subdomain Center API](https://api.subdomain.center)
- `jldc` - [JLDC Subdomains](https://jldc.me)
- `virustotal` - [VirusTotal API](https://www.virustotal.com)
- `alienvault` - [AlienVault OTX](https://otx.alienvault.com)
- `urlscan` - [URLScan.io](https://urlscan.io)
- `certspotter` - [CertSpotter API](https://api.certspotter.com)
- `hackertarget` - [HackerTarget API](https://api.hackertarget.com)
- `crtsh` - [crt.sh Certificate Search](https://crt.sh)
- `trickest` - [Trickest Inventory](https://github.com/trickest/inventory)
- `subdomainfinder` - [C99 Subdomain Finder](https://subdomainfinder.c99.nl)
- `chaos` - [Chaos Project Discovery](https://chaos.projectdiscovery.io)
- `merklemap` - [MerkleMap API](https://api.merklemap.com)
- `shodan` - [Shodan API](https://api.shodan.io)
- `reverseipdomain` - [Reverse IP Domain](https://sub-scan-api.reverseipdomain.com)
- `dnsdumpster` - [DNS Dumpster](https://dnsdumpster.com)
- `bugbountydata` - [BugBountyData](https://github.com/rix4uni/BugBountyData)

## Usage Examples

### Basic scan with all sources
```yaml
echo "example.com" | subdog
```

### Specific source only
```
echo "target.com" | subdog --source crtsh,certspotter
```

### Exclude specific sources
```
echo "target.com" | subdog --exclude-source shodan,virustotal
```

### Parallel processing for speed
```
echo "target.com" | subdog --silent --parallel
```

### Save results to a file
```
echo "target.com" | subdog --output target.com.txt
```

### Verbose output with summary table
```
echo "target.com" | subdog --verbose
```

### Save results and show verbose summary
```
echo "target.com" | subdog --verbose --output target.com.txt
```

### Scan multiple domains
```yaml
cat subs.txt | subdog
```
