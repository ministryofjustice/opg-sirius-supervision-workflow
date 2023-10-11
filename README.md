# OPG SIRIUS WORKFLOW

### Major dependencies

- [Go](https://golang.org/) (>= 1.19)
- [docker compose](https://docs.docker.com/compose/install/) (>= 2.0.0)

#### Installing dependencies locally:
(This is only necessary if dunning without docker)

- `yarn install`
- `go mod download`
---

## Local development

The application ran through Docker can be accessed on `localhost:8888/supervision/workflow/`.

**Note: Sirius is required to be running in order to authenticate. However, it also runs its own version of Workflow on port `8080`.
Ensure that after logging in, you redirect back to the correct port (`8888`)** 

To enable debugging and hot-reloading of Go files:

`make dev-up`

If you are using VSCode, you can then attach a remote debugger on port `2345`. The same is also possible in Goland.
You will then be able to use breakpoints to stop and inspect the application.

Additionally, hot-reloading is provided by Air, so any changes to the Go code (including templates) 
will rebuild and restart the application without requiring manually stopping and restarting the compose stack.

### Without docker

Alternatively to set it up not using Docker use below. This hosts it on `localhost:1234`
  
- `yarn install && yarn build `
- `go build main.go `
- `./main `

### Enabling code completion in .gotmpl files in GoLand

Go to `Settings -> Editor -> File Types -> Go template files` in your IDE and add `*.gotmpl` to the list of file name patterns.

Define the type of `{{ . }}` in the context of your template by adding a line like this at the top of the template:
`{{- /*gotype: github.com/ministryofjustice/opg-sirius-workflow/internal/server.WorkflowVars*/ -}}`

  -------------------------------------------------------------------
## Run Cypress tests

`make cypress`

To see the UI output you can still run 

`yarn && yarn cypress`

-------------------------------------------------------------------
## Run the unit/functional tests

`make unit-test`

-------------------------------------------------------------------
## Run Trivy scanning

`make scan`

