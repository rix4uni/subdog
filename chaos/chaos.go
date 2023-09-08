package main

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

func processDomain(domain string) {
	// Fetch bugbounty list
	resp, err := http.Get("https://raw.githubusercontent.com/projectdiscovery/public-bugbounty-programs/main/chaos-bugbounty-list.json")
	if err != nil {
		fmt.Println("Error fetching bugbounty list:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	var root Root
	err = json.Unmarshal(body, &root)
	if err != nil {
		fmt.Println("Error unmarshaling json:", err)
		return
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
		fmt.Println("Error fetching chaos data index:", err)
		return
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	var chaosData []ChaosData
	err = json.Unmarshal(body, &chaosData)
	if err != nil {
		fmt.Println("Error unmarshaling json:", err)
		return
	}

	for _, data := range chaosData {
		if data.Name == programName {
			// Download and extract
			resp, err := http.Get(data.URL)
			if err != nil {
				fmt.Println("Error downloading data:", err)
				return
			}
			defer resp.Body.Close()

			dataBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading data response:", err)
				return
			}

			zipReader, err := zip.NewReader(bytes.NewReader(dataBody), int64(len(dataBody)))
			if err != nil {
				fmt.Println("Error reading zip:", err)
				return
			}

			for _, file := range zipReader.File {
				if file.Name == domain+".txt" {
					rc, err := file.Open()
					if err != nil {
						fmt.Println("Error opening zip file:", err)
						return
					}

					io.Copy(os.Stdout, rc) // Print to stdout
					rc.Close()
				}
			}
		}
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		domain := scanner.Text()
		processDomain(domain)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
	}
}
