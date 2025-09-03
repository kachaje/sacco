INSERT INTO
  account (name, number, normal)
VALUES
  ("Assets", 100, 1),
  ("Cash", 110, 1),
  ("Merchandise", 120, 1),
  ("Liabilities", 200, -1),
  ("Deferred Revenue", 210, -1),
  ("Revenues", 300, -1),
  ("Expenses", 400, 1),
  ("Cost of Goods Sold", 410, 1),
  ("Equity", 500, -1),
  ("Capital", 510, -1);

INSERT INTO
  transactions (id, date, amount, account, direction)
VALUES
  (0, "2022-01-01", 500.0, 110, 1),
  (0, "2022-01-01", 500.0, 510, -1),
  (1, "2022-01-01", 100.0, 120, 1),
  (1, "2022-01-01", 100.0, 110, -1),
  (2, "2022-02-01", 15.0, 110, 1),
  (2, "2022-02-01", 15.0, 210, -1),
  (3, "2022-02-05", 15.0, 210, 1),
  (3, "2022-02-05", 15.0, 300, -1),
  (4, "2022-02-05", 3.0, 410, 1),
  (4, "2022-02-05", 3.0, 120, -1);