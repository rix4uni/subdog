package merklemap

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"
)

// FetchDomains fetches subdomains for a given domain from the MerkleMap API
func FetchDomains(domain string) ([]string, error) {
	url := fmt.Sprintf("https://api.merklemap.com/search?query=%s&stream=true", domain)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	var subdomains []string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "data: ") {
			line = strings.TrimPrefix(line, "data: ")
			// Extract domain field from the JSON response
			if strings.Contains(line, `"domain"`) {
				domainStart := strings.Index(line, `"domain":"`) + 10
				domainEnd := strings.Index(line[domainStart:], `"`) + domainStart
				subdomains = append(subdomains, line[domainStart:domainEnd])
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return subdomains, nil
}
