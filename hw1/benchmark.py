import time
import subprocess

import algorithms

texts = [
    '../data/10/wiki_00',
    '../data/50/wiki_00',
    '../data/100/wiki_00',
    '../data/pu/doc.txt',
]
term_path = '../data/pu/term.txt'
out_path = '/dev/null'

algs = [
    (algorithms.sp_py, 'sp_py'),
    (algorithms.mp_py, 'mp_py'),
    (algorithms.st_go, 'st_go'),
    (algorithms.mt_go, 'mt_go'),
]

for text_path in texts:
    print(text_path, ':')
    for alg, name in algs:
        start_time = time.time()
        alg(term_path, text_path, out_path)
        end_time = time.time()
        print(name, '\t\t', end_time - start_time)
    print('-' * 50)