package cmd

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
)

// FetchSubdomainsSubfinder fetches subdomains using the external subfinder tool
func FetchSubdomainsSubfinder(domain string) ([]string, error) {
	// Run: echo "domain" | subfinder -duc -silent -all
	cmd := exec.Command("subfinder", "-duc", "-silent", "-all", "-d", domain)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error running subfinder: %v", err)
	}

	// Parse the output line by line
	var subdomains []string
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			normalized := NormalizeSubdomain(line)
			if normalized != "" {
				subdomains = append(subdomains, normalized)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning subfinder output: %v", err)
	}

	return subdomains, nil
}
