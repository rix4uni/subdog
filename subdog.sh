#!/usr/bin/env bash


# COLORS
BRED='\033[1;31m'
BBLUE='\033[1;34m'
BGREEN='\033[1;32m'
BYELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
RESET='\033[0m'

banner(){
    echo -e "${GREEN}
\t\t ___  __  __  ____  ____  _____  ___ 
\t\t/ __)(  )(  )(  _ \(  _ \(  _  )/ __)
\t\t\__ \ )(__)(  ) _ < )(_) ))(_)(( (_-.
\t\t(___/(______)(____/(____/(_____)\___/
\t\t                    coded by ${YELLOW}@rix4uni${RED} in INDIA${RESET}"
}


if [[ ! -x "$(command -v jq)" ]]; then
	printf "${bblue}Installing jq\n\n"
	if [ -f /etc/os-release ]; then apt install jq -y &>/dev/null;
	elif [ -f /etc/redhat-release ]; then yum install jq -y &>/dev/null;
	elif [ -f /etc/arch-release ]; then sudo pacman -S jq -y &>/dev/null;
	fi
fi


passive_scan(){
	curl -s "https://rapiddns.io/subdomain/$domain?full=1#result" | grep "td" | grep ".$domain" | awk '{print $1}' | sed 's/<\/\?[^>]\+>//g' | sed '/<a$/d' | sort -u >> ./subdog_subdomain_old.txt
	curl -s "https://api.threatminer.org/v2/domain.php?q=$domain&rt=5" | jq -r '.results[]' 2>/dev/null | grep -o "\w.*$domain" | sort -u >> ./subdog_subdomain_old.txt
	curl -s "https://riddler.io/search/exportcsv?q=pld:$domain" | grep -o "\w.*$domain"| awk -F, '{print $6}' | sort -u >> ./subdog_subdomain_old.txt
	curl -s "https://otx.alienvault.com/api/v1/indicators/domain/$domain/passive_dns" | jq '.passive_dns[].hostname' 2>/dev/null | grep -o "\w.*$domain"| sort -u >> ./subdog_subdomain_old.txt
	curl -s "https://www.threatcrowd.org/searchApi/v2/domain/report/?domain=$domain" | jq -r '.subdomains' 2>/dev/null | grep -o "\w.*$domain" | sort -u >> ./subdog_subdomain_old.txt

	# sorting subdomains
	cat ./subdog_subdomain_old.txt | sort -u > ./subdog_subdomain.txt | rm -rf ./subdog_subdomain_old.txt
	echo -e "${GREEN}$(cat ./subdog_subdomain.txt)${RESET}"
	echo -e "${YELLOW}[+] Number of domains found: $(cat ./subdog_subdomain.txt | wc -l | awk '{print $1}'${RESET})"
	echo -e "${YELLOW}[+] Saved output: $PWD/subdog_subdomain.txt${RESET}"
}


PROGARGS=$(getopt -o ':d:ph::' --long 'domain:passive,help' -n 'subdog' -- "$@")


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

        '-p'|'--passive')
        	passive_scan
        	shift 2
        	exit 1
        	;;

        '-h'|'--help')
            showhelp
            exit 1
            ;;
    esac
done
