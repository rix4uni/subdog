from bs4 import BeautifulSoup
import re
import requests
import sys

def get_subdomains():
    domains = sys.stdin.read().splitlines()
    for domain in domains:
        site_url = f"https://subdomainfinder.c99.nl/search.php?domain={domain}"
        response = requests.get(site_url)
        soup = BeautifulSoup(response.text, "html.parser")

        for a in soup.find_all('a', href=True):
            if re.search("scans", a['href']):
                url = "https:" + a['href']
                response = requests.get(url)
                soup = BeautifulSoup(response.text, "html.parser")
                tabs = soup.find_all(class_="link")
                for tab in tabs:
                    if re.search(f"{domain}", tab.text):
                        print(tab.text)
get_subdomains()