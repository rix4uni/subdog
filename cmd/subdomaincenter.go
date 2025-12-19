package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// FetchSubdomainsSubdomaincenter fetches subdomains for a given domain from the API
func FetchSubdomainsSubdomaincenter(domain string) ([]string, error) {
	url := fmt.Sprintf("https://api.subdomain.center/?domain=%s", domain)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rawSubdomains []string
	err = json.Unmarshal(body, &rawSubdomains)
	if err != nil {
		return nil, err
	}

	// Use a map to track unique filtered subdomains
	subdomainMap := make(map[string]bool)
	for _, subdomain := range rawSubdomains {
		if isSubdomainOrDomain(subdomain, domain) {
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
