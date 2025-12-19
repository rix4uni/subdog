package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type AlienVaultResponse struct {
	PassiveDNS []struct {
		Hostname string `json:"hostname"`
	} `json:"passive_dns"`
}

// FetchSubdomainsAlienVault fetches subdomains for a given domain from the AlienVault API
func FetchSubdomainsAlienVault(domain string) ([]string, error) {
	url := fmt.Sprintf("https://otx.alienvault.com/api/v1/indicators/domain/%s/passive_dns", domain)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var avResponse AlienVaultResponse
	err = json.Unmarshal(body, &avResponse)
	if err != nil {
		return nil, err
	}

	// Use a map to track unique filtered subdomains
	subdomainMap := make(map[string]bool)
	for _, record := range avResponse.PassiveDNS {
		if isSubdomainOrDomain(record.Hostname, domain) {
			normalized := NormalizeSubdomain(record.Hostname)
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
