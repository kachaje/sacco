CREATE TABLE
  IF NOT EXISTS account (
    name TEXT UNIQUE NOT NULL,
    number INTEGER UNIQUE NOT NULL,
    normal INTEGER NOT NULL
  );

CREATE TABLE
  IF NOT EXISTS transactions (
    id INTEGER,
    date TEXT NOT NULL,
    amount REAL,
    account INTEGER,
    direction INTEGER
  );