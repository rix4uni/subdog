package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

// FetchSubdomainsSubdomainFinder fetches subdomains for a given domain from the subdomainfinder API
func FetchSubdomainsSubdomainFinder(domain string) ([]string, error) {
	url := "https://subdomainfinder.c99.nl/"

	// Prepare the POST request
	payload := fmt.Sprintf("CSRF9843433218797932=pirate107704869&is_admin=false&jn=JS+aan%2C+T+aangeroepen%2C+CSRF+aangepast&domain=%s&lol-stop-reverse-engineering-my-source-and-buy-an-api-key=cf917529992fd6f916e2b4ef8b37c6d97f040eba&scan_subdomains=", domain)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return nil, err
	}

	// Set the headers
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.6533.100 Safari/537.36")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %s", resp.Status)
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Extract the subdomains from the response body
	// Modified regex pattern to avoid lookbehind
	subdomainRegex := regexp.MustCompile(`href='//([^']+)'`)
	subdomainMatches := subdomainRegex.FindAllStringSubmatch(string(body), -1)

	// Use a map to track unique filtered subdomains
	subdomainMap := make(map[string]bool)
	for _, match := range subdomainMatches {
		if len(match) > 1 {
			subdomain := match[1]
			if isSubdomainOrDomain(subdomain, domain) {
				normalized := NormalizeSubdomain(subdomain)
				if normalized != "" {
					subdomainMap[normalized] = true
				}
			}
		}
	}

	uniqueSubdomains := make([]string, 0, len(subdomainMap))
	for sub := range subdomainMap {
		uniqueSubdomains = append(uniqueSubdomains, sub)
	}

	return uniqueSubdomains, nil
}
