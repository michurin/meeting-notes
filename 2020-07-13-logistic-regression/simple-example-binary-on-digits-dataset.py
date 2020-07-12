#!/usr/bin/env python3

import numpy as np

from sklearn import datasets
from sklearn.linear_model import LogisticRegression
# In this particular case scaling
# improves precision and accuracy
# but reduces recall
# from sklearn.preprocessing import StandardScaler


def main():
    """
    Summary looks like
    recall precis accur  params
    0.8103 0.8924 0.9722 solver='newton-cg', max_iter=10000
    0.8046 0.8805 0.9705 solver='liblinear', max_iter=10000
    0.7989 0.8797 0.9699 solver='sag', max_iter=10000
    """
    meta_params = [
        dict(solver='saga', max_iter=10000),  # L2 — default
        dict(solver='saga', max_iter=10000, penalty='l1'),
        # TODO try saga+elasticnet
        dict(solver='newton-cg', max_iter=10000),  # it supports L2 only
        dict(solver='lbfgs', max_iter=10000),  # L2 only
        dict(solver='liblinear', max_iter=10000),
        dict(solver='sag', max_iter=10000),
    ]

    print('Loading...')
    x, y = datasets.load_digits(return_X_y=True)  # 64 features
    # x = StandardScaler().fit_transform(x)  # it's better to normolize data, however, we can start without it
    y = (y == 8).astype(np.int)  # '8' is our fraudstars

    summary = []

    for n, model_params in enumerate(meta_params):
        clf = LogisticRegression(**model_params)
        print(f'Training... ({n+1}/{len(meta_params)})')
        clf.fit(x, y)
        print(f'Training results: casses={clf.classes_}; coeffs={", ".join(map(str, clf.coef_.ravel()[:3]))}, ...; coeffs.shape={clf.coef_.ravel().shape}')
        p = clf.predict(x)
        # yp = np.vstack((y, p)).T
        # print(yp)
        tp = np.sum(np.logical_and(p == 1, y == 1))
        fp = np.sum(np.logical_and(p == 1, y == 0))
        tn = np.sum(np.logical_and(p == 0, y == 0))
        fn = np.sum(np.logical_and(p == 0, y == 1))
        recall = tp/(tp+fn)
        precision = tp/(tp+fp)
        accuracy = (tp+tn)/(tp+tn+fp+fn)  # clf.score(x, y)
        res_str = f'TP={tp}, FP={fp}, TN={tn}, FN={fn}, Recall={recall} (tp/(tp+fn)), Precision={precision} (tp/(tp+fp)), Accuracy={accuracy} ((tp+tn)/(tp+tn+fp+fn))'
        print(res_str)
        summary.append((model_params, recall, precision, accuracy))

    print('=== summary ===')
    print('recall precis accur  params')
    for params, recall, prec, acc in summary:
        print(f'{recall:6.4f} {prec:6.4f} {acc:6.4f} {", ".join((f"{k}={v!r}" for k, v in params.items()))}')


if __name__ == '__main__':
    main()
