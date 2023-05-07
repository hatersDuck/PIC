import binance
import numpy as np
import pandas as pd
import datetime 

class Analysis():
    def __init__(self, **kwargs) -> None:
        self.interval = kwargs.get('interval', '1h')
        self.symbol = kwargs.get('symbol', 'BTCUSDT')

        self.rsi_period = kwargs.get('rsi_period', 14)
        self.stochastic_period = {
            'n': kwargs.get('stochastic_n', 14),
            'k': kwargs.get('stochastic_k', 3),
            'd': kwargs.get('stochastic_d', 3),
        }
        self.directional_period = {
            'n': kwargs.get('directional_n', 14),
            'm': kwargs.get('directional_m', 14),
            'k': kwargs.get('directional_k', 14),
        }
        self.cci_period = kwargs.get('cci_period', 20)
        self.adx_period = kwargs.get('adx_period', 14)

        self.prices = binance.Client().get_klines(symbol=self.symbol, interval=self.interval)
        self.close_prices = np.array([float(close[4]) for close in self.prices], dtype=float)
        self.high_prices = np.array([float(high[2]) for high in self.prices], dtype=float)
        self.low_prices = np.array([float(low[3]) for low in self.prices], dtype=float)
        self.dates = np.array([datetime.datetime.fromtimestamp(date[6]/1000) for date in self.prices])

    @property
    def rsi(self, period = None):
        """
        Вычисляет RSI (Relative Strength Index)

        :param period: Период для вычисления среднего значения (int).
        :return: Массив значений RSI (float)
        """
        if period is None:
            period = self.rsi_period
        
        prices = self.close_prices
        diffs = np.diff(prices)

        positive_diff = np.where(diffs > 0, diffs, 0)
        negative_diff = np.where(diffs < 0, -diffs, 0)

        average_gain = np.convolve(positive_diff, np.ones((period,)) / period, mode='valid')
        average_loss = np.convolve(negative_diff, np.ones((period,)) / period, mode='valid')

        rs = np.divide(average_gain, average_loss, out=np.zeros_like(average_gain), where=average_loss != 0)
        rsi = 100 - 100 / (1 + rs)

        # Создаем массив значений RSI
        rsi_values = np.full(len(prices), np.nan)
        rsi_values[period:] = rsi

        return rsi_values.tolist()

    @property
    def stochastic_oscillator(self, n = None, k_period = None, d_period = None):
        """
        Рассчитывает стохастический осциллятор для заданного набора цен

        :param prices: Список или массив исторических цен
        :param n: период для скользящего среднего
        :param k_period: Количество периодов, используемых для расчета %K. По умолчанию равно 14
        :param d_period: Количество периодов, используемых для расчета %D. По умолчанию равно 3
        :return: Кортеж, содержащий значения %K и %D в виде массивов numpy
        """

        # Преобразуем цены в массив numpy
        
        if n is None:
            n = self.stochastic_period['n']
        if k_period is None:
            k_period = self.stochastic_period['k']
        if d_period is None:
            d_period = self.stochastic_period['d']

        prices = self.close_prices
        high_prices = np.max(prices[-n:])
        low_prices = np.min(prices[-n:])
        
        # Расчет %K
        current_close_price = prices[-1]
        k = ((current_close_price - low_prices) / (high_prices - low_prices)) * 100
        
        # Расчет скользящего среднего по %K
        k_ma = np.mean(prices[-k_period:])
        
        # Расчет %D
        d = np.mean(prices[-d_period-k_period:-k_period])
        
        return k, k_ma, d

    @property
    def cci(self, period=None):
        """
        Вычисляет Commodity Channel Index (CCI)

        :param period: Период для вычисления CCI (int).
        :return: Массив значений CCI (float)
        """
        if period is None:
            period = self.cci_period
        
        typical_prices = (self.close_prices + self.high_prices + self.low_prices) / 3
        moving_average = np.mean(typical_prices[-period:])
        mean_deviation = np.mean(np.abs(typical_prices[-period:] - moving_average))

        cci_values = (typical_prices - moving_average) / (0.015 * mean_deviation)

        return cci_values.tolist()

    @property
    def average_directional_index(self, period=None):
        """
        #####Начало массива заполнено Nan#######
        Функция для расчета AVERAGE DIRECTIONAL INDEX (ADX)
        
        period: количество дней для расчета
        
        Возвращает ADX для каждого дня в заданном периоде.
        """
        if period is None:
            period = self.adx_period
        
        close_prices = self.close_prices
        high_prices = self.high_prices
        low_prices = self.low_prices

        plus_dm = np.zeros(len(high_prices))
        minus_dm = np.zeros(len(high_prices))
        
        for i in range(1, len(high_prices)):
            up_move = high_prices[i] - high_prices[i-1]
            down_move = low_prices[i-1] - low_prices[i]
            
            if up_move > down_move and up_move > 0:
                plus_dm[i] = up_move
                
            if down_move > up_move and down_move > 0:
                minus_dm[i] = down_move
        
        # Рассчитываем True Range (TR) и Directional Movement Index (DI)
        tr_list = np.maximum(high_prices[1:] - low_prices[1:], 
                            np.abs(high_prices[1:] - close_prices[:-1]), 
                            np.abs(low_prices[1:] - close_prices[:-1]))
        tr_smoothed = np.concatenate([[np.nan], pd.Series(tr_list).rolling(window=period).sum().values])
        
        plus_di = 100 * pd.Series(plus_dm).ewm(span=period, min_periods=period).mean().values / tr_smoothed
        minus_di = 100 * pd.Series(minus_dm).ewm(span=period, min_periods=period).mean().values / tr_smoothed
        
        # Рассчитываем ADX
        dx = 100 * np.abs((plus_di - minus_di) / (plus_di + minus_di))
        adx = pd.Series(dx).ewm(span=period, min_periods=period).mean().values
        
        return adx

    @property
    def ema(self, values = None, n=10):
        if values is None:
            values = self.close_prices
        
        alpha = 2/(n+1)
        ema = np.zeros_like(values)
        for i in range(len(ema)):
            if i == 0:
                ema[i] = values[i]
            else:
                ema[i] = alpha*values[i] + (1-alpha)*ema[i-1]
        return ema

    @property
    def sma(self, prices = None, window = 3):
        if prices is None:
            prices = self.close_prices
        
        weights = np.repeat(1.0, window) / window
        smas = np.convolve(prices, weights, 'valid')
        return smas

    def get_plot(self, arr):
        df = pd.DataFrame(arr, self.dates)
        df.plot()