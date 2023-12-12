# IHC 

## Quick start

**Note**: `docker compose` is an alias for `docker-compose`. Depending on your docker version and/or OS you might have to use the ladder.

To start the local dev setup:

```sh
docker compose up -d
```

This starts all components except for the cli.
When developing one of the services the corresponding components container can be stopped manually or via the `--scale` option.
```sh
docker compose up -d --scale=<service name>=0
```

## Additional

## Ports

Keycloak: 8080   
DynamoDB: 8000

### Remove all data

```sh
rm -R docker
docker compose up -d --force-recreate
```


