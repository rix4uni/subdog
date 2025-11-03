package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// CsrData represents the JSON response structure from crt.sh
type CsrData struct {
	CommonName string `json:"common_name"`
	NameValue  string `json:"name_value"`
}

// FetchSubdomainsCrtsh fetches subdomains for a given domain from the crt.sh API
func FetchSubdomainsCrtsh(domain string) ([]string, error) {
	url := fmt.Sprintf("https://crt.sh/?q=%s&output=json", domain)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var certs []CsrData
	if err := json.Unmarshal(body, &certs); err != nil {
		return nil, err
	}

	var results []string
	for _, cert := range certs {
		// Split the NameValue field by commas and then by spaces
		names := strings.FieldsFunc(cert.NameValue, func(r rune) bool {
			return r == ',' || r == ' '
		})
		results = append(results, names...)
	}

	// Print results directly without duplicates
	for _, result := range results {
		trimmedResult := strings.TrimSpace(result)
		if trimmedResult != "" {
			fmt.Println(trimmedResult)
		}
	}

	return results, nil // Return results in case it's needed elsewhere
}
