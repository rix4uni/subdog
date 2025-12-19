package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

// FetchSubdomainsDNSDumpster fetches subdomains from DNSDumpster for a given domain
func FetchSubdomainsDNSDumpster(domain string) ([]string, error) {
	url := "https://dnsdumpster.com/"
	csrfToken := "vZGRO1YfdzdviMYTYZqLrw0PxsV5mlAnVGFadIqkjIAhiyNgi5w70hIj7uuzdmXx" // Hardcoded CSRF token

	// Create a new request
	client := &http.Client{}
	form := fmt.Sprintf("csrfmiddlewaretoken=%s&targetip=%s&user=free", csrfToken, domain)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(form)))
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", url)
	req.Header.Set("Cookie", "csrftoken=" + csrfToken)

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Regex to match the subdomain entries
	re := regexp.MustCompile(`<tr><td class="col-md-4">(.*?)<br</td>`)
	matches := re.FindAllStringSubmatch(string(body), -1)

	// Use a map to track unique filtered subdomains
	subdomainMap := make(map[string]bool)
	for _, match := range matches {
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

	// Convert map to slice
	subdomains := make([]string, 0, len(subdomainMap))
	for subdomain := range subdomainMap {
		subdomains = append(subdomains, subdomain)
	}

	return subdomains, nil
}
