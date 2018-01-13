import math
import pickle
from collections import defaultdict
from tqdm import tqdm

with open('data/words.pkl', 'rb') as f:
    words = pickle.load(f)

with open('data/freq.pkl', 'rb') as f:
    freq = pickle.load(f)


cnt = defaultdict(int)
for word in tqdm(words):
    cnt[len(word)] += freq[word]
for k, v in cnt.items():
    if v == 0:
        cnt[k] = -1e9
    else:
        cnt[k] = math.log(v)

pb = {}
for word in tqdm(words):
    if freq[word] < 50:
        continue
    else:
        pb[word] = math.log(freq[word]) - cnt[len(word)]

with open('data/log_pb.pkl', 'wb') as f:
    pickle.dump(pb, f)