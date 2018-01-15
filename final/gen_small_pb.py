import math
import pickle
from tqdm import tqdm
from collections import defaultdict

with open('data/words.pkl', 'rb') as f:
    words = pickle.load(f)

with open('data/freq.pkl', 'rb') as f:
    freq = pickle.load(f)

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


unigram = set([w for w in words if len(w) == 1 and freq[w] > 300])
bigrams = set([
    w for w in words
    if len(w) == 2 and freq[w] > 500 and w[0] in unigram and w[1] in unigram
])

print('#unigram:', len(unigram))
print('#bigrams:', len(bigrams))

cnt = dict()
ttl = 0
for w in unigram:
    cnt[w] = freq[w]
    ttl += freq[w]
for w in bigrams:
    cnt[w] = freq[w]

pb1 = dict()
for w in unigram:
    pb1[w] = math.log(cnt[w] / ttl)
pb2 = dict()
for w in bigrams:
    pb2[w] = math.log(cnt[w] / cnt[w[1]])

res = defaultdict(list)
for w, _ in pb1.items():
    res[get_keys(w)].append(w)
for ks, _ in res.items():
    res[ks].sort(key=lambda x: -pb1[x])

with open('data/pb1.pkl', 'wb') as f:
    pickle.dump(pb1, f)

with open('data/pb2.pkl', 'wb') as f:
    pickle.dump(pb2, f)

with open('data/small_table.pkl', 'wb') as f:
    pickle.dump(res, f)