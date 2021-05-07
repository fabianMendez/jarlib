# Jarlib

A simple tool to generate a self-contained jar of a dependency


## Install

```
go install github.com/fabianMendez/jarlib@latest
```

## Examples

Pack `com.google.cloud.sql:postgres-socket-factory:1.2.1` and all its dependencies into the file `postgres-socket-factory.jar`.

```
jarlib \
  --dependency com.google.cloud.sql:postgres-socket-factory:1.2.1 \
  --output postgres-socket-factory.jar
```

### Run a local server

Run a local http server on port `1234` (optional).

```
PORT=1234 jarlib serve
```

Then you can use curl or your browser to generate the library:

```
curl "http://localhost:1234/com.google.cloud:google-cloud-storage:1.112.0?javaVersion=1.7&projectName=google-cloud-storage" --output google-cloud-storage.jar
```
