# Hash and Encode Password API

## Prerequisite
- golang v.1.16 or later

## Build
```
# For Linux
> GOOS=linux go build -o main_linux

# For MacOSX
> GOOS=darwin go build -o main_darwin
```

## Run
```
# For Linux
> ./main_linux

# For MacOSX
> ./main_darwin
```

## REST API samples
### Hash and Encode a Password String
```
> curl http://localhost:8080/hash --data "password=angryMonkey"
1

> curl http://localhost:8080/hash --data "password=chaosMonkey"
2

# failure cases
> curl http://localhost:8080/hash --data "password="
password field is required (400 Bad Request)

> curl http://localhost:8080/hash --data "invalid_key=value"
password field is required (400 Bad Request)

```

### Get Password
```
> curl http://localhost:8080/hash/1
ZEHhWB65gUlzdVwtDQArEyx-KVLzp_aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A-gf7Q==%

# failure case
> curl http://localhost:8080/hash/100
id 100 not found (404 Not Found)
```

### Get Statistics
```
> curl http://localhost:8080/stats
{"total":2,"totalTime":240,"average":120}

# failure case
> curl http://localhost:8080/stats
{"message":"stats POST:/hash not found"} (404 Not Found)
```

### Shutdown Gracefully
```
> curl http://localhost:8080/shutdown
```

## Test
This shell script run `go test` for all directories and show you test result in html
```
./test.sh
```

## Code structure
```
.
├── common
│   └── configs.go
├── handlers
│   ├── handler_test.go
│   ├── handler_utils.go
│   ├── handler.go
│   └── middleware.go
├── router
│   └── router.go
├── storages
│   ├── password_test.go
│   ├── password.go
│   ├── stats_test.go
│   └── stats.go
├── testutils
│   └── testutils.go
├── go.mod
├── go.sum
├── main.go
├── README.md
└── test.sh
```

- `/common/`<br>
  common constants, utilities functions
- `/handlers/`<br>
  http handlers
  middlwares
- `/router/`<br>
  router for REST API
- `/storages/`<br>
  data objects to store Password and Stats data
- `/testutils/`<br>
  go test utilities functions
- `main.go`<br>
  entrypoint to start server
- `test.sh`<br>
  script to run go test for all directories