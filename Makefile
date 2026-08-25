all: test build cypress down

.PHONY: cypress

test-results:
	mkdir -p -m 0777 test-results cypress/screenshots

setup-directories: test-results

test:
	$(MAKE) -j 3 go-lint gosec unit-test

go-lint:
	docker compose run --rm go-lint

gosec: setup-directories
	docker compose run --rm gosec

unit-test: setup-directories
	docker compose run --rm test-runner gotestsum --junitfile test-results/unit-tests.xml -- ./... -coverprofile=test-results/test-coverage.txt

build-all:
	docker compose build --parallel workflow json-server test-runner cypress

build:
	docker compose build workflow

cypress: setup-directories
	docker compose run --build --rm cypress

cypress-single: setup-directories
	docker compose run --rm cypress run --spec cypress/e2e/$(SPEC)

up:
	docker compose -f docker-compose.yml up --build -d local-proxy

dev-up:
	docker compose -f docker-compose.yml -f docker/docker-compose.dev.yml build workflow
	docker compose -f docker-compose.yml -f docker/docker-compose.dev.yml up workflow watch-assets local-proxy

sirius-up:
	docker compose -f docker-compose.yml -f docker/docker-compose.standalone.yml up --build -d workflow

down:
	docker compose down
