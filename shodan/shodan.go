package shodan

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// FetchSubdomains fetches subdomains for a given domain from the Shodan
func FetchSubdomains(domain string) ([]string, error) {
	url := fmt.Sprintf("https://api.shodan.io/dns/domain/%s", domain)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response struct {
		Subdomains []string `json:"subdomains"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	// Append the full domain to each subdomain
	fullSubdomains := make([]string, len(response.Subdomains))
	for i, sub := range response.Subdomains {
		fullSubdomains[i] = fmt.Sprintf("%s.%s", sub, domain)
	}

	return fullSubdomains, nil
}
