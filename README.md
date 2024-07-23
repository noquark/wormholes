<h1 align="center">wormholes</h1>
<p align='center'>
  <img alt="GitHub" src="https://img.shields.io/github/actions/workflow/status/noquark/wormholes/docker.yml?logo=github&style=flat-square" />
  <img alt="GitHub" src="https://img.shields.io/github/license/noquark/wormholes?logo=gnu&style=flat-square" />
</p>
<p align="center">
  <b>Lighting Fast and Distributed Link Shortener</b><br />
  <sub>With powerful event tracking and analytics</sub>
</p>
<br />

## Getting Started

You can run wormholes from provided `docker-compose.yml` under `./deploy` directory.

```sh
cd deploy
docker compose up -d
```

### API Endpoints

1. **PUT** `:5002`
2. **POST** `:5002/:id`
3. **GET** `:5002/:id`
4. **DELETE** `:5002/:id`
5. **GET** `:5000/:id`

## Confiuration

### Customizing ports

Wormholes comprises three services: **Director**, **Generator**, and **Creator**. In distributed mode, each service runs on a specific port. You can customize these ports using `PORT` environment variable.

**Default Ports**

- **Director** - `5000`
- **Generator** - `5001`
- **Creator** - `5002`

### Customizing database connections

Wormholes uses TimescaleDB, PostgreSQL and Redis. You can customize connection to these databses using environment variables as follows &mdash;

**TimescaleDB**

- `TS_URI`- This controls the URI for connecting to TimescaleDB. The default is `postgres://postgres:postgres@localhost:5433/postgres`.

**PostgreSQL**

- `PG_URI` - This controls the URI for connecting to PostgreSQL. The default is `postgres://postgres:postgres@localhost:5432/postgres`.
- `PG_MAX_CONN` - This controls the max connections for PostgreSQL. The default is `5000`.

**Redis**

- `REDIS_URI` - THis controls the URI connecting to Redis and the default is `redis://:redis@localhost:6379/0`.

### Links Ingestion

Links generated in **Creator** are ingested in a batch to avoid excessive database connections. We can control it's behaviour with following environment variables &mdash;

- `BATCH_SIZE` - This controls number of links ingested in a batch. The default value is `10000`.

### Events ingestion

When a link in clicked, an event is generated in **Director**. All such events are ingested to a time-series database ( **TimescaleDB** by default ). We can control it with following envrionment variables &mdash;

- `BATCH_SIZE` - This controls number of events ingested in a batch. The default value is `10000`.
- `STREAMS` - I this controls the number of streams that process events. The default value is `8`.

### Customizing Generator

- `ID_SIZE` - This controls the size of generated IDs. The default value is `7`.
- `BLOOM_MAX` - This configures bloomfilters based on approx number of IDs to store. The default value is `1000000`.
- `BLOOM_ERROR` - This controls the rate of false positives in bloom filter and the default is `0.0000001`.
- `BUCKET_SIZE` - Inside generator, IDs to be used are stored in buckets. This controls the number of buckets to store IDs `8`.
- `BUCKET_CAP` - This controls the number of IDs to store in a single bucket which is `100000 ` by default.

### Connecting Creator to Generator

The **Creator** connects with **Generator** over grpc and the url to generator instance can be controlled with `GEN_ADDR`. The default is `localhost:5001`.

## Contributing

Feel free to open an issue or pull request.
