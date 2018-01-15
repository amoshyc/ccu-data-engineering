import pickle
import readline

with open('data/table.pkl', 'rb') as f:
    table = pickle.load(f)

with open('data/freq.pkl', 'rb') as f:
    freq = pickle.load(f)

with open('data/rev.pkl', 'rb') as f:
    rev = pickle.load(f)


def query_keys(keys):
    if keys[0] == 'r':
        keys = keys[1:]
        if keys in rev:
            return rev[keys]
        else:
            return 'Not in dict'
    elif keys in table:
        res = table[keys]
        res.sort(key=lambda x: -freq[x])
        res = res[:min(10, len(res))]
        return res
    else:
        return 'Not in dict'


def interactive():
    while True:
        keys = input('> ')
        if keys == '':
            continue
        if keys == 'q':
            break
        print(query_keys(keys))


if __name__ == '__main__':
    interactive()