#!/usr/bin/python3

import json
import sys

N = 3

def tgramm(s):
    s = '^' + s
    #s = s[:3] + '1' + s[3:]
    #s = s[:6] + '2' + s[6:]
    for i in range(len(s) - N + 1):
        yield s[i:i+N]

def apply(model, wrd):
    r = []
    for k, v in model.items():
        s = sum(v.get(x, 0) for x in tgramm(wrd))
        r.append((s, k))
    r.sort(reverse=True)
    s = sum(t[0] for t in r)
    return list((e[0]/s, e[1]) for e in r)

def main():
    with open('model.json') as fh:
        model = json.load(fh)
    for x in sys.argv[1:]:
        print(x)
        for r in apply(model, x):
            print('{:4.2f} {}'.format(*r))

if __name__ == '__main__':
    main()
