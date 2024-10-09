package urlscan

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

// FetchSubdomains fetches subdomains for a given domain from the URLScan API
func FetchSubdomains(domain string) ([]string, error) {
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

	var subdomains []string
	for _, result := range usResponse.Results {
		subdomains = append(subdomains, result.Task.Domain)
		subdomains = append(subdomains, result.Task.ApexDomain)
		subdomains = append(subdomains, result.Page.Domain)
		subdomains = append(subdomains, result.Page.ApexDomain)
	}

	return subdomains, nil
}
