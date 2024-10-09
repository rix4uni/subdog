package virustotal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type VirusTotalResponse struct {
	Data []struct {
		ID string `json:"id"`
	} `json:"data"`
}

// FetchSubdomains fetches subdomains for a given domain from VirusTotal API
func FetchSubdomains(domain string) ([]string, error) {
	url := fmt.Sprintf("https://www.virustotal.com/ui/domains/%s/subdomains?limit=1000&relationships=resolutions", domain)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Tool", "vt-ui-main")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	req.Header.Set("Accept-Ianguage", "en-US,en;q=0.9,es;q=0.8")
	req.Header.Set("X-VT-Anti-Abuse-Header", "MTY1NjA5Nzk1NjAtWkc5dWRDQmlaU0JsZG1scy0xNjgzNDI2MDY4Ljc2MQ==")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var vtResponse VirusTotalResponse
	err = json.Unmarshal(body, &vtResponse)
	if err != nil {
		return nil, err
	}

	var subdomains []string
	for _, data := range vtResponse.Data {
		subdomains = append(subdomains, data.ID)
	}

	return subdomains, nil
}
