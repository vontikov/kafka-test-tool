# Kafka test tool

## Build
```shell
$ make build
```

## Image
```shell
$ make image
```

## Run
```shell
$ kafka_test_tool -help
Usage of kafka_test_tool:
  -kafka-bootstrap string
        Kafka bootstrap servers (default "localhost:9092")
  -listen-address string
        Listen address. (default ":8080")
  -log-level string
        Log level: trace|debug|info|warn|error|none (default "info")
```

## Usage
```shell
$ curl -X POST localhost:8080/topics/some-topic 'some message'
```
