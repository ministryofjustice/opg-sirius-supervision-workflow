# json-server

json-server is the mock API we use for local development and testing. It is a simple Node Express app reading from a JSON file, which makes it flexible for our needs.

If you need data, add it into `db.json` in the format you require it, making sure to include an `id` if you are using a plural route (i.e. a route that could return many different entries).
json-server provides functionality for nested routes and parent/child relationships but if you require custom routing (e.g. you always want the same data returned, regardless of the id), you can add these to `routes.json`.

If you want to inspect the data, json-server is served on port 3000 using the `docker-compose.dev.yml` so you can visit it in a browser.

For more advanced customisation, you can create your own Express middleware and include it in the `serve` script in `package.json`.

## Middleware

### Error Rerouter

To allow us to test validation and error handling where we expect bad requests to be returned from the API, the `error-rerouter` middleware enables requests to be rerouted to an errors table in `db.json` and returned with a `400` status code.
To add this to your Cypress test, simply set the `fail-route` cookie and add the expected error to the database, using the cookie's value as the id.

e.g.:

`cy.setCookie("fail-route", "notes")`

This will reroute to `/errors/notes` and return the data stored in the `errors` object with the id `notes`.

---

### Success Rerouter

This is the same as above, except to handle successful requests that are hard to fit into the JSON model. Follow the same pattern as the error rerouter but use the `success-route` and put the data in the `successes` object.

Please read the very helpful readme on the [json-server GitHub](https://github.com/typicode/json-serve) for more information.
