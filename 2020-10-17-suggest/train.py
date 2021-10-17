#!/usr/bin/python3

import json
import re
import collections
import sys

def tgramm(s):
    s = '^' + s  # BOW marker
    for i in range(len(s) - 3 + 1):
        yield s[i:i+3]

def main():
    _, data_file, model_file, text_file = sys.argv
    with open(data_file) as fh:
        data = json.load(fh)
    model = {}
    for cls, texts in data.items():
        tag = collections.Counter()
        for text in texts + [cls]:
            text = text.lower()
            text = text.replace('ё', 'е')
            for wrd in set(re.sub('[^а-яa-z]+', ' ', text).split()):
                if len(wrd) < 2:
                    continue
                tag.update(tgramm(wrd))
        model[cls] = tag
    with open(text_file, 'x') as tfh, open(model_file, 'x') as mfh:
        json.dump(data, tfh, sort_keys=True, indent=4, ensure_ascii=False)
        json.dump(model, mfh, sort_keys=True, indent=4, ensure_ascii=False)

if __name__ == '__main__':
    main()
