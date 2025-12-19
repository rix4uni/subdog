package cmd

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// FetchSubdomainsHackerTarget fetches subdomains for a given domain from the HackerTarget API
func FetchSubdomainsHackerTarget(domain string) ([]string, error) {
	url := fmt.Sprintf("https://api.hackertarget.com/hostsearch/?q=%s", domain)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse the CSV response
	reader := csv.NewReader(strings.NewReader(string(body)))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// Use a map to track unique filtered subdomains
	subdomainMap := make(map[string]bool)
	for _, record := range records {
		if len(record) > 0 {
			subdomain := record[0] // First column contains the subdomain
			if isSubdomainOrDomain(subdomain, domain) {
				normalized := NormalizeSubdomain(subdomain)
				if normalized != "" {
					subdomainMap[normalized] = true
				}
			}
		}
	}

	// Convert map to slice
	subdomains := make([]string, 0, len(subdomainMap))
	for subdomain := range subdomainMap {
		subdomains = append(subdomains, subdomain)
	}

	return subdomains, nil
}
