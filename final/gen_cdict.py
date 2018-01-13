import re
import pickle
from collections import defaultdict
import pandas as pd
from tqdm import tqdm

symbol_regex = re.compile(r'[，；．˙]', re.UNICODE)
change_regex = re.compile(r'[\(（][變又語讀](.*)', re.UNICODE)
bracket_regex = re.compile(r'[\(（](.*)[\)）]', re.UNICODE)


def process_name(name):
    if '.gif' in name:
        return None
    name = name.replace(' ', '')
    name = re.sub(symbol_regex, '', name)
    name = re.sub(bracket_regex, '', name)
    return name


def process_pron(pron):
    pron = pron.replace('　', ' ')
    pron = re.sub(change_regex, '', pron)
    pron = re.sub(bracket_regex, '', pron)
    res = []
    for token in pron.split(' '):
        if len(token) == 0 or token == '　':
            continue
        if len(token) > 1 and token[-1] == 'ㄦ' and token[-2] != ' ':
            res.append(token[:-1])
            res.append('ㄦ')
        else:
            res.append(token)
    return res


df1 = pd.read_excel('data/cdict1.xls')[['字詞名', '注音一式']]
df2 = pd.read_excel('data/cdict2.xls')[['字詞名', '注音一式']]
df3 = pd.read_excel('data/cdict3.xls')[['字詞名', '注音一式']]
df = pd.concat([df1, df2, df3])
df = df.dropna(axis=0)

with tqdm(ascii=True, total=len(df)) as pbar:
    d = dict()
    for idx, row in df.iterrows():
        name, pron = row['字詞名'], row['注音一式']
        new_name = process_name(name)
        new_pron = process_pron(pron)

        if not new_name or not new_pron:
            continue

        if len(new_name) != len(new_pron):
            print('|{}| => |{}|'.format(name, pron))
            print('|{}| => {}'.format(new_name, new_pron))
            assert False

        if new_name not in d:
            d[new_name] = new_pron

        pbar.update()

with open('data/dict.pkl', 'wb') as f:
    pickle.dump(d, f)