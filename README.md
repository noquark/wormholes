<h1 align="center">wormholes</h1>
<p align='center'>
  <img alt="GitHub" src="https://img.shields.io/github/actions/workflow/status/wormholesdev/wormholes/docker.yml?logo=github&style=flat-square" />
  <img alt="Go Report Card" src="https://img.shields.io/badge/go%20report-A-green.svg?style=flat-square&logo=go&logoColor=white" />
  <img alt="GitHub" src="https://img.shields.io/github/license/wormholesdev/wormholes?logo=gnu&style=flat-square" />
  <img alt="Go Version" src="https://img.shields.io/github/go-mod/go-version/wormholesdev/wormholes?logo=go&logoColor=white&style=flat-square" />
</p>
<br />

## Features

- **Fast**. Create millions of links within minute.
- **Easy to use**. With unified mode, run everything except databases in a single image.
- **Scalable**. With distributed mode, run multiple instances of creator, and redirector.
- **Analytics**. Detailed Events are stored in TimescaleDB. Dashboard is WIP.

## Getting Started

### Preparing

Default database configs are inside `deploy/conf`. Verify those. Included postgres config is tuned for my system for 5000 connections. Generate your own with [pgtune](https://pgtune.leopard.in.ua/#/). You may also want to mount volumes for database which you can do inside default compose files.

### Running in Unified Mode

Run wormholes with docker-compose

```sh
cd deploy
docker compose -f compose/unified.yml up -d
```

Following are the API endpoints in unified mode.

1. **PUT** `:5000/v1/links`
2. **POST** `:5000/v1/links/:id`
3. **GET** `:5000/v1/links/:id`
4. **DELETE** `:5000/v1/links/:id`
5. **GET** `:5000/l/:id`

### Running in Distributed Mode

Run wormholes with docker-compose

```sh
cd deploy
docker compose -f compose/distributed.yml up -d
```

Following are the API endpoints in distributed mode.

1. **PUT** `:5002`
2. **POST** `:5002/:id`
3. **GET** `:5002/:id`
4. **DELETE** `:5002/:id`
5. **GET** `:5000/:id`

## Additional Help

For detailed documentation visit [Docs](https://wormholes.dev/docs)

## Contributing

Feel free to open an issue or pull request.
