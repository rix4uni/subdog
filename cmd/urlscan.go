package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type URLScanResponse struct {
	Results []struct {
		Task struct {
			Domain     string `json:"domain"`
			ApexDomain string `json:"apexDomain"`
		} `json:"task"`
		Page struct {
			Domain     string `json:"domain"`
			ApexDomain string `json:"apexDomain"`
		} `json:"page"`
	} `json:"results"`
}

// FetchSubdomainsURLScan fetches subdomains for a given domain from the URLScan API
func FetchSubdomainsURLScan(domain string) ([]string, error) {
	url := fmt.Sprintf("https://urlscan.io/api/v1/search/?q=domain:%s", domain)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var usResponse URLScanResponse
	err = json.Unmarshal(body, &usResponse)
	if err != nil {
		return nil, err
	}

	// Use a map to track unique subdomains
	subdomainMap := make(map[string]bool)

	for _, result := range usResponse.Results {
		// Check Task.Domain (skip ApexDomain as it often contains unrelated domains)
		if isSubdomainOrDomain(result.Task.Domain, domain) {
			normalized := NormalizeSubdomain(result.Task.Domain)
			if normalized != "" {
				subdomainMap[normalized] = true
			}
		}
		// Check Page.Domain (skip ApexDomain as it often contains unrelated domains)
		if isSubdomainOrDomain(result.Page.Domain, domain) {
			normalized := NormalizeSubdomain(result.Page.Domain)
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
