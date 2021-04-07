<h1 align='center'>wormholes</h1>
<p align="center">
  <img alt="docker status" src="https://github.com/mohitsinghs/wormholes/actions/workflows/docker.yml/badge.svg">
</p>
<p align="center">
  <b>A self-hosted link shortner</b><br/>
</p>
<br />

### Features
 
- **Reliable** - Creating shortlinks is fast, fail safe and collisions free.
- **Extensible** - Redis and PostgresSQL are supported, more to be added.
- **Configurable** - Configurable via yaml, cli flags and environment variables.
- **Ready to use** - Just pull the docker image and run with desired database and config.

### APIs

- **PUT** `/api/v1/links` - to create links
- **POST** `/api/v1/links/<id>` - to update link
- **GET** `/api/v1/links/<id>` - to get link data
- **DELETE** `/api/v1/links/<id>` - to delete link

### Roadmap

- [ ] Proper documentation
- [ ] Add more databases (SQLite, MongoDB, MySQL)
- [ ] Root redirects
- [ ] Link templates
- [ ] Bulk link generation
- [ ] QR code generation

### Contributing

Pull request are more than welcome for anything on roadmap.
