import json
from pprint import pprint

class Crawler:
    def __init__(self, databaseName):
        self.links = []
        self.databaseName = databaseName
        self.readDatabase()

    def getLinks(self, dic):
        # get links in stories table
        for x in dic['stories']:
            if isinstance(x, dict):
                if x['link']:
                    self.links.append(x['link'])

        # get links in majors table


    def readDatabase(self):
        with open(self.databaseName, encoding='utf-8') as f:
            data = json.load(f)
        
        l = self.getLinks(data) 

c = Crawler('database.json')
pprint(c.links)
print(len(c.links))
