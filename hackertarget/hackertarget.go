package hackertarget

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// FetchSubdomains fetches subdomains for a given domain from the HackerTarget API
func FetchSubdomains(domain string) ([]string, error) {
	url := fmt.Sprintf("https://api.hackertarget.com/hostsearch/?q=%s", domain)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse the CSV response
	reader := csv.NewReader(strings.NewReader(string(body)))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var subdomains []string
	for _, record := range records {
		if len(record) > 0 {
			subdomains = append(subdomains, record[0]) // First column contains the subdomain
		}
	}

	return subdomains, nil
}
