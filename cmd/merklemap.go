package cmd

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"
)

// FetchDomainsMerkleMap fetches subdomains for a given domain from the MerkleMap API
func FetchDomainsMerkleMap(domain string) ([]string, error) {
	url := fmt.Sprintf("https://api.merklemap.com/search?query=%s&stream=true", domain)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	// Use a map to track unique filtered subdomains
	subdomainMap := make(map[string]bool)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "data: ") {
			line = strings.TrimPrefix(line, "data: ")
			// Extract domain field from the JSON response
			if strings.Contains(line, `"domain"`) {
				domainStart := strings.Index(line, `"domain":"`) + 10
				domainEnd := strings.Index(line[domainStart:], `"`) + domainStart
				domainName := line[domainStart:domainEnd]
				if isSubdomainOrDomain(domainName, domain) {
					normalized := NormalizeSubdomain(domainName)
					if normalized != "" {
						subdomainMap[normalized] = true
					}
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Convert map to slice
	subdomains := make([]string, 0, len(subdomainMap))
	for subdomain := range subdomainMap {
		subdomains = append(subdomains, subdomain)
	}

	return subdomains, nil
}
