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
    user_id BIGINT    NOT NULL PRIMARY KEY,
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

# Benchmark (Postgres 15.3) (Go 1.20) (PC) Intel(R) Core(TM) i7-8750H CPU @ 2.20GHz
```bash
make go-bench
```

| Name            | ns/op      | B/op    | allocs/op |
|-----------------|------------|---------|-----------|
| TxLoopUpdate    | 63_966_207 | 168_056 | 5_003     |
| TxLoopUpsert    | 69_837_876 | 168_056 | 5_003     |
| UnnestUpdate    | 7_950_833  | 234_930 | 2_027     |
| UnnestUpsert    | 7_997_338  | 234_930 | 2_027     |
| BatchExecUpdate | 18_686_485 | 495_235 | 5_030     |
| BatchExecUpsert | 19_463_064 | 503_235 | 5_030     |

| name            | time/op     |
|-----------------|-------------|
| TxLoopUpdate    | 71.1ms ±12% |
| TxLoopUpsert    | 73.2ms ± 8% |
| UnnestUpdate    | 8.41ms ± 5% |
| UnnestUpsert    | 8.27ms ± 3% |
| BatchExecUpdate | 20.2ms ±10% |
| BatchExecUpsert | 20.3ms ±10% |

| name            | B/op       |
|-----------------|------------|
| TxLoopUpdate    | 160kB ± 0% |
| TxLoopUpsert    | 168kB ± 0% |
| UnnestUpdate    | 235kB ± 0% |
| UnnestUpsert    | 235kB ± 0% |
| BatchExecUpdate | 495kB ± 0% |
| BatchExecUpsert | 503kB ± 0% |

| name            | allocs/op  |
|-----------------|------------|
| TxLoopUpdate    | 5.00k ± 0% |
| TxLoopUpsert    | 5.00k ± 0% |
| UnnestUpdate    | 2.03k ± 0% |
| UnnestUpsert    | 2.03k ± 0% |
| BatchExecUpdate | 5.03k ± 0% |
| BatchExecUpsert | 5.03k ± 0% |

# Benchmark (Postgres 16.0) (Go 1.21) (PC) Intel(R) Core(TM) i7-12700H
```bash
make go-bench
```

| Name            | ns/op      | B/op    | allocs/op |
|-----------------|------------|---------|-----------|
| TxLoopUpdate    | 19_786_396 | 160_135 | 5_005     |
| TxLoopUpsert    | 20_168_659 | 168_135 | 5_005     |
| UnnestUpdate    | 3_935_782  | 234_985 | 2_028     |
| UnnestUpsert    | 3_902_771  | 234_985 | 2_028     |
| BatchExecUpdate | 6_984_122  | 495_315 | 5_032     |
| BatchExecUpsert | 6_630_488  | 503_316 | 5_032     |

| name            | time/op       |
|-----------------|---------------|
| TxLoopUpdate    | 20.45ms ±  2% |
| TxLoopUpsert    | 26.59ms ± 24% |
| UnnestUpdate    | 3.997ms ±  1% |
| UnnestUpsert    | 3.998ms ±  2% |
| BatchExecUpdate | 7.044ms ±  1% |
| BatchExecUpsert | 7.004ms ±  8% |

| name            | B/op         |
|-----------------|--------------|
| TxLoopUpdate    | 156.4kB ± 0% |
| TxLoopUpsert    | 164.2kB ± 0% |
| UnnestUpdate    | 229.5kB ± 0% |
| UnnestUpsert    | 229.5kB ± 0% |
| BatchExecUpdate | 483.7kB ± 0% |
| BatchExecUpsert | 491.5kB ± 0% |

| name            | allocs/op   |
|-----------------|-------------|
| TxLoopUpdate    | 5.005k ± 0% |
| TxLoopUpsert    | 5.005k ± 0% |
| UnnestUpdate    | 2.028k ± 0% |
| UnnestUpsert    | 2.028k ± 0% |
| BatchExecUpdate | 5.032k ± 0% |
| BatchExecUpsert | 5.032k ± 0% |

# Benchmark (Postgres 16.0) (Go 1.21) ([vultr.com](https://www.vultr.com/?ref=8741375) VPS 131072.00 MB 8 cores / 16 threads @ 3.2 GHz) ($350/month) Intel(R) Xeon(R) E-2388G CPU @ 3.20GHz
```bash
make go-bench
```

| Name            | ns/op      | B/op    | allocs/op |
|-----------------|------------|---------|-----------|
| TxLoopUpdate    | 46_126_147 | 160_135 | 5_005     |
| TxLoopUpsert    | 45_719_610 | 168_135 | 5_005     |
| UnnestUpdate    | 5_123_888  | 234_985 | 2_028     |
| UnnestUpsert    | 5_127_720  | 234_985 | 2_028     |
| BatchExecUpdate | 11_179_808 | 495_315 | 5_032     |
| BatchExecUpsert | 11_252_240 | 503_316 | 5_032     |

| name            | time/op      |
|-----------------|--------------|
| TxLoopUpdate    | 46.83ms ± 1% |
| TxLoopUpsert    | 47.21ms ± 2% |
| UnnestUpdate    | 5.196ms ± 1% |
| UnnestUpsert    | 5.230ms ± 1% |
| BatchExecUpdate | 11.21ms ± 0% |
| BatchExecUpsert | 11.36ms ± 1% |

| name            | B/op         |
|-----------------|--------------|
| TxLoopUpdate    | 156.4kB ± 0% |
| TxLoopUpsert    | 164.2kB ± 0% |
| UnnestUpdate    | 229.5kB ± 0% |
| UnnestUpsert    | 229.5kB ± 0% |
| BatchExecUpdate | 483.7kB ± 0% |
| BatchExecUpsert | 491.5kB ± 0% |

| name            | allocs/op   |
|-----------------|-------------|
| TxLoopUpdate    | 5.005k ± 0% |
| TxLoopUpsert    | 5.005k ± 0% |
| UnnestUpdate    | 2.028k ± 0% |
| UnnestUpsert    | 2.028k ± 0% |
| BatchExecUpdate | 5.032k ± 0% |
| BatchExecUpsert | 5.032k ± 0% |
