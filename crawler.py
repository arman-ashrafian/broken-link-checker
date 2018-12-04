#!/usr/bin/python3.7

import json
import time
import subprocess as sub
import os
import concurrent.futures
import urllib.request

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

    # Retrieve a single page and report the URL and contents
    def load_url(self, url, timeout):
        with urllib.request.urlopen(url, timeout=timeout) as conn:
            return conn.getcode()

    def checkLinks(self):
        BASE_URL = "http://computingpaths.ucsd.edu"
        # create list of full urls
        full_urls = []
        for l in self.links:
            if l.url[0] == '/':
                l.url = BASE_URL + l.url
            elif l.url[0:4] != "http":
                l.url = BASE_URL + "/" + l.url
            full_urls.append(l.url.replace(' ', '%20'))

        # with statement to ensure threads are cleaned up promptly
        with concurrent.futures.ThreadPoolExecutor(max_workers=5) as executor:
            # Start the load operations and mark each future with its URL
            future_to_url = {executor.submit(self.load_url, l, 60): l for l in full_urls}
            for future in concurrent.futures.as_completed(future_to_url):
                url = future_to_url[future]
                try:
                    data = future.result()
                except Exception as exc:
                    print('Error: %r --- %s' % (url, exc))
                #else:
                #    print('%r ---- %d' % (url, data))
         
        

def main():
    c = Crawler('database.json')
    c.checkLinks() 
    
        
if __name__ == '__main__':
    main()

