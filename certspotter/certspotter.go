package certspotter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type CertspotterResponse []struct {
	DNSNames []string `json:"dns_names"`
}

// FetchDNSNames fetches DNS names for a given domain from the Certspotter API
func FetchDNSNames(domain string) ([]string, error) {
	url := fmt.Sprintf("https://api.certspotter.com/v1/issuances?domain=%s&include_subdomains=true&expand=dns_names", domain)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var csResponse CertspotterResponse
	err = json.Unmarshal(body, &csResponse)
	if err != nil {
		return nil, err
	}

	var dnsNames []string
	for _, issuance := range csResponse {
		dnsNames = append(dnsNames, issuance.DNSNames...)
	}

	return dnsNames, nil
}
