# rename the .env.sample to .env and fill the values
# it will export all the variables in the .env file
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

.PHONY: gen-env
gen-env:
	@echo "MAX_REQUESTS_PER_SECOND_IP=2" > .env
	@echo "MAX_REQUESTS_PER_SECOND_API_TOKEN=5" >> .env
	@echo "FREZEE_TIME_IN_SECONDS=60" >> .env


.PHONY: start-redis
start-redis:
	docker compose -f docker-compose.redis.yaml up -d

.PHONY: run-limiter
run-limiter: start-redis
	go run cmd/rate-limiter/main.go

.PHONY: run-server
run-server: start-redis
	go run cmd/server/main.go

# make sure to run docker compose up before running this command
FLAGS ?= -u http://localhost:8080/healthz -a FAKE_API_KEY -r 21 -c 10
.PHONY: run-cli-load
run-cli-load:
	go run cmd/cli/main.go load $(FLAGS)

.PHONY: test
test:
	go test ./...

LOAD_URL ?= http://server:8080/healthz
API_KEY ?=
REQUESTS ?= 10
CONCURRENCY ?= 5
.PHONY: load-test
load-test:
	LOAD_URL=$(LOAD_URL) API_KEY=$(API_KEY) REQUESTS=$(REQUESTS) CONCURRENCY=$(CONCURRENCY) \
	docker compose -f docker-compose.load-test.yaml up --abort-on-container-exit && docker compose -f docker-compose.load-test.yaml down


