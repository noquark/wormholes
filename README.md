<h1 align="center">wormholes</h1>
<p align='center'>
  <img alt="GitHub" src="https://img.shields.io/github/actions/workflow/status/noquark/wormholes/docker.yml?logo=github&style=flat-square" />
  <img alt="GitHub" src="https://img.shields.io/github/license/noquark/wormholes?logo=gnu&style=flat-square" />
</p>
<p align="center">
  <b>Lighting Fast Link Shortener</b><br />
</p>
<br />

## Getting Started

You can run wormholes from provided `compose.yml`.

```sh
git clone --depth=1 git@github.com:noquark/wormholes.git
cd wormholes
# Change PostgresQL and Redis configuration
docker compose up -d
```

### Redirection Endpoint

1. **GET** `:5000/:id`

### API Endpoints

1. **PUT** `:5000/api/`
2. **POST** `:5000/api/:id`
3. **GET** `:5000/api/:id`
4. **DELETE** `:5000/api/:id`
5. **GET** `:5000/api/:id`

## Configuration

### Customizing Ports

- `PORT` - Application port. Default value is `5000`.
- `GEN_PORT` - Generator port. Default value is `5001`

### Customizing database connections

Wormholes uses PostgreSQL and Redis. You can customize connection to these using environment variables as follows &mdash;

**PostgreSQL**

- `PG_URI` - This controls the URI for connecting to PostgreSQL. The default is `postgres://postgres:postgres@localhost:5432/postgres`.
- `PG_MAX_CONN` - This controls the max connections for PostgreSQL. The default is `5000`.

**Redis**

- `REDIS_URI` - THis controls the URI connecting to Redis and the default is `redis://:redis@localhost:6379/0`.

### Links Ingestion

Links are ingested in a batch to avoid excessive database connections. We can control it's behavior with following environment variables &mdash;

- `BATCH_SIZE` - This controls number of links ingested in a batch. The default value is `10000`.

### Customizing ID Generation

- `ID_SIZE` - This controls the size of generated IDs. The default value is `7`.
- `BLOOM_MAX` - This configures bloom-filters based on approx number of IDs to store. The default value is `1000000`.
- `BLOOM_ERROR` - This controls the rate of false positives in bloom filter and the default is `0.0000001`.
- `BUCKET_SIZE` - Inside generator, IDs to be used are stored in buckets. This controls the number of buckets to store IDs `8`.
- `BUCKET_CAP` - This controls the number of IDs to store in a single bucket which is `100000 ` by default.

## Contributing

Feel free to open an issue or pull request.
