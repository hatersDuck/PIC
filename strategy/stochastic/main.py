"""
    Стратегия на основе стохастического осциллятора:
    Если значение %K становится больше 80, то мы можем продать (SELL) активы, 
    так как они считаются перекупленными. 
    Если значение %K меньше 20, то мы можем купить (BUY) активы, 
    так как они считаются перепроданными.
"""

import os
import sys

path = os.path.abspath(os.path.join(os.path.dirname(__file__), '..', '..'))
sys.path.append(path)

from strategy.module_analData import Analysis

cfg = [
    {
        'interval': '1m',
        'stochastic_k': 5,
        'stochastic_d': 3,
        'wait': [90, 10],
    },
    {
        'interval': '15m',
        'stochastic_k': 14,
        'stochastic_d': 3,
        'wait': [80, 20],
    },
    {
        'interval': '1h',
        'stochastic_k': 21,
        'stochastic_d': 4,
        'wait': [75, 25],
    },
    {
        'interval': '8h',
        'stochastic_k': 30,
        'stochastic_d': 4,
        'wait': [70, 30],
    },
    {
        'interval': '1d',
        'stochastic_k': 60,
        'stochastic_d': 5,
        'wait': [80, 20],
    }
]

if __name__ == "__main__":
    buy = 0
    sell = 0
    skip = 0
    
    for c in cfg:
        al = Analysis(**c)
        rsi = al.cci
        if c['wait'][0] <= rsi[-1]:
            buy += 1
        elif c['wait'][1] >= rsi[-1]:
            sell += 1
        else:
            skip += 1
    
    if buy > sell and buy > skip:
        print("BUY", end="")
    elif buy < sell and sell > skip:
        print("SELL", end="")
    else:
        print("SKIP", end="")
    
        

