all: go-lint unit-test build scan cypress down

.PHONY: cypress

test-results:
	mkdir -p -m 0777 test-results cypress/screenshots .trivy-cache

setup-directories: test-results

go-lint:
	docker compose run --rm go-lint

build:
	docker compose build --parallel workflow

build-all:
	docker compose build --parallel workflow json-server test-runner cypress

unit-test: setup-directories
	docker compose run --rm test-runner gotestsum --junitfile test-results/unit-tests.xml -- ./... -coverprofile=test-results/test-coverage.txt

scan: setup-directories
	docker compose run --rm trivy image --format table --exit-code 0 311462405659.dkr.ecr.eu-west-1.amazonaws.com/sirius/sirius-workflow:latest
	docker compose run --rm trivy image --format sarif --output /test-results/trivy.sarif --exit-code 1 311462405659.dkr.ecr.eu-west-1.amazonaws.com/sirius/sirius-workflow:latest

cypress: setup-directories
	docker compose up -d --wait workflow
	docker compose run --rm cypress run --env grepUntagged=true

up:
	docker compose up --build -d workflow

dev-up:
	docker compose run --rm yarn
	docker compose run --rm yarn build
	docker compose -f docker-compose.yml -f docker/docker-compose.dev.yml build --no-cache workflow
	docker compose -f docker-compose.yml -f docker/docker-compose.dev.yml up workflow json-server

down:
	docker compose down
