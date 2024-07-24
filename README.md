# find-city-by-population

Example go app, that can leaves room to optimize cpu / memory usage.

![Screenshot find city by population](screenshot.png)


## Develop

Requires a golang environment and then it can be run as follows:

```
# Download city dataset locally
$ make data

# After running the app, it will listen on port http://localhost:8081
$ go run ./

# Run happy path test
$ make test
```

## Build docker image

```
$ make build
```

## Build & push docker image

```
$ make push IMAGE_PREFIX=my-docker-user IMAGE_TAG=iteration-1
```

## City dataset

The dataset is from https://www.geonames.org/ and it is licensed under a Creative Commons Attribution 4.0 License.
