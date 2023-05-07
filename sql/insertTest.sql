-- Заполнение таблицы tradeStrategy
INSERT INTO tradeStrategy (title, path_file, status, symbol)
VALUES
  ('Strategy A', '/strategy/strategy_a.py', 'active', 'BTCUSDT'),
  ('Strategy B', '/strategy/strategy_b.py', 'pause', 'ETHUSDT'),
  ('Strategy C', '/strategy/strategy_c.py', 'off', 'BNBUSDT');

-- Заполнение таблицы historyData
INSERT INTO historyData (symbol, open_price, high_price, low_price, close_price, interval)
VALUES
  ('BTCUSDT', 58200.00, 58450.00, 58050.00, 58350.00, '1m'),
  ('BTCUSDT', 58150.00, 58300.00, 58000.00, 58275.00, '1m'),
  ('ETHUSDT', 3500.00, 3550.00, 3450.00, 3525.00, '5m'),
  ('ETHUSDT', 3545.00, 3575.00, 3500.00, 3560.00, '5m'),
  ('BNBUSDT', 600.00, 620.00, 590.00, 615.00, '15m'),
  ('BNBUSDT', 610.00, 625.00, 605.00, 620.00, '15m');

-- Заполнение таблицы users
INSERT INTO users (api_key, secret_key, strategy_id, username)
VALUES
  ('api_key_1', 'secret_key_1', 1, 'user_1'),
  ('api_key_2', 'secret_key_2', 3, 'user_2'),
  ('api_key_3', 'secret_key_3', 2, 'user_3');

-- Заполнение таблицы orders
INSERT INTO orders (type, strategy_id, history_id)
VALUES
  ('buy', 1, 1),
  ('sell', 2, 4),
  ('buy', 1, 2),
  ('sell', 3, 5);

-- Заполнение таблицы transactions
INSERT INTO transactions (order_id, user_id, quantity)
VALUES
  (1, 1, 0.01),
  (2, 2, 0.05),
  (3, 3, 0.02),
  (4, 1, 0.03),
  (4, 2, 0.01);
