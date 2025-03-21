# Rate Limiter & CLI

## Project structure

```bash
├── cmd                # Application entry points
│   ├── cli            # Starts the stress-test CLI
│   ├── rate-limiter   # Executes the rate-limiting logic
│   └── server         # Runs the REST API server
├── internal           # Internal application implementations
│   ├── cli            # CLI-specific logic and utilities
│   │   ├── cmd        # Command definitions and execution logic
│   │   └── flags       # CLI flags parsing and configuration
│   └── server         # Server-specific logic and components
│       └── middleware # HTTP middleware components
└── pkg                # Reusable shared libraries
    ├── goten          # Stress-testing utilities and core logic
    ├── limiter        # Rate-limiting logic and interfaces
    └── redis          # Redis client wrapper implementing required interfaces
```

## Overview

This monorepo provides two core components:
- **CLI**: A command-line tool responsible for executing HTTP stress tests. Built with Cobra, it offers a single `load` command. The actual stress-testing logic resides within the reusable `pkg/goten` package, promoting modularity and ease of maintenance.
- **Server**: A web server that exposes a health check route (`/healthz`) and applies rate-limiting middleware to HTTP requests. It utilizes the Echo framework and depends on two key packages:
  - `pkg/limiter`: Implements comprehensive rate-limiting logic.
  - `pkg/redis`: Provides a Redis client wrapper conforming to interfaces expected by the limiter, ensuring seamless integration and ease of swapping implementations if required.


## Rate Limiter 

### How to run
```bash
# install dependencies
go mod tidy

# create a .env file with default values, you can edit these values as you wish
make gen-env

# execute only the rate limiter logic
make run-limiter

# run server that implements the rate limiter middleware
make run-server

# execute request to the server
curl localhost:8080/healthz

# execute request to the server with api_key header
curl -H 'api_key: 123' localhost:8080/healthz

# execute request to the server with api_key header until you force the stop 
# if you are using the default values in .env it should start failing soon
while true; do curl -H 'api_key: 12345' localhost:8080/healthz; sleep 0.15; done
```

### Testing
```bash
# unit testing
make test

# load test
# make sure the port 8080 is free
# you can update the .env to obtain different results
make load-test

# load test passing custom values
make load-test LOAD_URL=http://server:8080/healthz API_KEY=333 REQUESTS=20 CONCURRENCY=10
```

## CLI

### How to run
```bash
# install dependencies
go mod tidy

# run
go run cmd/cli/main.go load -u https://pokeapi.co/api/v2/pokemon/1 -r 50 -c 10

# with api key in header
go run cmd/cli/main.go load -u http://localhost:8080/healthz -a 222 -r 50 -c 10

# run with binary
go build -o bin/cli ./cmd/cli/main.go
./bin/cli load -u https://pokeapi.co/api/v2/pokemon/1 -r 50 -c 10

#
# USING DOCKER
#
# build docker image
docker build -t cli -f Dockerfile.cli .

# run docker image
docker run --rm cli -u https://pokeapi.co/api/v2/pokemon/1 -r 50 -c 10
```

