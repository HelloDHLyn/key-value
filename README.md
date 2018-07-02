# Key Value

Manage key-value data via HTTP.

## Development

### Prerequisite

  - Go 1.X
  - Redis

### Environment Variables

  - `REDIS_HOST`
  - `REDIS_PORT`

### Run

```bash
go run main.go
```

### Deploy

```bash
# Build binary
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-s' main.go

# Build docker image
docker build -t key-value .

# Run
docker run -e 'REDIS_HOST=localhost' -e 'REDIS_PORT=6379' -p 8080:8080 key-value
```


## Specification

### Common Errors

  - 500 Internal Server Error  
    There is an unexpected error. In most cases, 200 Ok will return when you request again.
  - 503 Service Unavailable  
    There is a problem with some backends.

### GET /v1/value

#### Query Parameter

  - `key` (string)

#### Response

  - 200 Ok
    ```json
    {
      "key": "string",
      "value": "string"
    }
    ```
  - 400 Bad Request : No `key` specified.

### POST /v1/value

#### Request Body

```json
{
  "key": "string",
  "value": "string"
}
```

#### Response

  - 200 Ok
    ```json
    {
      "key": "string",
      "value": "string"
    }
    ```
  - 400 Bad Request : Request body is invalid.

### GET /v1/list

#### Query Parameter

  - `key` (string)

#### Response

  - 200 Ok
    ```json
    {
      "key": "string",
      "list": [
        "string"
      ]
    }
    ```
  - 400 Bad Reqeust : No `key` specified.

### POST /v1/list

Append value to the list.

#### Request Body

```json
{
  "key": "string",
  "value": "string"
}
```

#### Response

  - 200 Ok
    ```json
    {
      "key": "string",
      "value": "string"
    }
    ```
  - 400 Bad Request : Request body in invalid.

### DELETE /v1/list

#### Query Parameter

  - `key` (string)
  - `value` (string) : Required if `delete_key` is `false`.
  - `delete_key` (boolean) : If `true`, key (including all values) will be deleted. Default `false`.

#### Response

  - 200 Ok
  - 400 Bad Request