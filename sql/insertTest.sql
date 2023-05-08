-- Заполнение таблицы tradeStrategy
INSERT INTO tradeStrategy (title, path_file, status, symbol)
VALUES
  ('NOSTRATEGY', '/', 'off', 'RUBRUB'),
  ('ADX', '/strategy/adx/main.py', 'active', 'BTCUSDT'),
  ('RSI', '/strategy/rsi/main.py', 'pause', 'BTCUSDT'),
  ('stochastic', '/strategy/stochastic/main.py', 'off', 'BTCUSDT');
