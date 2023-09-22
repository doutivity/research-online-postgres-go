# Efficiently store online with PostgreSQL and Go
- [Збереження стану онлайну користувача в Redis](https://dou.ua/forums/topic/35260/)
- [Hash, Set чи Sorted set. Який тип даних вибрати для збереження стану онлайну користувача в Redis?](https://dou.ua/forums/topic/44655/)

# Support Ukraine 🇺🇦
- Help Ukraine via [SaveLife fund](https://savelife.in.ua/en/donate-en/)
- Help Ukraine via [Dignitas fund](https://dignitas.fund/donate/)
- Help Ukraine via [National Bank of Ukraine](https://bank.gov.ua/en/news/all/natsionalniy-bank-vidkriv-spetsrahunok-dlya-zboru-koshtiv-na-potrebi-armiyi)
- More info on [war.ukraine.ua](https://war.ukraine.ua/) and [MFA of Ukraine](https://twitter.com/MFA_Ukraine)

# Testing
```bash
make env-up
make docker-go-version
make docker-pg-version
make migrate-up
make go-test
make go-bench
make env-down
```

# Schema
```sql
CREATE TABLE user_online
(
    user_id BIGINT PRIMARY KEY,
    online  TIMESTAMP NOT NULL
);
```

# Examples
```sql
TRUNCATE user_online;

INSERT INTO user_online (user_id, online)
VALUES (1, '2023-08-07 10:01:00'),
       (2, '2023-08-07 10:02:00'),
       (3, '2023-08-07 10:03:00'),
       (4, '2023-08-07 10:04:00'),
       (5, '2023-08-07 10:05:00'),
       (6, '2023-08-07 10:06:00'),
       (7, '2023-08-07 10:07:00'),
       (8, '2023-08-07 10:08:00'),
       (9, '2023-08-07 10:09:00'),
       (10, '2023-08-07 10:10:00'),
       (11, '2023-08-07 10:11:00'),
       (12, '2023-08-07 10:12:00')
ON CONFLICT (user_id) DO UPDATE
    SET online = excluded.online;

START TRANSACTION;
INSERT INTO user_online (user_id, online)
VALUES (1, '2023-08-07 11:01:00')
ON CONFLICT (user_id) DO UPDATE
    SET online = excluded.online;

INSERT INTO user_online (user_id, online)
VALUES (2, '2023-08-07 11:02:00')
ON CONFLICT (user_id) DO UPDATE
    SET online = excluded.online;
COMMIT;

START TRANSACTION;
UPDATE user_online
SET online = '2023-08-07 12:03:00'
WHERE user_id = 3;

UPDATE user_online
SET online = '2023-08-07 12:04:00'
WHERE user_id = 4;
COMMIT;

UPDATE user_online
SET online = CASE user_id
                 WHEN 5 THEN '2023-08-07 13:05:00'::TIMESTAMP
                 WHEN 6 THEN '2023-08-07 13:06:00'::TIMESTAMP
    END
WHERE user_id IN (5, 6);

UPDATE user_online AS to_t
SET online = from_t.online
FROM (
         VALUES (7, '2023-08-07 14:07:00'::TIMESTAMP),
                (8, '2023-08-07 14:08:00'::TIMESTAMP)
     ) AS from_t (user_id, online)
WHERE to_t.user_id = from_t.user_id;

-- version 1
SELECT *
FROM unnest(
             ARRAY[9, 10],
             ARRAY['2023-08-07 15:09:00'::TIMESTAMP, '2023-08-07 15:10:00'::TIMESTAMP]
         ) AS from_t (user_id, online);

-- version 2 supported https://github.com/sqlc-dev/sqlc/issues/958
SELECT unnest(ARRAY[9, 10])                                                              AS user_id,
       unnest(ARRAY['2023-08-07 15:09:00'::TIMESTAMP, '2023-08-07 15:10:00'::TIMESTAMP]) AS online;

-- version 1
UPDATE user_online AS to_t
SET online = from_t.online
FROM unnest(
             ARRAY[9, 10],
             ARRAY['2023-08-07 15:09:00'::TIMESTAMP, '2023-08-07 15:10:00'::TIMESTAMP]
         ) AS from_t (user_id, online)
WHERE to_t.user_id = from_t.user_id;

-- version 2
UPDATE user_online AS to_t
SET online = from_t.online
FROM (
         SELECT unnest(ARRAY[9, 10])                                                              AS user_id,
                unnest(ARRAY['2023-08-07 15:09:00'::TIMESTAMP, '2023-08-07 15:10:00'::TIMESTAMP]) AS online
     ) AS from_t
WHERE to_t.user_id = from_t.user_id;

-- version 1
INSERT INTO user_online (user_id, online)
VALUES (unnest(ARRAY[11, 12]),
        unnest(ARRAY['2023-08-07 16:11:00'::TIMESTAMP, '2023-08-07 16:12:00'::TIMESTAMP]))
ON CONFLICT (user_id) DO UPDATE
    SET online = excluded.online;

-- version 2
INSERT INTO user_online (user_id, online)
SELECT user_id, online
FROM unnest(
             ARRAY[11, 12],
             ARRAY['2023-08-07 16:11:00'::TIMESTAMP, '2023-08-07 16:12:00'::TIMESTAMP]
         ) AS from_t (user_id, online)
ON CONFLICT (user_id) DO UPDATE
    SET online = excluded.online;

-- version 3
INSERT INTO user_online (user_id, online)
SELECT user_id, online
FROM (
         SELECT unnest(ARRAY[11, 12])                                                             AS user_id,
                unnest(ARRAY['2023-08-07 16:11:00'::TIMESTAMP, '2023-08-07 16:12:00'::TIMESTAMP]) AS online
     ) AS from_t
ON CONFLICT (user_id) DO UPDATE
    SET online = excluded.online;

SELECT *
FROM user_online
ORDER BY user_id;
```

# Benchmark
```bash
make bench
```
| Name            | ns/op      | B/op    | allocs/op |
|-----------------|------------|---------|-----------|
| TxLoopUpsert    | 69_837_876 | 168_056 | 5_003     |
| TxLoopUpdate    | 63_966_207 | 168_056 | 5_003     |
| BatchExecUpsert | 19_463_064 | 503_235 | 5_030     |
| BatchExecUpdate | 18_686_485 | 495_235 | 5_030     |
| UnnestUpsert    | 7_997_338  | 234_930 | 2_027     |
| UnnestUpdate    | 7_950_833  | 234_930 | 2_027     |

```text
name             time/op
TxLoopUpsert     73.2ms ± 8%
TxLoopUpdate     71.1ms ±12%
BatchExecUpsert  20.3ms ±10%
BatchExecUpdate  20.2ms ±10%
UnnestUpsert     8.27ms ± 3%
UnnestUpdate     8.41ms ± 5%

name              alloc/op
TxLoopUpsert      168kB ± 0%
TxLoopUpdate      160kB ± 0%
BatchExecUpsert   503kB ± 0%
BatchExecUpdate   495kB ± 0%
UnnestUpsert      235kB ± 0%
UnnestUpdate      235kB ± 0%

name              allocs/op
TxLoopUpsert      5.00k ± 0%
TxLoopUpdate      5.00k ± 0%
BatchExecUpsert   5.03k ± 0%
BatchExecUpdate   5.03k ± 0%
UnnestUpsert      2.03k ± 0%
UnnestUpdate      2.03k ± 0%
```