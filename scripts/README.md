# scripts

This directory contains various scripts for development and load testing.

## Managing database

To create databases

```sh
./scripts/create
```

And, to destroy them

```sh
./scripts/destroy
```

## Accessing shells

To access postgres shell, run

```sh
./script/db
```

and for timescale shell, run

```sh
./script/tsdb
```

## Load testing

This requires `wrk` to be installed.

Test links creation with

```sh
./scripts/load create
```

Now, Pick one link from links table. Access shell with

```sh
./scripts/db
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
