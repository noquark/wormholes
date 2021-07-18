<h1 align='center'>wormholes</h1>
<p align="center">
  <img alt="docker status" src="https://github.com/mohitsinghs/wormholes/actions/workflows/docker.yml/badge.svg">
</p>
<p align="center">
  <b>A self-hosted link shortener</b><br/>
  <sub>powered by Fiber and TimescaleDB</sub>
</p>
<br />

## Features

- Fast, fail safe and collisions free link creation.
- Highly configurable
- Easy to use docker image
- Analytics dashboard (WIP)
- API keys (WIP)

## Running

### Setting up docker

If you don't have docker installed already, you can install it [from here](https://docs.docker.com/get-docker/).

Another alternative for linux is

```sh
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh
```

Once docker is installed, create a common network to be used with timescale and wormholes.

```sh
docker network create wormholes
```

### Setting up TimescaleDB

Run TimescaleDB with our newly created network -

```sh
docker run -d \
  --name=timescale
  --network=wormholes
  -e POSTGRES_PASSWORD=postgres
  timescale/timescaledb:latest-pg13
```

### Setting up Wormholes

Wormholes can be customized with environment variables or config.
To run with default configuration, start with -

```sh
docker run -d \
  --name=wormholes
  --network=wormholes
  -p 3000:3000
  ghcr.io/mohitsinghs/wormholes
```

This exposes wormholes web interface on port 3000.
