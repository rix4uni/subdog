from concurrent.futures import ThreadPoolExecutor
from bs4 import BeautifulSoup
import argparse
import re
import requests

def get_subdomains(domain):
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

if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("-d", "--domain", type=str, help="The domain to find subdomains for")
    parser.add_argument("-t", "--threads", type=int, default=8, help="Number of threads to use")
    args = parser.parse_args()

    domain = args.domain
    threads = args.threads

    try:
        with ThreadPoolExecutor(max_workers=threads) as executor:
            executor.submit(get_subdomains, domain)
    except KeyboardInterrupt:
        exit(0)