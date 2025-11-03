package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// FetchSubdomainsReverseIPDomain fetches subdomains for a given domain from the Reverse IP Domain API
func FetchSubdomainsReverseIPDomain(domain string) ([]string, error) {
	url := fmt.Sprintf("https://sub-scan-api.reverseipdomain.com/?domain=%s", domain)

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
		Result struct {
			Domains []string `json:"domains"`
		} `json:"result"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response.Result.Domains, nil
}
