package jldc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// FetchSubdomains fetches subdomains for a given domain from the jldc.me API
func FetchSubdomains(domain string) ([]string, error) {
	url := fmt.Sprintf("https://jldc.me/anubis/subdomains/%s", domain)

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
