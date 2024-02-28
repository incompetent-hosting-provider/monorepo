# Endpoints

## Status

| Endpoint                       | CLI                | Backend            |
|--------------------------------|--------------------|--------------------|
| GET   /user/                   | :white_check_mark: | :white_check_mark: |
| GET /user/balance              | :white_check_mark: | :white_check_mark: |
| POST /user/balance             | :white_check_mark: | :white_check_mark: |
| GET /instances                 |                    | :white_check_mark: |
| GET /instances/<instanceid>    |                    | :white_check_mark: |
| POST /instances/custom         |                    | :white_check_mark: |
| POST /instances/preset         |                    | :white_check_mark: |
| DELETE /instances/<instanceid> |                    | :white_check_mark: |
| GET /service/available-presets |                    | :white_check_mark: |

Temporary definition - subject to change.

**Note**: JWT token in Authorization Header in format

"Authorization": "Bearer <jwt>"

Error responses:

```json
{
  "error": "This went wrong because abc"
}
```

Instance object:

```json
{
    "type": "custom"|"preset",
    "name": "name",
    "id": "kjahsjdjaksd",
    "container image": {
        "name": "postgres",
        "version": "12.3"
    },
    "status": "RUNNING" | "PENDING" | "TERMINATED"
}
```

Instance object detailed:

```json
{
    ...instance object,
    "started_at": "timestamp",
    "created_at": "timestamp",
    "open_ports": [123, 456],
    "description": "This is my roblox server guys!"
}
```

## Account

GET `/user/`

Response:

```
{
    "email": "user@email.com",
    "balance": 1000,
}
```

GET `/user/balance`

Response:

```json
{
  "balance": 1000
}
```

---

POST `/user/balance`

Body:

```json
{
  "amount": 1000
}
```

Response:

```json
{
  "balance": 1000
}
```

## Container/Service

GET `/instances`

Response:

```json
{
    "instances":[
        instance object
    ]
}
```

---

GET `/instances/<id>`

Response:

```json
detailed instance object
```

---

POST `/instances/preset`

Body:

```json
{
  "preset": 1,
  "name": "my user defined name",
  "description?": "description"
}
```

=> 202
Response

```json
{
  "id": "dasdhjsk",
  "env_vars": {
    "..": ".."
  }
}
```

---

POST `/instances/custom`

Body:

```json
{
    "name": "hello",
    "description": "",
    "image":{
        "name": "asjkdas",
        "version": "3.21.1"
    },
    "env_vars": {
        "<var name>": "<var value>",
        ...
    },
    "ports": [123123,132312]
}
```

=> 202
Response

```json
{
  "id": "dasdhjsk"
}
```

---

DELETE `/instances/<id>`

=> 202

## Service

GET `/service/available-presets`

Response:

```json
{
  "presets": [
    {
      "name": "hello",
      "id": 1,
      "description": "Hello I am a mysql instance"
    }
  ]
}
```
