all: test build scan cypress down

.PHONY: cypress

test-results:
	mkdir -p -m 0777 test-results cypress/screenshots .trivy-cache

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

scan: setup-directories
	docker compose run --rm trivy image --format table --exit-code 0 311462405659.dkr.ecr.eu-west-1.amazonaws.com/sirius/sirius-workflow:latest
	docker compose run --rm trivy image --format sarif --output /test-results/trivy.sarif --exit-code 1 311462405659.dkr.ecr.eu-west-1.amazonaws.com/sirius/sirius-workflow:latest

cypress: setup-directories
	docker compose up -d --wait workflow json-server
	docker compose run --build --rm cypress

up:
	docker compose up --build -d workflow

dev-up:
	docker compose run --rm yarn
	docker compose -f docker-compose.yml -f docker/docker-compose.dev.yml build workflow
	docker compose -f docker-compose.yml -f docker/docker-compose.dev.yml up workflow yarn json-server

down:
	docker compose down
