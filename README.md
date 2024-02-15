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
- **Distributed Mode** : Run multiple creator and redirector instances for scale
- **Powerful Analytics** : Detailed event tracking and analytics

## Getting Started

> [!NOTE]
> The provided docker compose uses timescale for storing both events and links, but you can always choose to use postgres for storing links and timescale for events. See [config](https://noquark.com/docs/wormholes/configuration#customizing-database-connections)

You can run wormholes from provided `docker-compose.yml` under `./deploy` directory.

```sh
docker compose up -d
```

## API Endpoints

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
