<h1 align="center">wormholes</h1>
<p align='center'>
  <img alt="GitHub" src="https://img.shields.io/github/actions/workflow/status/noquark/wormholes/docker.yml?logo=github&style=flat-square" />
  <img alt="Go Report Card" src="https://img.shields.io/badge/go%20report-A-green.svg?style=flat-square&logo=go&logoColor=white" />
  <img alt="GitHub" src="https://img.shields.io/github/license/noquark/wormholes?logo=gnu&style=flat-square" />
  <img alt="Go Version" src="https://img.shields.io/github/go-mod/go-version/noquark/wormholes?logo=go&logoColor=white&style=flat-square" />
</p>
<br />

## Features

- **Lightning Fast** : Create millions of short links in minutes
- **Unified Mode** : Run everything except databases in a single image
- **Distributed Mode** : Run multiple creators and redirectors for scale
- **Powerful Analytics** : Detailed event tracking and analytics

## Getting Started

### Preparing

To get started with Wormholes, follow these steps:

- Verify the default database configurations located in `deploy/conf`. Make any necessary adjustments.
- Consider generating your custom Postgres configuration with [pgtune](https://pgtune.leopard.in.ua/#/).
- If needed, configure volumes for the database, which can be done within the default compose files.

### Running in Unified Mode

To run Wormholes in Unified Mode using Docker Compose:

```sh
cd deploy
docker compose -f compose/unified.yml up -d
```

### Running in Distributed Mode

Run wormholes with docker-compose

```sh
cd deploy
docker compose -f compose/distributed.yml up -d
```

## API Endpoints

### Unified Mode

1. **PUT** `:5000/v1/links`
2. **POST** `:5000/v1/links/:id`
3. **GET** `:5000/v1/links/:id`
4. **DELETE** `:5000/v1/links/:id`
5. **GET** `:5000/l/:id`

### Distributed Mode

1. **PUT** `:5002`
2. **POST** `:5002/:id`
3. **GET** `:5002/:id`
4. **DELETE** `:5002/:id`
5. **GET** `:5000/:id`

## Detailed Docs

For detailed documentation visit [Docs](https://noquark.com/docs/wormholes)

## Contributing

Feel free to open an issue or pull request.
