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

	// Use a map to track unique filtered subdomains
	subdomainMap := make(map[string]bool)
	
	// Helper function to process a name string (can contain multiple subdomains)
	processNameString := func(nameStr string) {
		if nameStr == "" {
			return
		}
		// Split by newlines first, then by commas, then by spaces
		// This handles all possible delimiters in the name_value field
		lines := strings.Split(nameStr, "\n")
		for _, line := range lines {
			// Split by commas
			commaParts := strings.Split(line, ",")
			for _, commaPart := range commaParts {
				// Split by spaces
				spaceParts := strings.Fields(commaPart)
				for _, name := range spaceParts {
					trimmedName := strings.TrimSpace(name)
					if trimmedName != "" && isSubdomainOrDomain(trimmedName, domain) {
						normalized := NormalizeSubdomain(trimmedName)
						if normalized != "" {
							subdomainMap[normalized] = true
						}
					}
				}
			}
		}
	}
	
	for _, cert := range certs {
		// Process CommonName field
		if cert.CommonName != "" {
			trimmedCN := strings.TrimSpace(cert.CommonName)
			if trimmedCN != "" && isSubdomainOrDomain(trimmedCN, domain) {
				normalized := NormalizeSubdomain(trimmedCN)
				if normalized != "" {
					subdomainMap[normalized] = true
				}
			}
		}
		
		// Process NameValue field (can contain multiple subdomains)
		processNameString(cert.NameValue)
	}

	// Convert map to slice
	results := make([]string, 0, len(subdomainMap))
	for subdomain := range subdomainMap {
		results = append(results, subdomain)
	}

	return results, nil
}
