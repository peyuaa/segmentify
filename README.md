# segmentify

Package `segmentify` provides a service for dynamic user segmentation.

[![Go Reference](https://pkg.go.dev/badge/github.com/peyuaa/segmentify.svg)](https://pkg.go.dev/github.com/peyuaa/segmentify)

# How to start service

## Docker (recommended, everything is set up)
```
docker-compose up
```

## Manually (not recommended)
1) You need to have a running PostgreSQL instance. Init script is located in `./db/init.sql`.
Service uses environment variable [DB_CONNECTION_STRING](https://pkg.go.dev/github.com/lib/pq#hdr-Connection_String_Parameters) to connect to the database.

2) Install [go-swagger](https://goswagger.io/install.html).
It's necessary since the service uses swagger-generated file for hosting documentation.
3) Run `make run` in the root of the project.

# Documentation
Service documentation is available at `/docs` after starting the service.
By default, it's available at `http://localhost:9090/docs`. It contains richer description of endpoints and models.

## Get all existing segments (deleted included)
Returns all segments that were ever created in the system.
### Request
```http request
GET /segments HTTP/1.1
Host: localhost:9090
```
### Response
```http request
HTTP/1.1 200 OK
Content-Type: application/json
Connection: close

[{"id":1,"slug":"AVITO_VOICE_MESSAGES","is_deleted":false},{"id":2,"slug":"AVITO_RED_BUTTON","is_deleted":false}]
```

## Get segment by slug
Returns segment by slug.
### Request
```http request
GET /segments/AVITO_VOICE_MESSAGES HTTP/1.1
Host: localhost:9090
```
### Response
```http request
HTTP/1.1 200 OK
Content-Type: application/json
Connection: close

{"id":1,"slug":"AVITO_VOICE_MESSAGES","is_deleted":false}
```

## Create new segment
### Request
```http request
POST /segments HTTP/1.1
Content-Type: application/json; charset=utf-8
Host: localhost:9090

{"slug":"AVITO_RED_BUTTON"}
```
### Response
```http request
HTTP/1.1 201 Created
Content-Type: application/json
Location: http://localhost:9090/segments/AVITO_RED_BUTTON
Connection: close

{"id":2,"slug":"AVITO_RED_BUTTON","is_deleted":false}

```

## Delete segment
Mark segment as deleted.
### Request
```http request
DELETE /segments/AVITO_VOICE_MESSAGES HTTP/1.1
Host: localhost:9090
```

### Response
```http request
HTTP/1.1 204 No Content
Connection: close
```

## Change user segments
Add and remove segments for user.
Field `expired` is optional and specifies the date when segment should be removed from user.
### Request
```http request
POST /segments/users HTTP/1.1
Content-Type: application/json; charset=utf-8
Host: localhost:9090

{"id":73234,"add":[{"slug":"AVITO_RESEARCH_AMOGUS","expired":"2025-01-02T15:04:06Z"},{"slug":"AVITO_CHINESE_MARKET"}],"remove":[{"slug":"AVITO_RED_BUTTON"}]}
```

### Response
```http request
HTTP/1.1 200 OK
Content-Type: application/json
Connection: close

{"segments":[{"slug":"AVITO_RESEARCH_AMOGUS"},{"slug":"AVITO_CHINESE_MARKET"}]}
```

## Get user segments (expired not included)
Returns all segments that are currently assigned to user.
```http request
GET /segments/users/73234 HTTP/1.1
Host: localhost:9090
```

### Response
```http request
HTTP/1.1 200 OK
Content-Type: application/json
Connection: close

[{"slug":"AVITO_RESEARCH_AMOGUS"},{"slug":"AVITO_CHINESE_MARKET"}]
```

## Get user history
Returns all changes of user segments in the specified time range.
Changes are sorted by date in ascending order. 

By default, if you run the service using docker-compose, time zone is set to Europe/Moscow.
```http request
GET /segments/users/73234/history?from=2023-08-30&to=2023-08-31 HTTP/1.1
Host: localhost:9090
```

### Response
```http request
HTTP/1.1 200 OK
Content-Type: text/plain; charset=utf-8
Connection: close

{"link":"http://localhost:9090/history/73234/2023-08-30/2023-08-31/history.csv"}
```

### CSV-history file example
```csv
73234,AVITO_RED_BUTTON,add,2023-08-30T17:36:28Z
73234,AVITO_RED_BUTTON,remove,2023-08-30T17:38:11Z
73234,AVITO_RESEARCH_AMOGUS,add,2023-08-30T17:38:11Z
73234,AVITO_CHINESE_MARKET,add,2023-08-30T17:38:11Z
```