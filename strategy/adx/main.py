"""
    Стратегия на основе индекса направленного движения (ADX):
    Если значение ADX становится больше 25, то мы можем продать (SELL) активы, 
    так как это может свидетельствовать о начале тренда вниз.
    Если значение ADX меньше 20, то мы можем купить (BUY) активы, 
    так как это может свидетельствовать об отсутствии тренда.
"""

import os
import sys

path = os.path.abspath(os.path.join(os.path.dirname(__file__), '..', '..'))
sys.path.append(path)

from strategy.module_analData import Analysis
import random

cfg = [

    {
        'interval': '1m',
        'adx_period': [7,14][random.randint(0,1)],
        'wait': [20, 30],
    },
    {
        'interval': '15m',
        'adx_period': [14,21][random.randint(0,1)],
        'wait': [20, 25],
    },
    {
        'interval': '1h',
        'adx_period': [21,30][random.randint(0,1)],
        'wait': [15, 20],
    },
    {
        'interval': '8h',
        'adx_period': [30,44][random.randint(0,1)],
        'wait': [20, 25],
    },
    {
        'interval': '1d',
        'adx_period': [60,90][random.randint(0,1)],
        'wait': [20, 25],
    }
]

if __name__ == "__main__":
    buy = 0
    sell = 0
    skip = 0

    for c in cfg:
        al = Analysis(**c)
        adx = al.average_directional_index
        if c['wait'][0] >= adx[-1]:
            buy += 1
        elif c['wait'][1] <= adx[-1]:
            sell += 1
        else:
            skip += 1
    
    if buy > sell and buy > skip:
        print("BUY")
    elif buy < sell and sell > skip:
        print("SELL")
    else:
        print("SKIP")
    
        

