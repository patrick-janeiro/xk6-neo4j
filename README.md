# xk6-neo4j

A k6 extension to interact with neo4j.

This can be particularly useful in seeding/cleaning data before performance testing. See example usage.

## Currently Supported Commands 

- Supports http and https connections.
- Support basic auth.
- RunCypherInSession: Runs a cypher string within a new session and returns the result.

## Examples: 

```js

import http from 'k6/http';
import neo4j from 'k6/x/neo4j'
import { check, sleep } from 'k6';

const dbConfig= {
	address: "localhost:7474"
	user: "neo4j", 
	httpPort: 0 ,   // If HTTPS port is provided httpPort will not be used
	httpsPort: 7687,
	password: "admin",
	maxTransactionRetryTime: 30
}

export function setup() {
   
  const cypher = `
    LOAD CSV WITH HEADERS FROM 'https://data.neo4j.com/northwind/customers.csv' AS row 
    MERGE (c:Company {companyId: row.Id, companyName: row.companyName});
  `
 

  const driver  = new neo4j.Driver(dbConfig)
  driver.runCypherInSession(cypher)
}

export default function() {
  // run test here :)
  const response = http.get("https://my-api/customers");
  sleep(1);

  check(response, {
	  'status code was 200': (r) => r.status === 200,
	})
}

export function teardown() {
  const cypher = `
    MATCH (n)
    DETACH DELETE n;
  ` 

  const driver  = new neo4j.Driver(dbConfig)
  driver.runCypherInSession(cypher)
}

```
