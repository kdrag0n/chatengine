#!/usr/bin/env python3

import sys
import shutil
import random
from datetime import datetime
import bson
from tqdm import tqdm

# {"in_response_to": [{"text": "Although practicality beats purity.", "occurrence": 1}], "created_at": 1487420546.689, "extra_data": {}, "occurrence": 1}
def main(db_path='chat.bson'):
    'Main function'
    fpath = ' '.join(sys.argv[1:])

    with open(fpath, 'rb') as handle:
        log_data = handle.read().decode('utf-8').split('\n')

    with open(db_path, 'rb') as handle:
        data = bson.loads(handle.read())

    if not data:
        data = []

    for idx, line in tqdm(enumerate(log_data), total=len(log_data)):
        matches = [m for m in data if m['text'] == line]
        if matches:
            line_obj = matches[0]
            line_obj['occurrence'] += 1
            if idx != 0:
                irt_line = log_data[idx - 1]
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
                'extra_data': [],
                'occurrence': 1
            }
            if idx != 0:
                irt_obj = {
                    'text': log_data[idx - 1],
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
