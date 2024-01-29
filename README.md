# Fast Loki connection for Slog

## Installation
```bash
go get -u github.com/wensiet/logmod
```

## Usage example
```go
package main

import "github.com/wensiet/logmod"

func main() {
	logger := logmod.New(logmod.Options{
		Env:     "local", // local OR production OR test
		Service: "",      // service name-tage
		Loki: struct {
			Host string
			Port int
		}{ // Loki's connection string (host:port), by default endpoint is http://localhost:3100/loki/api/v1/push
			Host: "localhost",
			Port: 3100,
		}, 
	}) // IF LOKI IS NOT AVAILABLE, IT WILL USE STDOUT

	logger.Info("Hello World!")
}
```

## Loki fast [deploy](https://grafana.com/docs/loki/latest/setup/install/docker/) v2.9.4 (don't forget to change docker network)

```bash
wget https://raw.githubusercontent.com/grafana/loki/v2.9.4/production/docker-compose.yaml -O docker-compose.yaml
docker-compose -f docker-compose.yaml up
```


## Original Loki client and handlers
https://github.com/samber/slog-loki
https://github.com/grafana/loki-client-go
