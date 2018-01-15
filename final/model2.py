import pickle
import readline
import heapq
import time
from copy import copy
from pprint import pprint
from collections import defaultdict


with open('data/rev.pkl', 'rb') as f:
    rev = pickle.load(f)

with open('data/small_table.pkl', 'rb') as f:
    cands = pickle.load(f)

with open('data/pb1.pkl', 'rb') as f:
    pb1 = pickle.load(f)

with open('data/pb2.pkl', 'rb') as f:
    pb2 = pickle.load(f)


def split2keys(inp):
    return [inp[i:i + 2] for i in range(0, len(inp), 2)]


def build_graph(keys):
    g = defaultdict(list)

    S, T = '$', '^'
    cur = [S]

    for k in keys:
        nxt = set()

        for u in cur:
            for c in cands[k][:30]:
                if u == '$' or (u + c) not in pb2:
                    word = c
                    pb = pb1[word]
                else:
                    word = u + c
                    pb = pb2[word]

                g[u].append((c, pb))
                nxt.add(c)

        cur = list(nxt)
        print(cur)

    for u in cur:
        g[u].append((T, 0.0))

    # pprint(g)
    print(len(g))
    return g


def longest_path(g):
    dis = defaultdict(lambda: int(-1e9))
    prev = defaultdict(lambda: '$')
    name = dict()
    pq = list()

    S = '$'
    dis[S] = 0
    heapq.heappush(pq, (0, S))

    while not len(pq) == 0:
        d, u = heapq.heappop(pq)
        if d < dis[u]:
            continue

        for v, w in g[u]:
            if dis[v] < dis[u] + w:
                dis[v] = dis[u] + w
                prev[v] = u
                heapq.heappush(pq, (dis[v], v))

    u = prev['^']
    res = list()
    while u != '$':
        res.append(u)
        u = prev[u]

    return ''.join(res[::-1])


def query_keys(inp):
    if inp[0] == 'r':
        inp = inp[1:]
        return ''.join(rev[c] for c in inp)
    else:
        keys = split2keys(inp)
        g = build_graph(keys)
        res = longest_path(g)
        return res


# print(query_keys('ㄊ1ㄑ4ㄅ4ㄘ4'))
# print(query_keys('ㄓ1ㄨ2ㄕ1ㄖ4ㄈ3'))
# print(query_keys('ㄌ2ㄙ3ㄕ3ㄒ2'))
# print(query_keys('r離散數學'))

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