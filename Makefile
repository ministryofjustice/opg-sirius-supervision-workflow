.PHONY: cypress

all: go-lint build-all unit-test scan cypress down

build:
	docker compose -f docker/docker-compose.ci.yml build --parallel workflow

build-all:
	docker compose -f docker/docker-compose.ci.yml build --parallel workflow json-server test-runner cypress

go-lint:
	docker compose -f docker/docker-compose.ci.yml run --rm go-lint

test-results:
	mkdir -p -m 0777 test-results

setup-directories: test-results

unit-test: setup-directories
	docker compose -f docker/docker-compose.ci.yml run --rm test-runner gotestsum --junitfile test-results/unit-tests.xml -- ./... -coverprofile=test-results/test-coverage.txt

scan:
	trivy image sirius/sirius-workflow:latest

up:
	docker compose -f docker/docker-compose.ci.yml up --build -d deputy-hub

down:
	docker compose -f docker/docker-compose.ci.yml down