[![CircleCI](https://circleci.com/gh/MOOVE-Network/location_service/tree/eta_service.svg?style=svg&circle-token=3342db5a1501e6add8807b75b18683feea41aec6)](https://circleci.com/gh/MOOVE-Network/location_service/tree/eta_service)

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

This service currently needs the following environment variables 

* `LOCATION_DATABASE_URL`. It needs to be a connection string to the `mysql` server like so `user:pass@server:port/database`. Sometimes if the server name contains hyphens (`-`) as in the case of RDS, you might have to wrap the connection string like so `user:pass@tcp(server:port)/database`.
* `LOCATION_MAPS_API_KEY` is the google maps API key that is required
* `FCM_API_KEY` is the API Key to Firebase Cloud Messaging service
* `FCM_TOPIC_PREFIX` is the topic prefix used for messages sent with this service
* `LOCATION_REDIS_URL` is the redis url for caching trip locations it defaults to `localhost:6379`

### Configuring systemd

* Copy `location_service.service.example` to `/lib/systemd/system/location_service.service`. 
* Once done, make sure to `sudo chmod 755 /lib/systemd/system/location_service.service`.
* To enable`sudo systemctl enable location_service.service`
* To start `sudo systemctl start location_service`
* To look at the logs you can use `sudo journalctl -f -u location_service`
