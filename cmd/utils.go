package cmd

import "strings"

// isSubdomainOrDomain checks if a domain is either the target domain itself or a subdomain of it
// It matches the pattern ^(.*\.)?target\.domain$ (e.g., "dell.com" or "admin.dell.com")
func isSubdomainOrDomain(domain, targetDomain string) bool {
	if domain == "" || targetDomain == "" {
		return false
	}
	// Exact match
	if domain == targetDomain {
		return true
	}
	// Subdomain match: domain ends with "." + targetDomain
	return strings.HasSuffix(domain, "."+targetDomain)
}

// NormalizeSubdomain filters out email addresses (strings containing "@") and wildcard subdomains (starting with "*.")
// Returns empty string if the input contains "@" or starts with "*.", otherwise returns the lowercase version
func NormalizeSubdomain(subdomain string) string {
	if strings.Contains(subdomain, "@") {
		return ""
	}
	if strings.HasPrefix(subdomain, "*.") {
		return ""
	}
	return strings.ToLower(subdomain)
}
