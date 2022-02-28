<h1 align="center">wormholes</h1>
<p align='center'>
  <img alt="GitHub Workflow Status" src="https://img.shields.io/github/workflow/status/mohitsinghs/wormholes/docker?logo=github&style=for-the-badge" />
  <img alt="Go Report Card" src="https://img.shields.io/badge/go%20report-A-green.svg?style=for-the-badge&logo=go&logoColor=white" />
  <img alt="GitHub" src="https://img.shields.io/github/license/mohitsinghs/wormholes?logo=gnu&style=for-the-badge" />
  <img alt="Go Version" src="https://img.shields.io/github/go-mod/go-version/mohitsinghs/wormholes?logo=go&logoColor=white&style=for-the-badge" />
</p>
<br />

## Features

- Fail Safe
- Scalable and Distributed
- Highly Customizable
- High Performance
- API Access

## Why ?

There are a lot of link-shorteners out there, open-source and commercial. they provide easy access to resources in diverse fields. It's an easy problem to solve until it's not.

When I witnessed how a quick in-house implementation can cost a business and will refuse to scale, I was curious on how to scale these reliably. When a quick lookup didn't help me find any solution, I decided to write my own. I've detailed about this [here](https://mohitsingh.in/building-a-link-shortner) and [here](https://mohitsingh.in/a-distributed-link-shortner).

## Getting started

Although, a better approach will be to run wormholes on docker or kubernetes, For now, let's test if this works by running databases and cache on docker and wormholes services on our machine &mdash;

### 1. Running Postgres

This particular config is tuned for my system. Generate your own with [pgtune](https://pgtune.leopard.in.ua/#/).

**postgres/postgres.conf**

```properties
# tuned for load testing on my local system
max_connections = 5000
shared_buffers = 4GB
effective_cache_size = 12GB
maintenance_work_mem = 1GB
checkpoint_completion_target = 0.9
wal_buffers = 16MB
default_statistics_target = 100
random_page_cost = 1.1
effective_io_concurrency = 200
work_mem = 419kB
min_wal_size = 1GB
max_wal_size = 4GB
max_worker_processes = 4
max_parallel_workers_per_gather = 2
max_parallel_workers = 4
max_parallel_maintenance_workers = 2
listen_addresses = '*'
```

Now run a postgres instance with docker.

```bash
docker run -d \
  --name=postgres \
  -p 127.0.0.1:5432:5432 \
  -e POSTGRES_PASSWORD=postgres \
  -v "$PWD/postgres/postgres.conf":/etc/postgresql/postgresql.conf \
  postgres:alpine \
  -c config_file=/etc/postgresql/postgresql.conf
```

### 2. Running Timescale

Run timescale with docker. We will use this to store our analytical data.

```bash
docker run -d \
  --name=timescale \
  -p 127.0.0.1:5433:5432 \
  -e POSTGRES_PASSWORD=postgres \
  timescale/timescaledb:latest-pg13
```

### 3. Running Redis

We will be using redis for lru cache

**redis/redis.conf**

```properties
# set max memory
maxmemory 1gb
# evict least recently used keys
maxmemory-policy allkeys-lru
# require a password on connect
requirepass redis
```

Now let's run a redis docker instance for caching purpose.

```bash
docker run -d \
  --name=redis \
  -p 127.0.0.1:6379:6379 \
  redis:6-alpine \
  -v "$PWD/redis":/usr/local/etc/redis \
  redis-server /usr/local/etc/redis/redis.conf
```

### 4. Preparing wrk

Install [wrk](https://github.com/wg/wrk) in your system. We will be using this to test our setup.

Create a script for put requests.

**scripts/put.lua**

```lua
wrk.method = "PUT"
wrk.headers["Content-Type"] = "application/json"
wrk.body = '{ "tag": "google", "target": "https://google.com" }'
```

### 5. Creating required tables

Create tables in our postgres instance.

```sql
-- links
create table if not exists wh_links (
  id text primary key,
  tag text,
  target text,
  created_at timestamptz not null default now()
);
```

and our timescale instance

```sql
-- clicks
create table if not exists wh_clicks (
  time timestamptz not null,
  link text not null,
  tag text not null,
  cookie text not null,
  is_mobile boolean,
  is_bot boolean,
  browser text,
  browser_version text,
  os text,
  os_version text,
  platform text,
  ip text
);

-- create hypertable
select
  create_hypertable ('wh_clicks', 'time', if_not_exists => true);
```

### 6. Running

Finally run generator followed by director and creator

```sh
git clone https://github.com/mohitsinghs/wormholes
cd wormholes

# run services
go run . -as generator
go run . -as creator
go run . -as director
```

### 7. Load Testing

**Creating Links**

Let's create links using our wrk script we created earlier.

```sh
wrk -t8 -d10s -c100 -s "./scripts/put.lua" http://localhost:5002/api/v1/links
```

**Getting Links**

Now, let's load test our links API by retrieving a link.

```sh
 wrk -t8 -d10s -c100 http://localhost:5001/api/v1/links/<shortID>
```

Retrieve a shortID from our postgres instance. We will use this for further tests.

**Redirecting**

Finally, let's test our link redirection.

```sh
wrk -t8 -d10s -c100 -s "./load/get.lua" http://localhost:5000/<shortID>
```

### Stats

On my machine, results are something like this. The running duration for each test in below table was 1 minute.

| Task         | Cache Hit | Performance (Avg) | Total Ops | Latency (Avg) |
| ------------ | --------- | ----------------- | --------- | ------------- |
| Create Links |           | 121K reqs/sec     | 7.26M     | 47.39ms       |
| Get Link     | Redis     | 60K reqs/sec      | 3.60M     | 1.68ms        |
| Get Link     | None      | 25K reqs/sec      | 1.5M      | 4.30ms        |
| Redirect     | Memory    | 398K reqs/sec     | 23M       | 0.94ms        |
| Redirect     | Redis     | 39K reqs/sec      | 2.3M      | 23.33ms       |
| Redirect     | None      | 21K reqs/sec      | 1.3M      | 90.22ms       |

## Environment Variables

Wormholes are highly customizable with environment variables. The default values are there to make them work on my local setup, for production these needs modifications. Following are the environment variables and there usage &mdash;

### Director

| Name          | Purpose                             |                                                Default |
| ------------- | ----------------------------------- | -----------------------------------------------------: |
| `PORT`        | The port to run                     |                                                 `5000` |
| `BATCH_SIZE`  | Size of batch when ingesting events |                                                `10000` |
| `STREAM`      | Number of streams to ingest events  |                                                    `8` |
| `PG_URI`      | URI for connecting to PostgreSQL    | `postgres://postgres:postgres@localhost:5432/postgres` |
| `PG_MAX_CONN` | Maximum connections for PostgreSQL  |                                                 `5000` |
| `REDIS_URI`   | URI for connecting to Redis         |                      `redis://:redis@localhost:6379/0` |
| `TS_URI`      | URI for connecting to TimescaleDB   | `postgres://postgres:postgres@localhost:5433/postgres` |

### Generator

| Name              | Purpose                                   |                                                Default |
| ----------------- | ----------------------------------------- | -----------------------------------------------------: |
| `PORT`            | The port to run                           |                                                 `5001` |
| `ID_SIZE`         | Size of generated IDs                     |                                                    `7` |
| `MAX_LIMIT`       | Limit of IDs to be stored                 |                                              `1000000` |
| `ERROR_RATE`      | Rate of false positives in bloom filter   |                                            `0.0000001` |
| `BUCKET_SIZE`     | Number of buckets to store IDs            |                                                   `8 ` |
| `BUCKET_CAPACITY` | Number of IDs to store in a single bucket |                                              `100000 ` |
| `PG_URI`          | URI for connecting to PostgreSQL          | `postgres://postgres:postgres@localhost:5432/postgres` |
| `PG_MAX_CONN`     | Maximum connections for PostgreSQL        |                                                 `5000` |

### Creator

| Name          | Purpose                            |                                                Default |
| ------------- | ---------------------------------- | -----------------------------------------------------: |
| `PORT`        | The port to run                    |                                                 `5002` |
| `GEN_ADDR`    | Address of generator instance      |                                       `localhost:5001` |
| `BATCH_SIZE`  | Size of batch when ingesting links |                                                `10000` |
| `PG_URI`      | URI for connecting to PostgreSQL   | `postgres://postgres:postgres@localhost:5432/postgres` |
| `PG_MAX_CONN` | Maximum connections for PostgreSQL |                                                 `5000` |
| `REDIS_URI`   | URI for connecting to Redis        |                      `redis://:redis@localhost:6379/0` |

## Behind the name

Wormholes are tunnels made by earthworms and a link between two disparate points in spacetime. Short links are essentially a links to digital thing. This project is inspired by both and hence named after them.

## Having problems ?

Feel free to open an issue or feature request. I'll try to help you as long as you really care about open source and are being reasonable.
