import pickle
import readline
import heapq
import time
from copy import copy
from pprint import pprint

with open('data/table.pkl', 'rb') as f:
    table = pickle.load(f)

with open('data/log_pb.pkl', 'rb') as f:
    log_pb = pickle.load(f)

with open('data/freq.pkl', 'rb') as f:
    freq = pickle.load(f)


def split2keys(inp):
    return [inp[i:i + 2] for i in range(0, len(inp), 2)]


def find_max_prob(keys, n_cand=10):
    candidates = []
    tokens = []

    def dfs(s, total):
        nonlocal candidates, tokens

        if s == len(keys):
            heapq.heappush(candidates, (total, copy(tokens)))
            if len(candidates) > n_cand:
                heapq.heappop(candidates)
            return

        for t in range(s + 1, len(keys) + 1):
            key = ''.join(keys[s:t])
            if key not in table:
                continue

            for cand in table[key]:
                if cand not in log_pb:
                    break

                tokens.append(cand)
                dfs(t, total + log_pb[cand])
                tokens.pop()

    dfs(0, 0.0)

    pprint(candidates)

    return ''.join(candidates[0][1])


def query_keys(inp):
    print('query:', inp)
    if inp in table:
        res = table[inp]
        res.sort(key=lambda x: -freq[x])
        res = res[:min(10, len(res))]
        return res
    else:
        keys = split2keys(inp)
        res = find_max_prob(keys)
        return res


# print(query_keys('ㄅ3ㄈ1ㄕ1'))

while True:
    keys = input('> ')
    if keys == '':
        continue
    if keys == 'q':
        break
    start_time = time.time()
    print(query_keys(keys))
    end_time = time.time()
    print('Time used: {:.1f}s'.format(end_time - start_time))