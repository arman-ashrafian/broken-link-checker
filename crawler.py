# Crawl through computing paths and get the URL
# for every element with a src attribute. 

import json
from pprint import pprint

def readDatabase(databaseName):
    links = []
    with open(databaseName, encoding='utf-8') as f:
        data = json.load(f)
    
        
    for x in data:

        for y in data[x]:
            link = None
            if isinstance(y, dict):
                link = y.get('link')
            if link: print(link)

            for z in y:
                link2 = None
                if isinstance(z, dict):
                    link2 = z.get('link')
                if link2: print(link2)



    return links

links = readDatabase('database.json')
print(links)
for i in range(len(links)):
    print(i, ' ', end="")
    print(links[i])
