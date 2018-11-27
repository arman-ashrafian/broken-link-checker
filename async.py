#!/usr/bin/python3.7

import asyncio
import time
import aiohttp

async def fetch(url, session):
    async with session.get(url) as resp:
        resp = await resp.read()
        print(url)
        return resp

async def fetch_all(urls):
    tasks = []
    async with aiohttp.ClientSession() as sess:
        for url in urls:
            task = asyncio.ensure_future(fetch(url, sess))
            tasks.append(task)
            await asyncio.gather(*tasks)


def main():
    urls = [
        'http://python.org',
        'http://google.com',
        'http://youtube.com',
        'http://news.ycombinator.com'
    ]

    loop = asyncio.get_event_loop()
    future = asyncio.ensure_future(fetch_all(urls))
    loop.run_until_complete(future)

main()
