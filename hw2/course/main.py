from collections import defaultdict

data = []
while True:
    try:
        data.append(input())
    except:
        break

cnt = defaultdict(int)
for line in data:
    line = line.replace("\t", " ")
    for token in line.split(" "):
        cnt[token] += 1
    
res = list([(-v, k) for k, v in cnt.items()])
res = sorted(res)
for (v, k) in res:
    print('{} {}'.format(-v, k))