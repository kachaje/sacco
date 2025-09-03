-- Get Accounting Equation sides
SELECT
  '1. Accounting Equation';

SELECT
  '------------------------------------------------------------';

SELECT
  GROUP_CONCAT (name, ' + ') AS expression
FROM
  account
GROUP BY
  normal;

SELECT
  '';

SELECT
  '';

-- Get high level accounts
SELECT
  '2. High Level Accounts';

SELECT
  '------------------------------------------------------------';

SELECT
  GROUP_CONCAT (name, ' + ') AS expression
FROM
  account
WHERE
  MOD(number, 100) = 0
GROUP BY
  normal;

SELECT
  '';

SELECT
  '';

-- Full accounting equation
SELECT
  '3. Full accounting equation';

SELECT
  '------------------------------------------------------------';

SELECT
  MAX(left_side) || ' = ' || MAX(right_side) AS equation
FROM
  (
    SELECT
      GROUP_CONCAT (
        CASE
          WHEN normal = 1 THEN name
        end,
        ' + '
      ) AS left_side,
      GROUP_CONCAT (
        CASE
          WHEN normal = -1 THEN name
        end,
        ' + '
      ) AS right_side
    FROM
      account
    WHERE
      MOD(number, 100) = 0
    GROUP BY
      normal
  );

SELECT
  '';

SELECT
  '';