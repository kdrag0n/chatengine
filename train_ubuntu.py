#!/usr/bin/env python3

import sys
import csv
import shutil
import glob
import io
import random
from datetime import datetime
import bson
from tqdm import tqdm

# {"in_response_to": [{"text": "Although practicality beats purity.", "occurrence": 1}], "created_at": 1487420546.689, "extra_data": {}, "occurrence": 1}
def main(db_path='chat.bson'):
    'Main function'
    fpath = ' '.join(sys.argv[1:])

    data = []
    try:
        with open(db_path, 'rb') as handle:
            data = bson.loads(handle.read())
    except FileNotFoundError:
        pass

    if not data:
        data = []

    extracted_corpus_path = 'data/dialogs/**/*.tsv'
    print('Discovering files...')
    source_files = []
    for path in tqdm(glob.iglob(extracted_corpus_path)):
        source_files.append(path)

    for path in tqdm(source_files, total=len(source_files)):
        with open(path, 'rb') as tsv_file:
            tsv = io.StringIO(tsv_file.read().decode('utf-8'))
            reader = csv.reader(tsv, delimiter='\t')
            history = []

            for row in reader:
                line = row[3]
                matches = [m for m in data if m['text'] == line]

                if matches:
                    line_obj = matches[0]
                    line_obj['occurrence'] += 1
                    if history:
                        irt_line = history[-1]

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
                        'extra_data': [
                            {
                                'name': 'datetime',
                                'value': row[0]
                            },
                            {
                                'name': 'speaker',
                                'value': row[1]
                            }
                        ],
                        'occurrence': 1
                    }

                    if row[2].strip():
                        line_obj['extra_data'].append({
                            'name': 'addressing_speaker',
                            'value': row[2]
                        })

                    if history:
                        irt_obj = {
                            'text': history[-1],
                            'occurrence': 1
                        }

                        line_obj['in_response_to'].append(irt_obj)

                    data.append(line_obj)

                history.append(line)

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
