#!/usr/bin/env python3

import sys
import os
from fnmatch import filter
import shutil
import random
from datetime import datetime
import json
from chatterbot_corpus import __file__ as _cp__file__
import bson
from tqdm import tqdm

def rc_files(folder):
    rl = set()
    for root, fds, files in os.walk(folder):
        for fn in files:
            rl.add(os.path.join(root, fn))
    return rl

# {"in_response_to": [{"text": "Although practicality beats purity.", "occurrence": 1}], "created_at": 1487420546.689, "extra_data": {}, "occurrence": 1}
def main(db_path='chat.bson'):
    'Main function'

    corpus_pairs = []
    corpus_data_path = os.path.join(os.path.dirname(os.path.abspath(_cp__file__)),
                                    'data')
    for fpath in filter(rc_files(corpus_data_path), '*.corpus.json'):
        with open(fpath, 'rb') as handle:
            print(fpath)
            corpus_data = json.loads(handle.read().decode('utf-8'))
            for cp in corpus_data:
                corpus_pairs.extend(corpus_data[cp])

    with open(db_path, 'rb') as handle:
        data = bson.loads(handle.read())

    for pair in tqdm(corpus_pairs):
        irt_line = pair[0]
        line = pair[1]
        matches = [m for m in data if m['text'] == line]
        if matches:
            line_obj = matches[0]
            line_obj['occurrence'] += 1
            if irt_line in [irt['text'] for irt in line_obj['in_response_to']]:
                for irt in line_obj['in_response_to']:
                    if irt['text'] == line:
                        irt['occurrence'] += 1
            else:
                irt_obj = {
                    'text': irt_line,
                    'occurrence': 1
                }
                line_obj['in_response_to'].append(irt_obj)
        else:
            line_obj = {
                'text': line,
                'in_response_to': [],
                'created_at': datetime.now().timestamp(),
                'extra_data': {},
                'occurrence': 1
            }
            irt_obj = {
                'text': irt_line,
                'occurrence': 1
            }
            line_obj['in_response_to'].append(irt_obj)
            data.append(line_obj)

    print('Training finished!')
    print('Writing data atomically...')
    atom_name = db_path + '.t_asav' + str(random.randint(300, 4000))
    with open(atom_name, 'wb+') as handle:
        handle.write(bson.dumps(data))
    shutil.move(atom_name, db_path)
    print('Finished! Exiting...')
    exit(0)


if __name__ == '__main__':
    main()
