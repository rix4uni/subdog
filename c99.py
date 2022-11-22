import requests
import sys
from bs4 import BeautifulSoup

file1 = open('footlocker.txt', 'r')
Lines = file1.readlines()

site_url = "https://subdomainfinder.c99.nl/"

# Strips the newline character
for line in Lines:
    complete = site_url + line.strip()
    url = requests.get(complete).text
    s = BeautifulSoup(url,"lxml")
    tab = s.find('table',{'id':'result_table'})
    l = s.findAll("a")
    for l1 in l:
        check = l1.getText()
        if check.endswith(sys.argv[1]):
            print(check)
