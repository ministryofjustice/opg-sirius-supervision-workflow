# OPG SIRIUS WORKFLOW

### Major dependencies

- [Go](https://golang.org/) (>= 1.16)
- [Pact](https://github.com/pact-foundation/pact-ruby-standalone) (>= 1.88.3)
- [docker-compose](https://docs.docker.com/compose/install/) (>= 1.27.4)

#### Installing dependencies locally: 
- `yarn install`
- `go mod download`
-------------------------------------------------------------------

## Local development

The application ran through Docker can be accessed on `localhost:8888/supervision/workflow/`.

**Note: Sirius is required to be running in order to authenticate. However, it also runs its own version of Workflow on port `8080`.
Ensure that after logging in, you redirect back to the correct port (`8888`)** 

To enable debugging and hot-reloading of Go files:

`docker-compose -f docker/docker-compose.yml -f docker/docker-compose.dev.yml up --build`

If you are using VSCode, you can then attach a remote debugger on port `2345`. The same is also possible in Goland.
You will then be able to use breakpoints to stop and inspect the application.

Additionally, hot-reloading is provided by Air, so any changes to the Go code (including templates) 
will rebuild and restart the application without requiring manually stopping and restarting the compose stack.

### Without docker

Alternatively to set it up not using Docker use below. This hosts it on `localhost:1234`
  
- `yarn install && yarn build `
- `go build main.go `
- `./main `

  -------------------------------------------------------------------

## Run Pact and Cypress tests

Generate the pact file which mimics Sirius and runs workflow tests

`go test ./...`
 
To generate the pact file with no cache

`go test ./... -count=1`
 
Run Cypress tests against the pact copy of the service
 
`docker-compose -f docker/docker-compose.cypress.yml up -d --build `
 
`yarn && yarn cypress `
    
-------------------------------------------------------------------
## Run the Pact tests in more detail

test sirius files: `test-sirius`

test server files: `test-server`

-------------------------------------------------------------------
## Noted issues:
- Can't run locally if the pact stub is still running
