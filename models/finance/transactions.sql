SELECT
  id,
  date,
  name,
  CASE
    WHEN direction = 1 THEN amount
  END AS DR,
  CASE
    WHEN direction = -1 THEN amount
  END AS CR
FROM
  transactions
  LEFT JOIN account ON account = account.number
ORDER BY
  id,
  date,
  CR,
  DR;