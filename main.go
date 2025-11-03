package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/rix4uni/subdog/banner"
	cmd "github.com/rix4uni/subdog/cmd"
	"github.com/spf13/pflag"
)

var availableSources = []string{
	"subdomaincenter",
	"jldc",
	"virustotal",
	"alienvault",
	"urlscan",
	"certspotter",
	"hackertarget",
	"crtsh",
	"trickest",
	"subdomainfinder",
	"chaos",
	"merklemap",
	"shodan",
	"reverseipdomain",
	"dnsdumpster",
}

func main() {
	sources := pflag.StringP("source", "s", "all", "Choose source(s) to use, or 'all' for all sources. Use --list-sources to see available sources")
	excludeSources := pflag.StringP("exclude-source", "e", "", "Comma-separated list of sources to exclude when using --source all")
	listSources := pflag.BoolP("list-sources", "l", false, "List all available sources and exit")
	parallel := pflag.BoolP("parallel", "p", false, "Run all sources in parallel to speed up scanning")
	silent := pflag.Bool("silent", false, "Silent mode.")
	versionFlag := pflag.Bool("version", false, "Print the version of the tool and exit.")
	verbose := pflag.Bool("verbose", false, "enable verbose mode")
	pflag.Parse()

	if *listSources {
		fmt.Println("Available sources:")
		for _, source := range availableSources {
			fmt.Printf("  - %s\n", source)
		}
		fmt.Println("\nUse 'all' to run all sources")
		return
	}

	// Parse excluded sources into a map for quick lookup
	excludedMap := make(map[string]bool)
	if *excludeSources != "" {
		excluded := strings.Split(*excludeSources, ",")
		for _, e := range excluded {
			excludedMap[strings.TrimSpace(e)] = true
		}
	}

	// Helper function to check if a source should be executed
	shouldRun := func(sourceName string) bool {
		// If source is specifically requested, check if it's excluded
		if *sources == sourceName {
			return !excludedMap[sourceName]
		}
		// If "all" is requested, check if this source is excluded
		if *sources == "all" {
			return !excludedMap[sourceName]
		}
		return false
	}

	if *versionFlag {
		banner.PrintBanner()
		banner.PrintVersion()
		return
	}

	if !*silent {
		banner.PrintBanner()
	}

	// Output mutex for thread-safe printing when using --parallel
	var outputMutex sync.Mutex

	// Function to process a single source
	processSource := func(sourceName string, domain string, fetchFunc func(string) ([]string, error)) {
		if !shouldRun(sourceName) {
			return
		}
		if *verbose {
			outputMutex.Lock()
			fmt.Printf("Fetching from %s for %s\n", sourceName, domain)
			outputMutex.Unlock()
		}
		results, err := fetchFunc(domain)
		if err != nil {
			if *verbose {
				outputMutex.Lock()
				fmt.Printf("Error fetching subdomains from %s for %s: %v\n", sourceName, domain, err)
				outputMutex.Unlock()
			}
			return
		}
		outputMutex.Lock()
		for _, result := range results {
			fmt.Println(result)
		}
		outputMutex.Unlock()
	}

	// Function to process chaos source (special case - no return values)
	processChaosSource := func(domain string) {
		if !shouldRun("chaos") {
			return
		}
		if *verbose {
			outputMutex.Lock()
			fmt.Printf("Fetching from Chaos for %s\n", domain)
			outputMutex.Unlock()
		}
		outputMutex.Lock()
		cmd.ProcessDomainChaos(domain)
		outputMutex.Unlock()
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		domain := strings.TrimSpace(scanner.Text())

		if *parallel {
			// Run all sources in parallel
			var wg sync.WaitGroup

			wg.Add(1)
			go func() {
				defer wg.Done()
				processSource("subdomaincenter", domain, cmd.FetchSubdomainsSubdomaincenter)
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				processSource("jldc", domain, cmd.FetchSubdomainsJldc)
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				processSource("virustotal", domain, cmd.FetchSubdomainsVirusTotal)
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				processSource("alienvault", domain, cmd.FetchSubdomainsAlienVault)
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				processSource("urlscan", domain, cmd.FetchSubdomainsURLScan)
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				processSource("certspotter", domain, func(d string) ([]string, error) {
					return cmd.FetchDNSNamesCertspotter(d)
				})
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				processSource("hackertarget", domain, cmd.FetchSubdomainsHackerTarget)
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				processSource("crtsh", domain, cmd.FetchSubdomainsCrtsh)
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				processSource("trickest", domain, func(d string) ([]string, error) {
					return cmd.FetchHostnamesTrickest(d)
				})
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				processSource("subdomainfinder", domain, cmd.FetchSubdomainsSubdomainFinder)
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				processChaosSource(domain)
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				processSource("merklemap", domain, cmd.FetchDomainsMerkleMap)
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				processSource("shodan", domain, cmd.FetchSubdomainsShodan)
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				processSource("reverseipdomain", domain, cmd.FetchSubdomainsReverseIPDomain)
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				processSource("dnsdumpster", domain, cmd.FetchSubdomainsDNSDumpster)
			}()

			wg.Wait()
		} else {
			// Sequential execution
			if shouldRun("subdomaincenter") {
				if *verbose {
					fmt.Printf("Fetching from subdomaincenter for %s\n", domain)
				}
				subdomains, err := cmd.FetchSubdomainsSubdomaincenter(domain)
				if err != nil {
					if *verbose {
						fmt.Printf("Error fetching subdomains from subdomaincenter for %s: %v\n", domain, err)
					}
				} else {
					for _, subdomain := range subdomains {
						fmt.Println(subdomain)
					}
				}
			}

			if shouldRun("jldc") {
				if *verbose {
					fmt.Printf("Fetching from jldc for %s\n", domain)
				}
				subdomains, err := cmd.FetchSubdomainsJldc(domain)
				if err != nil {
					if *verbose {
						fmt.Printf("Error fetching subdomains from jldc for %s: %v\n", domain, err)
					}
				} else {
					for _, subdomain := range subdomains {
						fmt.Println(subdomain)
					}
				}
			}

			if shouldRun("virustotal") {
				if *verbose {
					fmt.Printf("Fetching from VirusTotal for %s\n", domain)
				}
				subdomains, err := cmd.FetchSubdomainsVirusTotal(domain)
				if err != nil {
					if *verbose {
						fmt.Printf("Error fetching subdomains from VirusTotal for %s: %v\n", domain, err)
					}
				} else {
					for _, subdomain := range subdomains {
						fmt.Println(subdomain)
					}
				}
			}

			if shouldRun("alienvault") {
				if *verbose {
					fmt.Printf("Fetching from AlienVault for %s\n", domain)
				}
				subdomains, err := cmd.FetchSubdomainsAlienVault(domain)
				if err != nil {
					if *verbose {
						fmt.Printf("Error fetching subdomains from AlienVault for %s: %v\n", domain, err)
					}
				} else {
					for _, subdomain := range subdomains {
						fmt.Println(subdomain)
					}
				}
			}

			if shouldRun("urlscan") {
				if *verbose {
					fmt.Printf("Fetching from URLScan for %s\n", domain)
				}
				subdomains, err := cmd.FetchSubdomainsURLScan(domain)
				if err != nil {
					if *verbose {
						fmt.Printf("Error fetching subdomains from URLScan for %s: %v\n", domain, err)
					}
				} else {
					for _, subdomain := range subdomains {
						fmt.Println(subdomain)
					}
				}
			}

			if shouldRun("certspotter") {
				if *verbose {
					fmt.Printf("Fetching from Certspotter for %s\n", domain)
				}
				dnsNames, err := cmd.FetchDNSNamesCertspotter(domain)
				if err != nil {
					if *verbose {
						fmt.Printf("Error fetching DNS names from Certspotter for %s: %v\n", domain, err)
					}
				} else {
					for _, dnsName := range dnsNames {
						fmt.Println(dnsName)
					}
				}
			}

			if shouldRun("hackertarget") {
				if *verbose {
					fmt.Printf("Fetching from HackerTarget for %s\n", domain)
				}
				subdomains, err := cmd.FetchSubdomainsHackerTarget(domain)
				if err != nil {
					if *verbose {
						fmt.Printf("Error fetching subdomains from HackerTarget for %s: %v\n", domain, err)
					}
				} else {
					for _, subdomain := range subdomains {
						fmt.Println(subdomain)
					}
				}
			}

			if shouldRun("crtsh") {
				if *verbose {
					fmt.Printf("Fetching from crt.sh for %s\n", domain)
				}
				subdomains, err := cmd.FetchSubdomainsCrtsh(domain)
				if err != nil {
					if *verbose {
						fmt.Printf("Error fetching subdomains from crt.sh for %s: %v\n", domain, err)
					}
				} else {
					for _, subdomain := range subdomains {
						fmt.Println(subdomain)
					}
				}
			}

			if shouldRun("trickest") {
				if *verbose {
					fmt.Printf("Fetching from Trickest for %s\n", domain)
				}
				hostnames, err := cmd.FetchHostnamesTrickest(domain)
				if err != nil {
					if *verbose {
						fmt.Printf("Error fetching hostnames from Trickest for %s: %v\n", domain, err)
					}
				} else {
					for _, hostname := range hostnames {
						fmt.Println(hostname)
					}
				}
			}

			if shouldRun("subdomainfinder") {
				if *verbose {
					fmt.Printf("Fetching from Subdomain Finder for %s\n", domain)
				}
				subdomains, err := cmd.FetchSubdomainsSubdomainFinder(domain)
				if err != nil {
					if *verbose {
						fmt.Printf("Error fetching subdomains from Subdomain Finder for %s: %v\n", domain, err)
					}
				} else {
					for _, subdomain := range subdomains {
						fmt.Println(subdomain)
					}
				}
			}

			if shouldRun("chaos") {
				if *verbose {
					fmt.Printf("Fetching from Chaos for %s\n", domain)
				}
				cmd.ProcessDomainChaos(domain)
			}

			if shouldRun("merklemap") {
				if *verbose {
					fmt.Printf("Fetching from shodan for %s\n", domain)
				}
				subdomains, err := cmd.FetchDomainsMerkleMap(domain)
				if err != nil {
					if *verbose {
						fmt.Printf("Error fetching domains from shodan for %s: %v\n", domain, err)
					}
				} else {
					for _, subdomain := range subdomains {
						fmt.Println(subdomain)
					}
				}
			}

			if shouldRun("shodan") {
				if *verbose {
					fmt.Printf("Fetching from MerkleMap for %s\n", domain)
				}
				subdomains, err := cmd.FetchSubdomainsShodan(domain)
				if err != nil {
					if *verbose {
						fmt.Printf("Error fetching domains from MerkleMap for %s: %v\n", domain, err)
					}
				} else {
					for _, subdomain := range subdomains {
						fmt.Println(subdomain)
					}
				}
			}

			if shouldRun("reverseipdomain") {
				if *verbose {
					fmt.Printf("Fetching from reverseipdomain for %s\n", domain)
				}
				subdomains, err := cmd.FetchSubdomainsReverseIPDomain(domain)
				if err != nil {
					if *verbose {
						fmt.Printf("Error fetching domains from reverseipdomain for %s: %v\n", domain, err)
					}
				} else {
					for _, subdomain := range subdomains {
						fmt.Println(subdomain)
					}
				}
			}

			if shouldRun("dnsdumpster") {
				if *verbose {
					fmt.Printf("Fetching from DNSDumpster for %s\n", domain)
				}
				subdomains, err := cmd.FetchSubdomainsDNSDumpster(domain)
				if err != nil {
					if *verbose {
						fmt.Printf("Error fetching subdomains from DNSDumpster for %s: %v\n", domain, err)
					}
				} else {
					for _, subdomain := range subdomains {
						fmt.Println(subdomain)
					}
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading input: %v\n", err)
	}
}
