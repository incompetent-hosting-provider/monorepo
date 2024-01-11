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

## Keycloak

### Setup
Make sure you have Docker installed.

## Access Keycloak Admin Interface
Once the containers are running, access the Keycloak Admin interface using a web browser. The default address is http://localhost:8080.

### User credentials

<details>
<summary>Admin</summary>
To log in to the admin console use the default admin credentials:

Username: admin\
Password: admin
</details>

<details>
<summary>Test users</summary>
When docker compose is started, a realm containing three test users is loaded:

- test-user-1
- test-user-2
- test-user-3

All these test users use a default password "Test123"
</details>


## Additional

## Ports

Keycloak: 8080   
DynamoDB: 8000

### Remove all data

```sh
rm -R docker
docker compose up -d --force-recreate
```


