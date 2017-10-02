#!/usr/bin/env python3
import requests
import random
import json

f = open('keys.json', 'rb')
keydat = json.loads(f.read().decode('utf-8'))

head = {'Authorization': keydat['api'][1]}

post_data = {
    'session': 'testclientpy-' + str(random.randint(1, 9223372036854775800)),
    'query': ''
}

while True:
    inp = input('> ')
    post_data['query'] = inp
    print(requests.post('http://localhost:2083/ask', headers=head, data=json.dumps(post_data)).text)
