<h1> OPG SIRIUS WORKFLOW </h1>

### Major dependencies

- [Go](https://golang.org/) (>= 1.16)
- [Pact](https://github.com/pact-foundation/pact-ruby-standalone) (>= 1.88.3)
- [docker-compose](https://docs.docker.com/compose/install/) (>= 1.27.4)

  <p> install dependencies: </p>
  ### `yarn install`
  <p> If any packages are added to go mod call </p>
  ### `go mod download `

-------------------------------------------------------------------
  ### `yarn test-sirius`
  <p> test sirius files </p>

   ### `yarn test-server`
  <p> test server files </p>

<h2> Set up the service </h2>
  <p> Make sure Sirius is running on localhost:8080 </p>
  <p> Once Sirius is running the run below command to launch Workflow locally, it should be on localhost:8888/supervision/workflow/ </p>
  
  ### `docker-compose -f docker/docker-compose.yml up -d --build `

  <h3> Alternatively to set it up not using Docker use below. This hosts it on localhost:1234 : </h3>
  
  ### `yarn install && yarn build `
  ### `go build main.go `
  ### `./main `

  -------------------------------------------------------------------

<h2> Run Pact and Cypress tests </h2>

<p> Generate the pact file which mimics Sirius and runs workflow tests</p>

 ### `go test ./...`
 
 <p> To generate the pact file with no cache </p>

 ### `go test ./... -count=1`
 
 <p> Run Cypress tests against the pact copy of the service </p>
 
 ### `docker-compose -f docker/docker-compose.cypress.yml up -d --build `
 
 ### `yarn && yarn cypress `
    
  -------------------------------------------------------------------
<h2> Run the Pact tests in more detail</h2>

  ### `test-sirius`
  <p> test sirius files </p>

  ### `test-server`
  <p> test server files </p>

    -------------------------------------------------------------------

  <h2> Noted issues: </h2>
  <ul>
  <li> Can't run locally if the pact stub is still running </li>
  </ul>

  -------------------------------------------------------------------


 