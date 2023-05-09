"""
    Стратегия на основе RSI:
    Если значение RSI становится больше 70, то мы можем продать (SELL) активы, 
    так как они считаются перекупленными. 
    Если значение RSI меньше 30, то мы можем купить (BUY) активы, 
    так как они считаются перепроданными.
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
        'rsi_period': [9, 14][random.randint(0,1)],
        'wait': [80, 20],
    },
    {
        'interval': '15m',
        'rsi_period': [14, 21][random.randint(0,1)],
        'wait': [70, 30],
    },
    {
        'interval': '1h',
        'rsi_period': [21, 30][random.randint(0,1)],
        'wait': [65, 35],
    },
    {
        'interval': '8h',
        'rsi_period': [30, 60][random.randint(0,1)],
        'wait': [60, 40],
    },
    {
        'interval': '1d',
        'rsi_period': [60, 90][random.randint(0,1)],
        'wait': [70, 30],
    }
]

if __name__ == "__main__":
    buy = 0
    sell = 0
    skip = 0
    
    for c in cfg:
        al = Analysis(**c)
        k, _, _ = al.stochastic_oscillator
        if c['wait'][0] <= k:
            buy += 1
        elif c['wait'][1] >= k:
            sell += 1
        else:
            skip += 1
    
    if buy > sell and buy > skip:
        print("BUY", end="")
    elif buy < sell and sell > skip:
        print("SELL", end="")
    else:
        print("SKIP", end="")
    
        

