package subdomainfinder

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

// FetchSubdomains fetches subdomains for a given domain from the subdomainfinder API
func FetchSubdomains(domain string) ([]string, error) {
	url := "https://subdomainfinder.c99.nl/"

	// Prepare the POST request
	payload := fmt.Sprintf("CSRF9843411078797932=a&jn=JS+aan%%2C+T+aangeroepen%%2C+CSRF+aangepast&domain=%s&lol-stop-reverse-engineering-my-source-and-buy-an-api-key=b7daebdf32b85e93bbbf95e54b00e9d474ee0579&scan_subdomains=", domain)
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

	// Collect unique subdomains
	subdomainMap := make(map[string]struct{})
	for _, match := range subdomainMatches {
		if len(match) > 1 {
			subdomainMap[match[1]] = struct{}{}
		}
	}

	uniqueSubdomains := make([]string, 0, len(subdomainMap))
	for sub := range subdomainMap {
		uniqueSubdomains = append(uniqueSubdomains, sub)
	}

	// If no subdomains were found, return an error
	if len(uniqueSubdomains) == 0 {
		return nil, fmt.Errorf("no subdomains found for %s", domain)
	}

	return uniqueSubdomains, nil
}
