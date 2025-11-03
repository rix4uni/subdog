package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// TrickestData represents the structure of the Trickest JSON response
type TrickestData struct {
	Domain   string `json:"domain"`
	Hostnames string `json:"hostnames"`
}

// FetchHostnamesTrickest fetches hostnames for a given domain from the Trickest API
func FetchHostnamesTrickest(domain string) ([]string, error) {
	// Fetch the Trickest targets JSON file
	resp, err := http.Get("https://raw.githubusercontent.com/rix4uni/targets-filter/refs/heads/main/trickest-targets.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON data
	var targets []TrickestData
	if err := json.Unmarshal(body, &targets); err != nil {
		return nil, err
	}

	var hostnames []string
	for _, target := range targets {
		if strings.EqualFold(target.Domain, domain) {
			// Fetch hostnames from the URL
			url := target.Hostnames
			resp, err := http.Get(url)
			if err != nil {
				return nil, fmt.Errorf("error fetching hostnames from Trickest: %v", err)
			}
			defer resp.Body.Close()

			// Read the hostnames from the response
			body, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}

			// Split the response into individual hostnames
			lines := strings.Split(string(body), "\n")
			for _, line := range lines {
				if line = strings.TrimSpace(line); line != "" {
					hostnames = append(hostnames, line)
				}
			}
			break // No need to continue after finding the correct domain
		}
	}

	return hostnames, nil
}
