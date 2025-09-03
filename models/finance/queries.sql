SELECT
  *
FROM
  transactions
  LEFT JOIN account ON transactions.account = account.number;

SELECT
  SUM(
    CASE
      WHEN direction = 1 THEN amount
    END
  ) AS DR,
  SUM(
    CASE
      WHEN direction = -1 THEN amount
    END
  ) AS CR
FROM
  transactions;

SELECT
  id,
  SUM(direction * amount) AS s
FROM
  transactions
GROUP BY
  id
HAVING
  s != 0;