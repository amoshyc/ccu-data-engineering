import pickle
from tqdm import tqdm

with open('data/words.pkl', 'rb') as f:
    words = pickle.load(f)

with open('data/freq.pkl', 'rb') as f:
    freq = pickle.load(f)


def ngrams(n, res):
    with open('data/text.txt', 'r') as f:
        for line in f:
            for i in range(0, len(line) - n):
                if line[i:i + n] in res:
                    res[line[i:i + n]] += 1


unigram = [w for w in words if freq[w] > 200 and len(w) == 1]
bigrams = [w1 + w2 for w1 in unigram for w2 in unigram]

print('#unigram:', len(unigram))
print('#bigrams:', len(bigrams))

res = dict()
for k in unigram:
    res[k] = 0
for k in bigrams:
    res[k] = 0

ngrams(1, res)
ngrams(2, res)

with open('data/ngrams.pkl', 'wb') as f:
    pickle.dump(res, f)