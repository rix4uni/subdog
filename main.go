package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/rix4uni/subdog/banner"
	cmd "github.com/rix4uni/subdog/cmd"
	"github.com/spf13/pflag"
)

var availableSources = []string{
	"subdomaincenter",
	"jldc",
	"virustotal",
	"alienvault",
	"urlscan",
	"certspotter",
	"hackertarget",
	"crtsh",
	"trickest",
	"subdomainfinder",
	"chaos",
	"merklemap",
	"shodan",
	"reverseipdomain",
	"dnsdumpster",
	"bugbountydata",
}

// SourceStat tracks statistics for a single source execution
type SourceStat struct {
	Name     string
	Count    int
	Duration time.Duration
	Error    error
	ErrorMsg string
}

// formatNumber formats a number with comma separators
func formatNumber(n int) string {
	s := fmt.Sprintf("%d", n)
	if n < 1000 {
		return s
	}
	var result strings.Builder
	for i, r := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result.WriteString(",")
		}
		result.WriteRune(r)
	}
	return result.String()
}

// formatDuration formats duration to string with appropriate precision
func formatDuration(d time.Duration) string {
	seconds := d.Seconds()
	if seconds < 1 {
		return fmt.Sprintf("%.3fs", seconds)
	}
	return fmt.Sprintf("%.3fs", seconds)
}

// printSummary prints a formatted summary table with box-drawing characters
func printSummary(stats []SourceStat, totalDuration time.Duration, domain string) {
	fmt.Println()
	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	// Format title to fit exactly 60 characters (including borders)
	// Title text: "Subdomain Enumeration Summary - " (35 chars) + domain (max 23 chars) = 58 chars total
	titlePrefix := "Subdomain Enumeration Summary - "
	maxDomainLen := 58 - len(titlePrefix)
	domainDisplay := domain
	if len(domainDisplay) > maxDomainLen {
		domainDisplay = domainDisplay[:maxDomainLen-3] + "..."
	}
	title := fmt.Sprintf("%s%s", titlePrefix, domainDisplay)
	// Pad to exactly 58 characters (60 - 2 borders)
	if len(title) < 58 {
		title = title + strings.Repeat(" ", 58-len(title))
	}
	fmt.Printf("║%s║\n", title)
	fmt.Println("╠═══════════════════╦═════════╦═══════════╦════════════════╣")
	fmt.Println("║ Source            ║ Count   ║ Time      ║ Status         ║")
	fmt.Println("╠═══════════════════╬═════════╬═══════════╬════════════════╣")

	totalCount := 0
	for _, stat := range stats {
		totalCount += stat.Count
		durationStr := formatDuration(stat.Duration)

		// Truncate source name if too long (max 17 chars)
		sourceName := stat.Name
		if len(sourceName) > 17 {
			sourceName = sourceName[:14] + "..."
		}

		// Format count (max 7 chars)
		countStr := formatNumber(stat.Count)
		if len(countStr) > 7 {
			countStr = countStr[:4] + "..."
		}

		// Format duration (max 9 chars)
		if len(durationStr) > 9 {
			durationStr = durationStr[:6] + "..."
		}

		// Format status (max 14 chars)
		status := "✓ Success"
		if stat.Error != nil {
			errorMsg := stat.ErrorMsg
			// Account for "✗ " prefix (2 chars), so max 12 chars for error message
			if len(errorMsg) > 12 {
				errorMsg = errorMsg[:9] + "..."
			}
			status = fmt.Sprintf("✗ %s", errorMsg)
		}
		if len(status) > 14 {
			status = status[:11] + "..."
		}

		fmt.Printf("║ %-17s ║ %-7s ║ %-9s ║ %-14s ║\n",
			sourceName, countStr, durationStr, status)
	}

	fmt.Println("╠═══════════════════╬═════════╬═══════════╬════════════════╣")

	// Format total count and duration to fit column widths
	totalCountStr := formatNumber(totalCount)
	if len(totalCountStr) > 7 {
		totalCountStr = totalCountStr[:4] + "..."
	}
	totalDurationStr := formatDuration(totalDuration)
	if len(totalDurationStr) > 9 {
		totalDurationStr = totalDurationStr[:6] + "..."
	}

	fmt.Printf("║ TOTAL             ║ %-7s ║ %-9s ║                ║\n",
		totalCountStr, totalDurationStr)
	fmt.Println("╚═══════════════════╩═════════╩═══════════╩════════════════╝")
	fmt.Println()
}

// uniqueFileWriter writes unique lines to stdout and file
type uniqueFileWriter struct {
	file           *os.File
	seenSubdomains map[string]bool
	writtenToFile  map[string]bool
	mutex          *sync.Mutex
}

// Ensure uniqueFileWriter implements io.Writer
var _ io.Writer = (*uniqueFileWriter)(nil)

// Write implements io.Writer - writes unique lines to stdout and file
func (w *uniqueFileWriter) Write(p []byte) (n int, err error) {
	text := strings.TrimSpace(string(p))
	if text == "" {
		return len(p), nil
	}

	// Normalize the text (filter wildcards, emails, lowercase)
	normalized := cmd.NormalizeSubdomain(text)
	if normalized == "" {
		return len(p), nil // Filter out emails and wildcards
	}

	w.mutex.Lock()
	defer w.mutex.Unlock()

	// Write to stdout if we haven't seen this subdomain before
	if !w.seenSubdomains[normalized] {
		_, err = fmt.Println(normalized)
		if err != nil {
			return len(p), err
		}
		w.seenSubdomains[normalized] = true
	}

	// Write to file if specified and not already written to file
	// Mark BEFORE writing to prevent race conditions (defensive approach)
	if w.file != nil && !w.writtenToFile[normalized] {
		w.writtenToFile[normalized] = true        // Mark first to prevent other goroutines from writing
		_, err = fmt.Fprintln(w.file, normalized) // Then write
		if err != nil {
			// If write fails, unmark to allow retry (though this is unlikely)
			w.writtenToFile[normalized] = false
			return len(p), err
		}
	}

	n = len(p)
	return n, err
}

func main() {
	sources := pflag.StringP("source", "s", "all", "Choose source(s) to use, or 'all' for all sources. Use --list-sources to see available sources")
	excludeSources := pflag.StringP("exclude-source", "e", "", "Comma-separated list of sources to exclude when using --source all")
	listSources := pflag.BoolP("list-sources", "l", false, "List all available sources and exit")
	parallel := pflag.BoolP("parallel", "p", false, "Run all sources in parallel to speed up scanning")
	silent := pflag.Bool("silent", false, "Silent mode.")
	versionFlag := pflag.Bool("version", false, "Print the version of the tool and exit.")
	verbose := pflag.Bool("verbose", false, "enable verbose mode")
	outputFile := pflag.StringP("output", "o", "", "Save subdomain results to a file")
	pflag.Parse()

	if *listSources {
		fmt.Println("Available sources:")
		for _, source := range availableSources {
			fmt.Printf("  - %s\n", source)
		}
		fmt.Println("\nUse 'all' to run all sources")
		return
	}

	// Parse excluded sources into a map for quick lookup
	excludedMap := make(map[string]bool)
	if *excludeSources != "" {
		excluded := strings.Split(*excludeSources, ",")
		for _, e := range excluded {
			excludedMap[strings.TrimSpace(e)] = true
		}
	}

	// Helper function to check if a source should be executed
	shouldRun := func(sourceName string) bool {
		// If source is specifically requested, check if it's excluded
		if *sources == sourceName {
			return !excludedMap[sourceName]
		}
		// If "all" is requested, check if this source is excluded
		if *sources == "all" {
			return !excludedMap[sourceName]
		}
		return false
	}

	if *versionFlag {
		banner.PrintBanner()
		banner.PrintVersion()
		return
	}

	if !*silent {
		banner.PrintBanner()
	}

	// Open output file if specified
	var file *os.File
	var fileErr error
	if *outputFile != "" {
		file, fileErr = os.Create(*outputFile)
		if fileErr != nil {
			fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", fileErr)
			return
		}
		defer func() {
			file.Sync() // Ensure all writes are flushed
			file.Close()
		}()
	}

	// Output mutex for thread-safe printing when using --parallel
	var outputMutex sync.Mutex

	// Track unique subdomains for stdout (prevents duplicate stdout output)
	seenSubdomains := make(map[string]bool)
	// Track what's been written to file (prevents duplicate file writes)
	writtenToFile := make(map[string]bool)

	// Create unique writer for chaos source (writes unique lines to stdout and file)
	uniqueWriter := &uniqueFileWriter{
		file:           file,
		seenSubdomains: seenSubdomains,
		writtenToFile:  writtenToFile,
		mutex:          &outputMutex,
	}

	// Helper function to write output to both file and stdout
	// Both stdout and file only get unique lines (atomic operation)
	writeOutput := func(text string) {
		// Normalize the text first (filter wildcards, emails, lowercase)
		normalized := cmd.NormalizeSubdomain(text)
		if normalized == "" {
			return // Skip if filtered out
		}

		outputMutex.Lock()
		defer outputMutex.Unlock()

		// Write to stdout if we haven't seen this subdomain before
		if !seenSubdomains[normalized] {
			fmt.Println(normalized)
			seenSubdomains[normalized] = true
		}

		// Write to file if specified and not already written to file
		// Mark BEFORE writing to prevent race conditions (defensive approach)
		if file != nil && !writtenToFile[normalized] {
			writtenToFile[normalized] = true // Mark first to prevent other goroutines from writing
			fmt.Fprintln(file, normalized)   // Then write
		}
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		domain := strings.TrimSpace(scanner.Text())

		// Statistics tracking for this domain
		var stats []SourceStat
		var statsMutex sync.Mutex
		domainStartTime := time.Now()

		// Function to process a single source
		processSource := func(sourceName string, domain string, fetchFunc func(string) ([]string, error)) {
			if !shouldRun(sourceName) {
				return
			}
			startTime := time.Now()
			if *verbose {
				outputMutex.Lock()
				fmt.Printf("Fetching from %s for %s\n", sourceName, domain)
				outputMutex.Unlock()
			}
			results, err := fetchFunc(domain)
			duration := time.Since(startTime)

			if err != nil {
				if *verbose {
					outputMutex.Lock()
					fmt.Printf("Error fetching subdomains from %s for %s: %v\n", sourceName, domain, err)
					outputMutex.Unlock()
				}
				statsMutex.Lock()
				stats = append(stats, SourceStat{
					Name:     sourceName,
					Count:    0,
					Duration: duration,
					Error:    err,
					ErrorMsg: err.Error(),
				})
				statsMutex.Unlock()
				return
			}

			// Count unique subdomains written by this source (must be done inside mutex)
			outputMutex.Lock()
			countBefore := len(seenSubdomains)
			outputMutex.Unlock()

			for _, result := range results {
				writeOutput(result)
			}

			outputMutex.Lock()
			countAfter := len(seenSubdomains)
			count := countAfter - countBefore
			outputMutex.Unlock()

			statsMutex.Lock()
			stats = append(stats, SourceStat{
				Name:     sourceName,
				Count:    count,
				Duration: duration,
				Error:    nil,
			})
			statsMutex.Unlock()
		}

		// Function to process chaos source (special case - no return values)
		processChaosSource := func(domain string) {
			if !shouldRun("chaos") {
				return
			}
			startTime := time.Now()
			if *verbose {
				outputMutex.Lock()
				fmt.Printf("Fetching from Chaos for %s\n", domain)
				outputMutex.Unlock()
			}

			// Track count before chaos execution (lock, read, unlock)
			outputMutex.Lock()
			countBefore := len(seenSubdomains)
			outputMutex.Unlock()

			// Call ProcessDomainChaos WITHOUT holding the mutex
			// uniqueFileWriter.Write will handle its own mutex locking
			cmd.ProcessDomainChaos(domain, uniqueWriter)

			// Count new unique subdomains added by chaos (lock, read, unlock)
			outputMutex.Lock()
			countAfter := len(seenSubdomains)
			chaosCount := countAfter - countBefore
			outputMutex.Unlock()

			duration := time.Since(startTime)
			statsMutex.Lock()
			stats = append(stats, SourceStat{
				Name:     "chaos",
				Count:    chaosCount,
				Duration: duration,
				Error:    nil,
			})
			statsMutex.Unlock()
		}

		if *parallel {
			// Run all sources in parallel
			var wg sync.WaitGroup

			wg.Add(1)
			go func() {
				defer wg.Done()
				processSource("subdomaincenter", domain, cmd.FetchSubdomainsSubdomaincenter)
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				processSource("jldc", domain, cmd.FetchSubdomainsJldc)
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				processSource("virustotal", domain, cmd.FetchSubdomainsVirusTotal)
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				processSource("alienvault", domain, cmd.FetchSubdomainsAlienVault)
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				processSource("urlscan", domain, cmd.FetchSubdomainsURLScan)
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				processSource("certspotter", domain, func(d string) ([]string, error) {
					return cmd.FetchDNSNamesCertspotter(d)
				})
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				processSource("hackertarget", domain, cmd.FetchSubdomainsHackerTarget)
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				processSource("crtsh", domain, cmd.FetchSubdomainsCrtsh)
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				processSource("trickest", domain, func(d string) ([]string, error) {
					return cmd.FetchHostnamesTrickest(d)
				})
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				processSource("subdomainfinder", domain, cmd.FetchSubdomainsSubdomainFinder)
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				processChaosSource(domain)
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				processSource("merklemap", domain, cmd.FetchDomainsMerkleMap)
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				processSource("shodan", domain, cmd.FetchSubdomainsShodan)
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				processSource("reverseipdomain", domain, cmd.FetchSubdomainsReverseIPDomain)
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				processSource("dnsdumpster", domain, cmd.FetchSubdomainsDNSDumpster)
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				processSource("bugbountydata", domain, cmd.FetchSubdomainsBugBountyData)
			}()

			wg.Wait()

			// Print summary for parallel execution
			if *verbose {
				totalDuration := time.Since(domainStartTime)
				printSummary(stats, totalDuration, domain)
			}
		} else {
			// Sequential execution
			// Helper function for sequential source processing with stats tracking
			processSequentialSource := func(sourceName string, fetchFunc func(string) ([]string, error)) {
				if !shouldRun(sourceName) {
					return
				}
				startTime := time.Now()
				if *verbose {
					fmt.Printf("Fetching from %s for %s\n", sourceName, domain)
				}
				subdomains, err := fetchFunc(domain)
				duration := time.Since(startTime)
				if err != nil {
					if *verbose {
						fmt.Printf("Error fetching subdomains from %s for %s: %v\n", sourceName, domain, err)
					}
					stats = append(stats, SourceStat{
						Name:     sourceName,
						Count:    0,
						Duration: duration,
						Error:    err,
						ErrorMsg: err.Error(),
					})
					return
				}
				// Count unique subdomains written by this source (must be done inside mutex)
				outputMutex.Lock()
				countBefore := len(seenSubdomains)
				outputMutex.Unlock()

				for _, subdomain := range subdomains {
					writeOutput(subdomain)
				}

				outputMutex.Lock()
				countAfter := len(seenSubdomains)
				count := countAfter - countBefore
				outputMutex.Unlock()
				stats = append(stats, SourceStat{
					Name:     sourceName,
					Count:    count,
					Duration: duration,
					Error:    nil,
				})
			}

			processSequentialSource("subdomaincenter", cmd.FetchSubdomainsSubdomaincenter)
			processSequentialSource("jldc", cmd.FetchSubdomainsJldc)
			processSequentialSource("virustotal", cmd.FetchSubdomainsVirusTotal)
			processSequentialSource("alienvault", cmd.FetchSubdomainsAlienVault)
			processSequentialSource("urlscan", cmd.FetchSubdomainsURLScan)
			processSequentialSource("certspotter", func(d string) ([]string, error) {
				return cmd.FetchDNSNamesCertspotter(d)
			})
			processSequentialSource("hackertarget", cmd.FetchSubdomainsHackerTarget)
			processSequentialSource("crtsh", cmd.FetchSubdomainsCrtsh)
			processSequentialSource("trickest", func(d string) ([]string, error) {
				return cmd.FetchHostnamesTrickest(d)
			})
			processSequentialSource("subdomainfinder", cmd.FetchSubdomainsSubdomainFinder)

			// Handle chaos source in sequential mode
			if shouldRun("chaos") {
				startTime := time.Now()
				if *verbose {
					fmt.Printf("Fetching from Chaos for %s\n", domain)
				}
				// Track count before chaos execution (lock, read, unlock)
				outputMutex.Lock()
				countBefore := len(seenSubdomains)
				outputMutex.Unlock()

				// Call ProcessDomainChaos WITHOUT holding the mutex
				// uniqueFileWriter.Write will handle its own mutex locking
				cmd.ProcessDomainChaos(domain, uniqueWriter)

				// Count new unique subdomains added by chaos (lock, read, unlock)
				outputMutex.Lock()
				countAfter := len(seenSubdomains)
				chaosCount := countAfter - countBefore
				outputMutex.Unlock()

				duration := time.Since(startTime)
				stats = append(stats, SourceStat{
					Name:     "chaos",
					Count:    chaosCount,
					Duration: duration,
					Error:    nil,
				})
			}

			processSequentialSource("merklemap", cmd.FetchDomainsMerkleMap)
			processSequentialSource("shodan", cmd.FetchSubdomainsShodan)
			processSequentialSource("reverseipdomain", cmd.FetchSubdomainsReverseIPDomain)
			processSequentialSource("dnsdumpster", cmd.FetchSubdomainsDNSDumpster)
			processSequentialSource("bugbountydata", cmd.FetchSubdomainsBugBountyData)

			// Print summary for sequential execution
			if *verbose {
				totalDuration := time.Since(domainStartTime)
				printSummary(stats, totalDuration, domain)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading input: %v\n", err)
	}
}
