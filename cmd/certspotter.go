package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type CertspotterResponse []struct {
	DNSNames []string `json:"dns_names"`
}

// FetchDNSNamesCertspotter fetches DNS names for a given domain from the Certspotter API
func FetchDNSNamesCertspotter(domain string) ([]string, error) {
	url := fmt.Sprintf("https://api.certspotter.com/v1/issuances?domain=%s&include_subdomains=true&expand=dns_names", domain)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var csResponse CertspotterResponse
	err = json.Unmarshal(body, &csResponse)
	if err != nil {
		return nil, err
	}

	// Use a map to track unique filtered subdomains
	subdomainMap := make(map[string]bool)

	for _, issuance := range csResponse {
		for _, dnsName := range issuance.DNSNames {
			if isSubdomainOrDomain(dnsName, domain) {
				normalized := NormalizeSubdomain(dnsName)
				if normalized != "" {
					subdomainMap[normalized] = true
				}
			}
		}
	}

	// Convert map to slice
	dnsNames := make([]string, 0, len(subdomainMap))
	for subdomain := range subdomainMap {
		dnsNames = append(dnsNames, subdomain)
	}

	return dnsNames, nil
}
