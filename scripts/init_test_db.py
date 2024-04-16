import sqlite3 as sql

conn = sql.connect("box.db")
cur = conn.cursor()

res = cur.execute("CREATE TABLE box(address, rawValue)")

conn.commit()

