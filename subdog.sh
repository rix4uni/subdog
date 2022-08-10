#!/usr/bin/env bash


if [[ ! -x "$(command -v jq)" ]]; then
	printf "Installing jq\n\n"
	if [ -f /etc/os-release ]; then apt install jq -y;
	elif [ -f /etc/redhat-release ]; then yum install jq -y;
	elif [ -f /etc/arch-release ]; then sudo pacman -S jq -y;
	fi
fi


passive_scan(){
	curl -s "https://rapiddns.io/subdomain/$domain?full=1#result" | grep "td" | grep ".$domain" | awk '{print $1}' | sed 's/<\/\?[^>]\+>//g' | sed '/<a$/d' | sort -u | tee -a subdog_subdomain.txt
	curl -s "https://api.threatminer.org/v2/domain.php?q=$domain&rt=5" | jq -r '.results[]' 2>/dev/null | grep -o "\w.*$domain" | sort -u | tee -a subdog_subdomain.txt
	curl -s "https://riddler.io/search/exportcsv?q=pld:$domain" | grep -o "\w.*$domain"| awk -F, '{print $6}' | sort -u | tee -a subdog_subdomain.txt
	curl -s "https://otx.alienvault.com/api/v1/indicators/domain/$domain/passive_dns" | jq '.passive_dns[].hostname' 2>/dev/null | grep -o "\w.*$domain"| sort -u | tee -a subdog_subdomain.txt
	curl -s "https://www.threatcrowd.org/searchApi/v2/domain/report/?domain=$domain" | jq -r '.subdomains' 2>/dev/null | grep -o "\w.*$domain" | sort -u | tee -a subdog_subdomain.txt


	echo "[+] Number of domains found: $(cat subdog_subdomain.txt | wc -l)"
}


PROGARGS=$(getopt -o ':d:h::' --long 'domain:,help' -n 'subdog' -- "$@")


# Note the quotes around "$PROGARGS": they are essential!
eval set -- "$PROGARGS"
unset PROGARGS


while true; do
    case "$1" in
        '-d'|'--domain')
            domain=$2
            shift 2
            continue
            ;;

        '-p'|'--passive'
        	passive_scan
        	shift 2
        	continue
        	;;

        '-o'|'--output'
        	touch $3
        	shift 2
        	continue
        	;;

        '-h'|'--help')
            showhelp
            exit 1
            ;;
    esac
done
