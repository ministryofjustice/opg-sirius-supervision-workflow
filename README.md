<h1> OPG SIRIUS WORKFLOW </h1>

  ### `yarn install`
  <p> install dependencies </p>

  ### `test-sirius`
  <p> test sirius files </p>

   ### `test-server`
  <p> test server files </p>

<h2> Set up the service </h2>
  <p> Make sure Sirius is running </p>
  
  ### `docker-compose -f docker/docker-compose.yml up -d --build `

<h2> Run the tests </h2>
<p> Generate the pact file which mimics Sirius and tests files in Sirius folder</p>

 ### ` `
 
 <p> Run Cypress tests against the pact copy of the service </p>
 
 ### `docker-compose -f docker/docker-compose.cypress.yml up -d --build `
 
 ### `yarn && yarn cypress `
    
  <p> If any packages are added to go mod call </p>

   ### `go mod download `

  <p> Noted issues: </p>
  <ul>
  <li> Can't run locally if the pact stub is still running </li>
  </ul>