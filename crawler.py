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
        for x in dic['majors']:
            # check image links
            self.links.append(x.get('image'))
            if x.get('moreInfo'):
                for d in x['moreInfo']:
                    # check links in moreInfo
                    self.links.append(d.get('link'))
        
        # get links in departments table
        for x in dic['departments'].values():
            self.links.append(x['link'])

        # get links in resources table
        for x in dic['resources']:
            if x.get('mapImage'):
                self.links.append(x.get('mapImage'))
            if x.get('mapLink'):
                self.links.append(x.get('mapLink'))
            if x.get('link'):
                self.links.append(x.get('link'))
        
        # get links in resourceBanner table
        for x in dic['resourceBanner']:
            img = x.get('image')
            link = x.get('link')
            if img: self.links.append(img)
            if link: self.links.append(link)
        
        # get links in projects table
        for x in dic['projects']:
            for v in x.get('videos'):
                self.links.append(v)
            if x.get('link'):
                self.links.append(x.get('link'))
        
        # get links in orgs table
        for x in dic['orgs']:
            link = x.get('link')
            if link: self.links.append(link)

    def readDatabase(self):
        with open(self.databaseName, encoding='utf-8') as f:
            data = json.load(f)
        
        l = self.getLinks(data) 

c = Crawler('database.json')
pprint(c.links)
print(len(c.links))
