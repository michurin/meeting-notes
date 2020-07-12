#!/usr/bin/env python3


from collections import namedtuple

import numpy as np

from sklearn import datasets
from sklearn.linear_model import LogisticRegression


# ----- common structures data, shared between Prod/Go and Train/Python wolds


Model = namedtuple('Model', ('wo', 'w'))


# ----- part on the Prod/GoLang side


def std_log(x):
    return 1 / (1 + np.exp(-x))


def fraud_proba(model, x):
    p = std_log(np.sum(model.w * x) + model.wo)
    return p


# ---- on the Train/Python side


def main():
    # 1. prepare train data
    x, y_orig = datasets.load_digits(return_X_y=True)  # 64 features
    y = (y_orig == 8).astype(np.int)  # '8' is our fraudstars

    # 2. train model
    clf = LogisticRegression(solver='newton-cg')
    clf.fit(x, y)

    # (you can overview results here)
    if not True:
        for n, (proba, value) in enumerate(np.vstack([clf.predict_proba(x)[:, 1], y_orig]).T):
            print('{:5d} {:7.5f} {:1.0f} {}'.format(n, proba, value, '*' * int(proba*20)))

    # 3. prepare model for GoLand world
    model = Model(clf.intercept_[0], clf.coef_[0, :])

    # (demo)
    vector1 = x[38]
    vector2 = x[123]
    for v in vector1, vector2:
        p_native = clf.predict_proba(v[np.newaxis, :])[0, 1]
        p_prod = fraud_proba(model, v)
        print(f'predictions native/manual: {p_native:7.5f} / {p_prod:7.5f}')


if __name__ == '__main__':
    main()
