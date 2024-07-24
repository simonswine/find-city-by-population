FROM golang:1.22.5 AS build

WORKDIR /app

RUN apt-get update && apt-get -y install unzip

COPY Makefile /app
RUN make data

COPY go.mod go.sum /app
RUN go mod download

COPY main.go /app
RUN CGO_ENABLED=0 go build -o /usr/bin/find-city-by-population  /app

FROM gcr.io/distroless/static

EXPOSE 8081
USER nobody
WORKDIR  /app

COPY --from=build /usr/bin/find-city-by-population /usr/bin/find-city-by-population
COPY --from=build /app/data /app/data

CMD ["find-city-by-population"]
