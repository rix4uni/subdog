package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// FetchSubdomainsShodan fetches subdomains for a given domain from the Shodan
func FetchSubdomainsShodan(domain string) ([]string, error) {
	url := fmt.Sprintf("https://api.shodan.io/dns/domain/%s", domain)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response struct {
		Subdomains []string `json:"subdomains"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	// Use a map to track unique filtered subdomains
	subdomainMap := make(map[string]bool)
	for _, sub := range response.Subdomains {
		fullSubdomain := fmt.Sprintf("%s.%s", sub, domain)
		if isSubdomainOrDomain(fullSubdomain, domain) {
			normalized := NormalizeSubdomain(fullSubdomain)
			if normalized != "" {
				subdomainMap[normalized] = true
			}
		}
	}

	// Convert map to slice
	fullSubdomains := make([]string, 0, len(subdomainMap))
	for subdomain := range subdomainMap {
		fullSubdomains = append(fullSubdomains, subdomain)
	}

	return fullSubdomains, nil
}
