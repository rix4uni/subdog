from bs4 import BeautifulSoup
import re
import requests
import sys

def get_subdomains():
    domains = sys.stdin.read().splitlines()
    for domain in domains:
        site_url = f"https://subdomainfinder.c99.nl/search.php?domain={domain}"
        response = requests.get(site_url).text
        soup = BeautifulSoup(response, "html.parser")
        a = soup.findAll("a", class_="text-decoration-none")
        for a in soup.find_all('a', href=True):
            h = a['href']
            if re.search("scans", h):
                url = "https:" + h
                response = requests.get(url).text
                s = BeautifulSoup(response, "html.parser")
                tab = s.find('table', {'id': 'result_table'})
                links = s.findAll("a")
                for link in links:
                    check = link.getText()
                    if check.endswith(f"{domain}"):
                        print(check)
get_subdomains()