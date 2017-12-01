import json
import pathlib
import multiprocessing as mp
import subprocess


def ngrams(data):
    line_id, lines, terms = data
    result = {term: [] for term in terms}

    for i, line in enumerate(lines):
        N = len(line)
        s, t = 0, min(7, N)

        while s < N - 1:
            query = line[s:t]
            # print(s, t, N, query)

            if query in result:
                result[query].append((line_id + i, s))
                s, t = t, min(t + 7, N)
            else:
                if t - s == 2:
                    s, t = s + 1, min(s + 8, N)
                else:
                    t = t - 1

    return result


def sp_py(term_path, text_path, output_path=None):
    with open(term_path, 'r') as f:
        terms = f.read().split('\n')
    with open(text_path, 'r') as f:
        text = f.read().split('\n')
    
    result = ngrams((0, text, terms))
    result = [{'Term': term, 'Pos': pos} for term, pos in result.items()]
    cmp_func = lambda x: (-len(x['Pos']), x['Term'])
    result = sorted(result, key=cmp_func)

    if output_path:
        output = pathlib.Path(output_path)
        output.parent.mkdir(parents=True, exist_ok=True)
        json.dump(result, output.open('w'), ensure_ascii=False)

    return result


def mp_py(term_path, text_path, output_path=None):
    with open(term_path, 'r') as f:
        terms = f.read().split('\n')

    def chunk_data(chunk_size=100000):
        with open(text_path, 'r') as f:
            eof = False
            chunk_cnt = 0

            while not eof:
                chunk = []
                for _ in range(chunk_size):
                    line = f.readline()
                    chunk.append(line)
                    if line == '':
                        eof = True
                        break

                yield chunk_cnt * chunk_size, chunk, terms
                chunk_cnt += 1
                chunk.clear()

    merged = {term: [] for term in terms}
    with mp.Pool(processes=4) as pool:
        chunk_results = pool.imap(ngrams, chunk_data())

        for res in chunk_results:
            for k, v in res.items():
                merged[k] += v

    result = [{'Term': term, 'Pos': pos} for term, pos in merged.items()]
    cmp_func = lambda x: (-len(x['Pos']), x['Term'])
    result = sorted(result, key=cmp_func)

    if output_path:
        output = pathlib.Path(output_path)
        output.parent.mkdir(parents=True, exist_ok=True)
        json.dump(result, output.open('w'), ensure_ascii=False)

    return result

def st_go(term_path, text_path, output_path):
    cmd = './st_go/st_go {} {} {}'.format(term_path, text_path, output_path)
    subprocess.run(cmd, shell=True)

def mt_go(term_path, text_path, output_path):
    cmd = './mt_go/mt_go {} {} {}'.format(term_path, text_path, output_path)
    subprocess.run(cmd, shell=True)


if __name__ == '__main__':
    from pprint import pprint
    term_path = '../data/pu/term.txt'
    text_path = '../data/pu/doc.txt'
    # term_path = './tmp_term.txt'
    # text_path = './tmp_text.txt'
    output_path = './out.json'
    sp_py(term_path, text_path, output_path)
    # mt_go(term_path, text_path, output_path)