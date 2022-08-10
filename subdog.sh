curl -s "https://rapiddns.io/subdomain/$domain?full=1#result" | grep "td" | grep ".$domain" | awk '{print $1}' | sed 's/<\/\?[^>]\+>//g' | sed '/<a$/d' | sort -u
curl -s "https://api.threatminer.org/v2/domain.php?q=$domain&rt=5" | jq -r '.results[]' 2>/dev/null | grep -o "\w.*$domain" | sort -u
curl -s "https://riddler.io/search/exportcsv?q=pld:$domain" | grep -o "\w.*$domain"| awk -F, '{print $6}' | sort -u
curl -s "https://otx.alienvault.com/api/v1/indicators/domain/$domain/passive_dns" | jq '.passive_dns[].hostname' 2>/dev/null | grep -o "\w.*$domain"| sort -u
curl -s "https://www.threatcrowd.org/searchApi/v2/domain/report/?domain=$domain" | jq -r '.subdomains' 2>/dev/null | grep -o "\w.*$domain" | sort -u
