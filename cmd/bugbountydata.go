package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// FetchSubdomainsBugBountyData fetches subdomains for a given domain from the BugBountyData GitHub repository
func FetchSubdomainsBugBountyData(domain string) ([]string, error) {
	url := fmt.Sprintf("https://raw.githubusercontent.com/rix4uni/BugBountyData/refs/heads/main/data/%s.txt", domain)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse the response body as plain text (split by newlines)
	lines := strings.Split(string(body), "\n")

	// Use a map to track unique filtered subdomains
	subdomainMap := make(map[string]bool)
	for _, line := range lines {
		subdomain := strings.TrimSpace(line)
		if subdomain != "" && isSubdomainOrDomain(subdomain, domain) {
			normalized := NormalizeSubdomain(subdomain)
			if normalized != "" {
				subdomainMap[normalized] = true
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

