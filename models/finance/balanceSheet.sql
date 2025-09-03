SELECT
  (account) AS a,
  name,
  SUM(amount * direction * normal) AS balance
FROM
  transactions
  LEFT JOIN account ON a = account.number
GROUP BY
  name
ORDER BY
  a,
  name;