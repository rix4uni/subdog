package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var (
	subfinderOnce sync.Once
	subfinderPath string
	subfinderErr  error
)

// ensureSubfinder locates the subfinder binary, installing it with
// `go install` if it is not found. It runs at most once per process.
func ensureSubfinder() (string, error) {
	subfinderOnce.Do(func() {
		if p, err := exec.LookPath("subfinder"); err == nil {
			subfinderPath = p
			return
		}

		fmt.Fprintln(os.Stderr, "subfinder not found, installing: go install -v github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest")
		install := exec.Command("go", "install", "-v", "github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest")
		install.Stdout = os.Stderr
		install.Stderr = os.Stderr
		if err := install.Run(); err != nil {
			subfinderErr = fmt.Errorf("subfinder is not installed and automatic installation failed: %v", err)
			return
		}

		if p, err := exec.LookPath("subfinder"); err == nil {
			subfinderPath = p
			return
		}

		// The Go bin directory may not be on PATH; look there directly.
		name := "subfinder"
		if runtime.GOOS == "windows" {
			name += ".exe"
		}
		var dirs []string
		if out, err := exec.Command("go", "env", "GOBIN").Output(); err == nil {
			dirs = append(dirs, strings.TrimSpace(string(out)))
		}
		if out, err := exec.Command("go", "env", "GOPATH").Output(); err == nil {
			for _, gp := range filepath.SplitList(strings.TrimSpace(string(out))) {
				dirs = append(dirs, filepath.Join(gp, "bin"))
			}
		}
		for _, d := range dirs {
			if d == "" {
				continue
			}
			candidate := filepath.Join(d, name)
			if _, err := os.Stat(candidate); err == nil {
				subfinderPath = candidate
				return
			}
		}
		subfinderErr = fmt.Errorf("subfinder was installed but could not be found; add your Go bin directory to PATH")
	})
	return subfinderPath, subfinderErr
}

// FetchSubdomainsSubfinder fetches subdomains using the external subfinder tool
func FetchSubdomainsSubfinder(domain string) ([]string, error) {
	// Run: echo "domain" | subfinder -duc -silent -all
	bin, err := ensureSubfinder()
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(bin, "-duc", "-silent", "-all", "-d", domain)

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
