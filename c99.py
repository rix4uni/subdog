import re
import requests
import sys

output_set = set()  # Set to store unique lines

for line in sys.stdin:
    domain = line.strip()

    headers = {
        'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36',
        'Content-Type': 'application/x-www-form-urlencoded'
    }

    data = {
        'CSRF1026678917967820': 'subnet_ip106570706',
        'CSRF1110069873195520': 'phisher107742278',
        'CSRF1069049887721126': 'malware101109442',
        'CSRF1006010564689816': 'cracker107427594',
        'CSRF1062459611835528': 'bot105178731',
        'CSRF1060746757887332': 'spammer100812419',
        'CSRF1022301678324918': 'network109677702',
        'CSRF1011654517265110': 'subnet_ip102785982',
        'CSRF1092569843373750': 'stalker102366075',
        'CSRF1036350231823663': 'stalker101432705',
        'CSRF1023618189141243': 'CSRF104142642',
        'CSRF1109746853977066': 'thief104170467',
        'CSRF1062460397164181': 'malware109068665',
        'CSRF1047996790303129': 'identitytheft105064706',
        'CSRF1045590297095199': 'cyberspace101176868',
        'CSRF1064158129702326': 'bot103978853',
        'CSRF1041828215900991': 'pirate106830137',
        'CSRF1108857192496320': 'cyber104694561',
        'CSRF1049128538789078': 'computer106801022',
        'CSRF1097304165602833': 'subnet_ip109993268',
        'CSRF1043759433233265': 'geek102207732',
        'CSRF1049274074868055': 'CSRF105833260',
        'CSRF1107037085276134': 'bogey103183448',
        'CSRF1008637876372338': 'security109782516',
        'CSRF1098551843787092': 'intruder104617209',
        'CSRF1105083555809864': 'vulnerability102553659',
        'CSRF1073652919491369': 'tenant111107018',
        'CSRF1020853807735606': 'subnet_ip106528200',
        'CSRF1027406919391413': 'espionage105253696',
        'CSRF1105429060838860': 'programmer107297692',
        'CSRF1062313749491890': 'intrusion101273125',
        'CSRF1063973752082886': 'bogey100008519',
        'CSRF1053452521549933': 'techie102123215',
        'CSRF1074997027815144': 'bogey103858462',
        'CSRF1101960335833182': 'CSRF103111343',
        'CSRF1046264093876229': 'hacking109867406',
        'CSRF1006532262417901': 'cyberwar106931993',
        'CSRF1063207470959835': 'bot105729253',
        'CSRF1099826046217150': 'tenant101838209',
        'CSRF1016062260861713': 'phisher108159458',
        'CSRF1095447115612290': 'CSRF100225669',
        'CSRF1070782222063318': 'pirate103343400',
        'CSRF1032372413574012': 'geek106471080',
        'CSRF1104901726419672': 'subnet101263964',
        'CSRF1008313945677950': 'breach104594539',
        'CSRF1097760797134098': 'hacking104254710',
        'CSRF1089815352997209': 'ipv4104455423',
        'CSRF1093049577609262': 'mask100569177',
        'CSRF1016983369924546': 'subnet103336076',
        'CSRF1018628175213948': 'attacker107143677',
        'CSRF1090199874315778': 'CSRF105562368',
        'CSRF1089293745478469': 'phisher101797530',
        'CSRF1066257680544260': 'pirate109244305',
        'CSRF1044463393610135': 'techie100664114',
        'CSRF1009515592398279': 'CSRF101147243',
        'CSRF1005873714074577': 'Trojan101975082',
        'CSRF1052967994876624': 'intruder110531115',
        'CSRF1025419042705406': 'subnet106501043',
        'CSRF1034377807751095': 'identitytheft100476181',
        'CSRF1025162750917734': 'spy103619647',
        'CSRF1084609853366834': 'Trojan105684296',
        'CSRF1062250102795577': 'infiltrator102682581',
        'CSRF1036954146235748': 'subnet108741174',
        'CSRF1070228926725691': 'infiltrator101655359',
        'CSRF1012726089203031': 'geek103699871',
        'CSRF1057370162294326': 'infiltrator100898912',
        'CSRF1071524921926269': 'CSRF104791502',
        'CSRF1069736650323106': 'bot102679558',
        'CSRF1053992552845268': 'phisher107429337',
        'CSRF1072228586864001': 'scammer110856167',
        'CSRF1072589852045715': 'thief109862314',
        'CSRF1018281597743757': 'subnet_ip109914022',
        'CSRF1104742883280212': 'computer107851254',
        'CSRF1017303942805450': 'bogey105050753',
        'CSRF98434202308098797932': 'tech101041640',
        'jn': 'JS aan, T aangeroepen, CSRF aangepast',
        'domain': domain,
        'lol-stop-reverse-engineering-my-source-and-buy-an-api-key': '7e1408afd06beca934ee57afcc3e15a48f551f65',
        'scan_subdomains': ''
    }

    response = requests.post('https://subdomainfinder.c99.nl/', headers=headers, data=data)

    line = next(line for line in response.text.splitlines() if "value='https://subdomainfinder.c99.nl/scans/" in line)

    pattern = r"(?<=value=')[^']*"
    match = re.search(pattern, line)
    if match:
        value = match.group(0)
        response2 = requests.get(f'{value}', headers=headers, data=data)
        subdomain_lines = re.findall(r'<a href="//subdomainfinder.c99.nl/scans/([^"]*)"', response2.text)
        for subdomain_line in subdomain_lines:
            c_url = f'https://subdomainfinder.c99.nl/scans/{subdomain_line}'
            if re.search(f"{domain}", c_url):
                output_set.add(c_url)  # Add unique lines to the set

# Print unique lines
for line in output_set:
    response3 = requests.get(f'{line}')
    outputs = re.findall(fr"//[^']*{domain}", response3.text)
    for output in outputs:
        result = output.split("/")[2]
        if f"{domain}" in result:
            print(result)