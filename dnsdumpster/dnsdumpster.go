package dnsdumpster

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

// FetchSubdomains fetches subdomains from DNSDumpster for a given domain
func FetchSubdomains(domain string) ([]string, error) {
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

	var subdomains []string
	for _, match := range matches {
		if len(match) > 1 {
			subdomains = append(subdomains, match[1])
		}
	}

	return subdomains, nil
}
