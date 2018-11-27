#!/usr/bin/python3.7

import json
import time


class Link:
    def __init__(self, url, parents):
        self.url = url
        self.parents = parents
    
    def __repr__(self):
        return '%s\nparents: %s\n' % (self.url, str(self.parents))

class Crawler:
    def __init__(self, databaseName):
        self.data = []
        self.links = []
        self.databaseName = databaseName
        self.readDatabase()
        # get links from DB and add to links[]
        self.getLinks_Stories(self.data)
        self.getLinks_Majors(self.data) 
        self.getLinks_Departments(self.data)
        self.getLinks_Resources(self.data)
        self.getLinks_ResourceBanner(self.data)
        self.getLinks_Projects(self.data)
        self.getLinks_Orgs(self.data)


    def getLinks_Stories(self, dic):
        # get links in stories table
        i = 0
        for x in dic['stories']:
            if isinstance(x, dict):
                if x['link']:
                    l = Link(x['link'], ['stories', i, 'link'])
                    self.links.append(l)
            i += 1
    
    def getLinks_Majors(self, dic):
        # get links in majors table
        i = 0
        j = 0
        for x in dic['majors']:
            # check image links
            l = Link(x.get('image'), ['majors', i, 'image'])
            self.links.append(l)
            if x.get('moreInfo'):
                j = 0
                for d in x['moreInfo']:
                    # check links in moreInfo
                    l = Link(d.get('link'), ['majors', i,'moreInfo', j, 'link'])
                    self.links.append(l)
                    j += 1
            i += 1

    def getLinks_Departments(self, dic):
        # get links in departments table
        for x in dic['departments']:
            l = Link(dic['departments'][x]['link'], ['departments', x, 'link'])
            self.links.append(l)

    def getLinks_Resources(self, dic):
        # get links in resources table
        i = 0
        for x in dic['resources']:
            if x.get('mapImage'):
                l = Link(x.get('mapImage'), ['resources', i, 'mapImage'])
                self.links.append(l)
            if x.get('mapLink'):
                l = Link(x.get('mapLink'), ['resources', i, 'mapLink'])
                self.links.append(l)
            if x.get('link'):
                l = Link(x.get('link'), ['resources', i, 'link'])
                self.links.append(l)
            i += 1 
        
    def getLinks_ResourceBanner(self, dic):
        # get links in resourceBanner table
        i = 0
        for x in dic['resourceBanner']:
            img = x.get('image')
            link = x.get('link')
            if img: 
                l = Link(img, ['resourceBanner', i, 'image'])
                self.links.append(l)
            if link: 
                l = Link(img, ['resourceBanner', i, 'link'])
                self.links.append(l)
            i += 1

    def getLinks_Projects(self, dic):
        # get links in projects table
        i = 0
        j = 0
        for x in dic['projects']:
            j = 0
            for v in x.get('videos'):
                l = Link(v, ['projects', i, 'videos', j])
                self.links.append(l)
                j += 1
            if x.get('link'):
                l = Link(x.get('link'), ['projects', i, 'link'])
                self.links.append(l)
            i += 1

    def getLinks_Orgs(self, dic):
        # get links in orgs table
        i = 0
        for x in dic['orgs']:
            linkurl = x.get('link')
            if linkurl:
                l = Link(linkurl, ['orgs', i, 'link'])
                self.links.append(l)
            i += 1

    def readDatabase(self):
        with open(self.databaseName, encoding='utf-8') as f:
            self.data = json.load(f)

def main():
    c = Crawler('database.json')
    f = open("links.txt", "w")
    for l in c.links:
        f.write(l.url)
        f.write('\n')

if __name__ == '__main__':
    main()

