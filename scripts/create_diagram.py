import argparse
import psycopg2
import matplotlib.pyplot as plt

parser = argparse.ArgumentParser()
parser.add_argument("id", help="user id for profits graph")
args = parser.parse_args()

user_id = args.id

conn = psycopg2.connect(database="pic", user="danila", password="Ckj;ysq", host="localhost")
cur = conn.cursor()

query = f"""
SELECT ts.title, SUM(t.quantity) AS profit
FROM tradeStrategy ts 
JOIN orders o ON ts.strategy_id = o.strategy_id 
JOIN transactions t ON o.order_id = t.order_id 
WHERE t.quantity > 0 AND t.user_id = {user_id}
GROUP BY ts.title;
"""
cur.execute(query)

strategy_titles = []
profits = []
for row in cur:
    strategy_titles.append(row[0])
    profits.append(row[1])

cur.close()
conn.close()

plt.figure(figsize=(10, 5))
plt.bar(strategy_titles, profits)
plt.title(f"Profits by Strategy for User {user_id}")
plt.xlabel("Strategy Title")
plt.ylabel("Profit")

import datetime
import os

if not os.path.exists(f"img"):
    os.makedirs(f"img")
if not os.path.exists(f"img/{user_id}"):
    os.makedirs(f"img/{user_id}")

now = datetime.datetime.now()
filename = f"img/{user_id}/{now.strftime('%y_%m_%d_%H_%M')}.png"
plt.savefig(filename)
print(filename, end="")