package cmd

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Program struct {
	Name    string   `json:"name"`
	Domains []string `json:"domains"`
}

type Root struct {
	Programs []Program `json:"programs"`
}

type ChaosData struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// ProcessDomainChaos fetches and processes domain data from chaos
// If writer is nil, it writes to os.Stdout
func ProcessDomainChaos(domain string, writer io.Writer) error {
	if writer == nil {
		writer = os.Stdout
	}
	// Fetch bugbounty list
	resp, err := http.Get("https://raw.githubusercontent.com/projectdiscovery/public-bugbounty-programs/main/chaos-bugbounty-list.json")
	if err != nil {
		return fmt.Errorf("error fetching bugbounty list: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bugbounty list API returned status %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading bugbounty response: %v", err)
	}

	var root Root
	err = json.Unmarshal(body, &root)
	if err != nil {
		return fmt.Errorf("error unmarshaling bugbounty json: %v", err)
	}

	var programName string
	for _, p := range root.Programs {
		for _, d := range p.Domains {
			if d == domain {
				programName = p.Name
				break
			}
		}
		if programName != "" {
			break
		}
	}

	// Fetch chaos data index
	resp, err = http.Get("https://chaos-data.projectdiscovery.io/index.json")
	if err != nil {
		return fmt.Errorf("error fetching chaos data index: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("chaos data index API returned status %d", resp.StatusCode)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading chaos index response: %v", err)
	}

	var chaosData []ChaosData
	err = json.Unmarshal(body, &chaosData)
	if err != nil {
		return fmt.Errorf("error unmarshaling chaos index json: %v", err)
	}

	for _, data := range chaosData {
		if data.Name == programName {
			// Download and extract
			resp, err := http.Get(data.URL)
			if err != nil {
				return fmt.Errorf("error downloading data: %v", err)
			}
			defer resp.Body.Close()

			dataBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("error reading data response: %v", err)
			}

			zipReader, err := zip.NewReader(bytes.NewReader(dataBody), int64(len(dataBody)))
			if err != nil {
				return fmt.Errorf("error reading zip: %v", err)
			}

			for _, file := range zipReader.File {
				if file.Name == domain+".txt" {
					rc, err := file.Open()
					if err != nil {
						return fmt.Errorf("error opening zip file: %v", err)
					}

					// Read and process lines to filter emails and normalize to lowercase
					scanner := bufio.NewScanner(rc)
					for scanner.Scan() {
						line := strings.TrimSpace(scanner.Text())
						if line != "" {
							normalized := NormalizeSubdomain(line)
							if normalized != "" {
								fmt.Fprintln(writer, normalized)
							}
						}
					}
					rc.Close()
				}
			}
		}
	}
	return nil
}
