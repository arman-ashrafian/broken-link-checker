#!/usr/bin/python3.7

import json
import time
import os
import grequests as gr # for concurrent http requests

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
        self.getLinks_Stories()
        self.getLinks_Majors() 
        self.getLinks_Departments()
        self.getLinks_Resources()
        self.getLinks_ResourceBanner()
        self.getLinks_Projects()
        self.getLinks_Orgs()


    def getLinks_Stories(self):
        # get links in stories table
        i = 0
        for x in self.data['stories']:
            if isinstance(x, dict):
                if x['link']:
                    l = Link(x['link'], ['stories', i, 'link'])
                    self.links.append(l)
            i += 1
    
    def getLinks_Majors(self):
        # get links in majors table
        i = 0
        j = 0
        for x in self.data['majors']:
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

    def getLinks_Departments(self):
        # get links in departments table
        for x in self.data['departments']:
            l = Link(self.data['departments'][x]['link'], ['departments', x, 'link'])
            self.links.append(l)

    def getLinks_Resources(self):
        # get links in resources table
        i = 0
        for x in self.data['resources']:
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
        
    def getLinks_ResourceBanner(self):
        # get links in resourceBanner table
        i = 0
        for x in self.data['resourceBanner']:
            img = x.get('image')
            link = x.get('link')
            if img: 
                l = Link(img, ['resourceBanner', i, 'image'])
                self.links.append(l)
            if link: 
                l = Link(img, ['resourceBanner', i, 'link'])
                self.links.append(l)
            i += 1

    def getLinks_Projects(self):
        # get links in projects table
        i = 0
        j = 0
        for x in self.data['projects']:
            j = 0
            for v in x.get('videos'):
                l = Link(v, ['projects', i, 'videos', j])
                self.links.append(l)
                j += 1
            if x.get('link'):
                l = Link(x.get('link'), ['projects', i, 'link'])
                self.links.append(l)
            i += 1

    def getLinks_Orgs(self):
        # get links in orgs table
        i = 0
        for x in self.data['orgs']:
            linkurl = x.get('link')
            if linkurl:
                l = Link(linkurl, ['orgs', i, 'link'])
                self.links.append(l)
            i += 1

    def readDatabase(self):
        with open(self.databaseName, encoding='utf-8') as f:
            self.data = json.load(f)

    def requestExceptionHandler(self, request, exception):
        print('%s timeout failure' % request.url)

    def checkLinks(self):
        # make sure all urls are in the form http://...
        BASE_PATH = 'http://computingpaths.ucsd.edu'
        urls = []
        for l in self.links:
            if l.url == '/':
                l.url = BASE_PATH + l.url
            elif l.url[0:4] != 'http':
                l.url = BASE_PATH + '/' + l.url
            urls.append(l.url)
        # remove duplicates
        urls = set(urls)
        
        # unsent requests
        reqs = (gr.get(url) for url in urls)
        responses = gr.map(reqs, exception_handler=self.requestExceptionHandler)
        for r in responses:
            if r == None:
                print("Error")
            elif r.status_code != 200:
                print('[%d] %s' % (r.status_code, r.url))

def main():
    c = Crawler('database.json')
    c.checkLinks()
        
if __name__ == '__main__':
    main()

