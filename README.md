# IHP

## Quick start

**Note**: `docker compose` is an alias for `docker-compose`. Depending on your docker version and/or OS you might have to use the ladder.

To start the local dev setup:
**Note**: `arm` architecture is currently not supported, due to incompatible dependencies in terraform providers.

```sh
docker compose up -d
```

This starts all components except for the cli.
When developing one of the services the corresponding components container can be stopped manually or via the `--scale` option.

```sh
docker compose up -d --scale=<service name>=0
```

| Service  | Endpoint                                |
| -------- | --------------------------------------- |
| Keycloak | [endpoint admin](http://localhost:8080) |
| Grafana  | [endpoint](http://localhost:3000)       |

## Credentials

### Keycloak

Admin:

```
admin:admin
```

<details>
<summary>Sample Users</summary>
When docker compose is started, a realm containing three test users is loaded:

- test-user-1
- test-user-2
- test-user-3

All these test users use a default password "Test123"

</details>

### Grafana

Admin:

```
admin:admin
```

## Additional

Note that dynamodb persists its db within the `docker/dynamodb` directory. This needs to be deleted to fully reset the application.

As of now the creation of custom containers as well as the deletion of running containers is not supported by the terraform service.
