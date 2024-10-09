package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"github.com/rix4uni/subdog/subdomaincenter"
	"github.com/rix4uni/subdog/jldc"
	"github.com/rix4uni/subdog/virustotal"
	"github.com/rix4uni/subdog/alienvault"
	"github.com/rix4uni/subdog/urlscan"
	"github.com/rix4uni/subdog/certspotter"
	"github.com/rix4uni/subdog/hackertarget"
	"github.com/rix4uni/subdog/crtsh"
	"github.com/rix4uni/subdog/trickest"
	"github.com/rix4uni/subdog/subdomainfinder"
	"github.com/rix4uni/subdog/chaos"
	"github.com/rix4uni/subdog/merklemap"
)

// prints the version message
const version = "0.0.2"

func printVersion() {
	fmt.Printf("Current subdog version %s\n", version)
}

// Prints the Colorful banner
func printBanner() {
	banner := `
                 __         __             
   _____ __  __ / /_   ____/ /____   ____ _
  / ___// / / // __ \ / __  // __ \ / __  /
 (__  )/ /_/ // /_/ // /_/ // /_/ // /_/ / 
/____/ \__,_//_.___/ \__,_/ \____/ \__, /  
                                  /____/ 
`
fmt.Printf("%s\n%50s\n\n", banner, "Current subdog version "+version)
}

func main() {
	tools := flag.String("tools", "all", "Choose tools: subdomaincenter, jldc, virustotal, alienvault, urlscan, certspotter, hackertarget, crtsh, trickest, subdomainfinder, chaos, merklemap, or all")
	silent := flag.Bool("silent", false, "silent mode.")
	version := flag.Bool("version", false, "Print the version of the tool and exit.")
	verbose := flag.Bool("verbose", false, "enable verbose mode")
	flag.Parse()

	// Print version and exit if -version flag is provided
	if *version {
		printBanner()
		printVersion()
		return
	}

	// Don't Print banner if -silnet flag is provided
	if !*silent {
		printBanner()
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		domain := strings.TrimSpace(scanner.Text())

		if *tools == "subdomaincenter" || *tools == "all" {
			if *verbose {
				fmt.Printf("Fetching from subdomaincenter for %s\n", domain)
			}
			subdomains, err := subdomaincenter.FetchSubdomains(domain)
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

		if *tools == "jldc" || *tools == "all" {
			if *verbose {
				fmt.Printf("Fetching from jldc for %s\n", domain)
			}
			subdomains, err := jldc.FetchSubdomains(domain)
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

		if *tools == "virustotal" || *tools == "all" {
			if *verbose {
				fmt.Printf("Fetching from VirusTotal for %s\n", domain)
			}
			subdomains, err := virustotal.FetchSubdomains(domain)
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

		if *tools == "alienvault" || *tools == "all" {
			if *verbose {
				fmt.Printf("Fetching from AlienVault for %s\n", domain)
			}
			subdomains, err := alienvault.FetchSubdomains(domain)
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

		if *tools == "urlscan" || *tools == "all" {
			if *verbose {
				fmt.Printf("Fetching from URLScan for %s\n", domain)
			}
			subdomains, err := urlscan.FetchSubdomains(domain)
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

		if *tools == "certspotter" || *tools == "all" {
			if *verbose {
				fmt.Printf("Fetching from Certspotter for %s\n", domain)
			}
			dnsNames, err := certspotter.FetchDNSNames(domain)
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

		if *tools == "hackertarget" || *tools == "all" {
			if *verbose {
				fmt.Printf("Fetching from HackerTarget for %s\n", domain)
			}
			subdomains, err := hackertarget.FetchSubdomains(domain)
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

		if *tools == "crtsh" || *tools == "all" {
			if *verbose {
				fmt.Printf("Fetching from crt.sh for %s\n", domain)
			}
			subdomains, err := crtsh.FetchSubdomains(domain)
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

		if *tools == "trickest" || *tools == "all" {
			if *verbose {
				fmt.Printf("Fetching from Trickest for %s\n", domain)
			}
			hostnames, err := trickest.FetchHostnames(domain)
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

		if *tools == "subdomainfinder" || *tools == "all" {
			if *verbose {
				fmt.Printf("Fetching from Subdomain Finder for %s\n", domain)
			}
			subdomains, err := subdomainfinder.FetchSubdomains(domain)
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

		if *tools == "chaos" || *tools == "all" {
			if *verbose {
				fmt.Printf("Fetching from Chaos for %s\n", domain)
			}
			chaos.ProcessDomain(domain)
		}

		if *tools == "merklemap" || *tools == "all" {
			if *verbose {
				fmt.Printf("Fetching from MerkleMap for %s\n", domain)
			}
			subdomains, err := merklemap.FetchDomains(domain)
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
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading input: %v\n", err)
	}
}
