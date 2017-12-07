Location service provides drop in replacement for existing Rails endpoints dealing with Heartbeat and current location

## :hammer: Building (without docker)

This project uses [dep](https://github.com/golang/dep) for dependency management. Please do not check in the `vendor` directory. You can obtain `dep` like so.

```
$ go get -u github.com/golang/dep/cmd/dep
$ dep ensure

```

This will add the dependencies to the `vendor` directory. You can then build the project with `make`.

```
$ make linux
```

To build for `darwin`, you can do `make mac`.

## :hammer: Building with docker

For building with docker, this project uses a multi-stage build introduced in docker 17.05.

```
$ docker build -t moove/location_service:latest .
```

## :rocket: Running

This service currently needs only one environment variable `LOCATION_DATABASE_URL`. It needs to be a connection string to the `mysql` server like so `user:pass@server:port/database`. Sometimes if the server name contains hyphens (`-`) as in the case of RDS, you might have to wrap the connection string like so `user:pass@tcp(server:port)/database`.
