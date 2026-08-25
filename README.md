# OPG SIRIUS WORKFLOW

### Major dependencies

- [Go](https://golang.org/) (>= 1.19)
- [docker compose](https://docs.docker.com/compose/install/) (>= 2.0.0)

#### Installing dependencies locally:
(This is only necessary if dunning without docker)

- `npm install`
- `go mod download`
---

## Local development

The application ran through Docker can be accessed on `localhost:8888/supervision/workflow/`.

To enable debugging and hot-reloading of Go files:

`make dev-up`

If you are using VSCode, you can then attach a remote debugger on port `2345`. The same is also possible in Goland.
You will then be able to use breakpoints to stop and inspect the application.

Additionally, hot-reloading is provided by Air, so any changes to the Go code (including templates)
will rebuild and restart the application without requiring manually stopping and restarting the compose stack.

Note that the `local-proxy` service is required to allow for both the front and back ends to access json-server as a mock
API. It also handily logs all requests to the console.

```
make build
# switch to opg-sirius repo
make dev-up
```

### With Sirius

If you want to run Workflow locally against a local Sirius instance, you can use `make sirius-up`. As this will not use
json-server as a mock, you will first need to log into Sirius in order for the user to have a valid session and then
redirect back to Workflow using the correct port (`8888`).

### Enabling code completion in .gotmpl files in GoLand

Go to `Settings -> Editor -> File Types -> Go template files` in your IDE and add `*.gotmpl` to the list of file name patterns.

Define the type of `{{ . }}` in the context of your template by adding a line like this at the top of the template:
`{{- /*gotype: github.com/ministryofjustice/opg-sirius-workflow/internal/server.WorkflowVars*/ -}}`

  -------------------------------------------------------------------
## Run *one* Cypress test headless (i.e. not in UI)
`make cypress-single SPEC=bonds/bonds_list.cy.js`

## Run *a suite* of Cypress tests headless e.g. bonds subfolder
`make cypress-single SPEC=bonds/`

## Run *all* the Cypress tests headless
`make build-all` (optional) <br>
`make cypress`


To see the UI output you can still run
`make up && npm run cypress`

-------------------------------------------------------------------
## Run the unit/functional tests
`make unit-test`
