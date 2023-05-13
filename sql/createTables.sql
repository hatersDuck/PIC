CREATE TABLE tradeStrategy (
  strategy_id SERIAL PRIMARY KEY,
  title CHAR(20),
  path_file TEXT NOT NULL DEFAULT '/strategy/default.py',
  status CHAR(6) NOT NULL DEFAULT 'off' CHECK (status IN ('active', 'pause', 'off')),
  symbol CHAR(10) NOT NULL DEFAULT 'BTCUSDT'
);

CREATE TABLE historyData (
  history_id SERIAL PRIMARY KEY,
  date TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
  symbol CHAR(10) NOT NULL DEFAULT 'BTCUSDT',
  price DECIMAL(18,8) NOT NULL
);

CREATE TABLE users (
  user_id BIGSERIAL PRIMARY KEY,
  api_key CHAR(64) DEFAULT 'empty',
  secret_key CHAR(64) DEFAULT 'empty',
  strategy_id SERIAL NOT NULL REFERENCES tradeStrategy(strategy_id) ON UPDATE CASCADE ON DELETE SET DEFAULT,
  status_trade CHAR(1) NOT NULL DEFAULT 'N',
  username CHAR(32) NOT NULL DEFAULT 'no',
  state_in_bot char(2) NOT NULL DEFAULT 'no',
  success BOOLEAN DEFAULT 'f'
);

CREATE TABLE orders (
  order_id SERIAL PRIMARY KEY,
  type CHAR(10) NOT NULL,
  date TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
  strategy_id SERIAL NOT NULL REFERENCES tradeStrategy(strategy_id) ON UPDATE CASCADE ON DELETE CASCADE,
  history_id SERIAL NOT NULL REFERENCES historyData(history_id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE transactions (
  transaction_id SERIAL PRIMARY KEY,
  order_id SERIAL NOT NULL REFERENCES orders(order_id) ON UPDATE CASCADE ON DELETE CASCADE,
  user_id BIGSERIAL NOT NULL REFERENCES users(user_id) ON UPDATE CASCADE,
  date TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
  quantity DECIMAL(18,8) NOT NULL
);