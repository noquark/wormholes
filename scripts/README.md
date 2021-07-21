# scripts

This directory contains various scripts for development and load testing.

## Managing database

```sh
./scripts/ctl create|destroy pg|ts
```

## Accessing shells

```sh
./script/shell pg|ts
```

## Load testing

This requires `wrk` to be installed.

Test links creation with

```sh
./scripts/load create
```

Now, Pick one link from links table. Access shell with

```sh
./scripts/shell pg
```

and get a link with

```sql
SELECT id from wh_links LIMIT 1;
```

Using link from previous step, test link redirects with

```sh
./scripts/load get <link>
```

Finally, test API access with

```sh
./scripts/load api <link>
```
