package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// FetchSubdomainsSubdomaincenter fetches subdomains for a given domain from the API
func FetchSubdomainsSubdomaincenter(domain string) ([]string, error) {
	url := fmt.Sprintf("https://api.subdomain.center/?domain=%s", domain)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var subdomains []string
	err = json.Unmarshal(body, &subdomains)
	if err != nil {
		return nil, err
	}

	return subdomains, nil
}
