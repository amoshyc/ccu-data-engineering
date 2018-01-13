import pickle
from collections import defaultdict

import pandas as pd
from tqdm import tqdm

with open('data/dict.pkl', 'rb') as f:
    pron = pickle.load(f)


def get_first_syllable(pron):
    if pron[0] == '˙':
        return pron[1] + '0'
    else:
        idx = 'ˊˇˋ'.find(pron[-1])
        if idx == -1:
            return pron[0] + '1'
        return pron[0] + str(idx + 2)


def get_keys(word):
    if word in pron:
        keys = [get_first_syllable(char) for char in pron[word]]
    else:
        if any((char not in pron) for char in word):
            return None
        keys = [get_first_syllable(pron[char][0]) for char in word]
    return ''.join(keys)


df = pd.read_csv(
    'data/essay.txt', sep='\t', header=None, names=['word', 'freq'])

essay_words = set(df['word'])
cdict_words = set(pron.keys())
words = essay_words | cdict_words

with tqdm(total=len(words)) as pbar:
    table = defaultdict(list)  # 碼表
    rev = dict()  # 反查表
    for word in words:
        keys = get_keys(word)
        if keys:
            table[keys].append(word)
            rev[word] = keys
        pbar.update()

with tqdm(total=len(df) + len(cdict_words)) as pbar:
    freq = {}  # 詞頻表
    for _, r in df.iterrows():
        freq[r['word']] = int(r['freq'])
        pbar.update()
    for word in cdict_words:
        if word not in freq:
            freq[word] = 0
        pbar.update()

with open('data/table.pkl', 'wb') as f:
    pickle.dump(table, f)

with open('data/rev.pkl', 'wb') as f:
    pickle.dump(rev, f)

with open('data/freq.pkl', 'wb') as f:
    pickle.dump(freq, f)

with open('data/words.pkl', 'wb') as f:
    pickle.dump(words, f)