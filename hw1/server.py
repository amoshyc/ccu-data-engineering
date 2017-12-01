import time
import json
import pathlib

import tornado.ioloop
import tornado.web

import algorithms


class IndexHandler(tornado.web.RequestHandler):
    def get(self):
        self.render('index.html')


class DataHandler(tornado.web.RequestHandler):
    def make_file(self, data, path):
        try:
            is_file = pathlib.Path(data).exists()
        except:
            is_file = False
        
        if not is_file:
            with open(path, 'w') as f: 
                f.write(data)
        else:
            path = data
        
        return is_file, pathlib.Path(path)

    def post(self):
        text = self.get_argument('text', '')
        term = self.get_argument('term', '')
        use_go = self.get_argument('use', '') == 'true'
        alg = algorithms.mt_go if use_go else algorithms.mp_py

        is_text_file, text_path = self.make_file(text, './tmp_text.txt')
        is_term_file, term_path = self.make_file(term, './tmp_term.txt')
        file_only = is_text_file
        out_path = pathlib.Path('./out.json') 

        start_time = time.time()
        alg(term_path, text_path, out_path)
        end_time = time.time()

        t = end_time - start_time

        if not file_only:
            json_data = json.load(out_path.open('r'))
            res = []
            for item in json_data:
                res.append([item['Term'], item['Pos']])
        else:
            res = []

        ret = {
            'time': '{:.3f}s'.format(t),
            'file': str(out_path),
            'file_only': is_text_file,
            'res': res
        }
        print('finish reading')
        self.write(ret)


def main():
    route = [
        (r'/', IndexHandler),
        (r'/data', DataHandler),
        (r'/static/(.*)', tornado.web.StaticFileHandler, { 'path': './static' }),
    ] # yapf: disable

    app = tornado.web.Application(route, debug=True)

    app.listen(8787)
    tornado.ioloop.IOLoop.current().start()


if __name__ == "__main__":
    main()