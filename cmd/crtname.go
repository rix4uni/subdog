package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// FetchSubdomainsCrtName fetches subdomains for a given domain from crt.name (plain-text, one host per line)
func FetchSubdomainsCrtName(domain string) ([]string, error) {
	url := fmt.Sprintf("https://crt.name/v1/search?apex=%s", domain)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("crt.name returned status %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	subdomainMap := make(map[string]bool)
	for _, line := range strings.Split(string(body), "\n") {
		line = strings.TrimSpace(line)
		if line != "" && isSubdomainOrDomain(line, domain) {
			if normalized := NormalizeSubdomain(line); normalized != "" {
				subdomainMap[normalized] = true
			}
		}
	}

	results := make([]string, 0, len(subdomainMap))
	for subdomain := range subdomainMap {
		results = append(results, subdomain)
	}

	return results, nil
}
