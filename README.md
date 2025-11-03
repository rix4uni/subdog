# subdog

A powerful subdomain enumeration tool that aggregates data from multiple sources to create comprehensive lists of root subdomains.

## ⚠️ Important Note

**subdog provides unfiltered data** that may include:
- Apex domains (e.g., `example.com`)
- Subdomains from related domains (e.g., `subdomain.example.org`)
- International variants (e.g., `example.com.tw`)

For targeted results, always use the recommended filtering command shown in the **Recommended Usage** section.

## Installation
```
go install github.com/rix4uni/subdog@latest
```

## Download prebuilt binaries
```
wget https://github.com/rix4uni/subdog/releases/download/v0.0.4/subdog-linux-amd64-0.0.4.tgz
tar -xvzf subdog-linux-amd64-0.0.4.tgz
rm -rf subdog-linux-amd64-0.0.4.tgz
mv subdog ~/go/bin/subdog
```
Or download [binary release](https://github.com/rix4uni/subdog/releases) for your platform.

## Compile from source
```
git clone --depth 1 github.com/rix4uni/subdog.git
cd subdog; go install
```

## Usage
```yaml
Usage of subdog:
  -e, --exclude-source string   Comma-separated list of sources to exclude when using --source all
  -l, --list-sources            List all available sources and exit
  -p, --parallel                Run all sources in parallel to speed up scanning
      --silent                  Silent mode.
  -s, --source string           Choose source(s) to use, or 'all' for all sources. Use --list-sources to see available sources (default "all")
      --verbose                 enable verbose mode
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

## 🎯 Recommended Usage

### Filtered Output for Target Domain
```yaml
echo "example.com" | subdog --silent | grep -aE "^(.*\.)?example\.com$" | unew
```

**Why filtering is necessary:**
Without filtering, you may get results from related domains:
```yaml
www.example-services.com    # Related domain
api.example.com            # Target subdomain
cdn.us.example.com         # Target subdomain
example.com.tw            # International variant
```

The grep filter ensures you only get subdomains directly under your target domain.

## 💡 Examples

### Single Domain Enumeration
```yaml
# Basic scan with all sources
echo "target.com" | subdog

# Specific source only
echo "target.com" | subdog --source crtsh,certspotter

# Exclude specific sources
echo "target.com" | subdog --exclude-source shodan,virustotal

# Parallel processing for speed
echo "target.com" | subdog --silent --parallel

# Verbose output for debugging
echo "target.com" | subdog --verbose
```

### Multiple Domain Enumeration
```yaml
# Scan multiple domains
cat targets.txt | subdog --silent

# Parallel processing with multiple domains
cat targets.txt | subdog --silent --parallel

# Specific sources for multiple domains
cat targets.txt | subdog --source chaos,certspotter
```

### Advanced Usage
```yaml
# Full pipeline with filtering and sorting
echo "company.com" | subdog --parallel --silent | grep -aE "^(.*\.)?company\.com$" | unew | company_subdomains.txt

# Multiple domains with individual filtering
cat domains.txt | while read domain;do echo "$domain" | subdog --silent | grep -aE "^(.*\.)?${domain//./\\.}$" | unew | "${domain}_subdomains.txt";done
```