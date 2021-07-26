![Logo_Lock_Transparent](https://user-images.githubusercontent.com/12732454/126949778-df08d3b0-5233-42d7-9757-c4652e0bdbff.png)


# Whaler-api

Whaler-api is a Go application backend serving up REST and GraphQL endpoints for the Whaler client app. WebSockets are used to power custom real-time collaborative note taking and real-time model updates.

## What is Whaler? 

Whaler is a platform consisting of a native mac app ([Whaler](https://github.com/ZHRhodes/Whaler)) and a heroku-hosted backend. The goal is to enable real time sales outreach collaboration right on top of the organization's Salesforce data. This is accomplished by allowing the user to sign in to their organization's Salesforce, importing subsets of their data into the frontend app, and enhancing that data with additional constructs that power our features. 

These features include: 
* Kanban-style progress tracking for accounts and contacts
* Per-account tasks feature with ability to create, edit, set due date, and complete tasks
* Assign accounts, contacts, and tasks to users within your organization
* Real time collaborative editor for account notes, communicating over websockets with operational transform conflict resolution
* Real time data model updates over websockets

Because a picture speaks a thousand words, here are two key pages within the app:

<img width="1345" alt="Screen Shot 2021-07-25 at 7 55 51 PM" src="https://user-images.githubusercontent.com/12732454/126926937-3ffb596c-f8af-489c-8fef-d0a213cafa11.png">

<img width="1340" alt="Screen Shot 2021-07-25 at 7 56 45 PM" src="https://user-images.githubusercontent.com/12732454/126926982-1609ec78-9c7d-4ae2-a72f-d9a14073c856.png">

<img width="1340" alt="Screen Shot 2021-07-25 at 7 58 02 PM" src="https://user-images.githubusercontent.com/12732454/126927102-e9f19bf7-eb2d-4506-a55d-505dde54820e.png">

## Whaler-api Technical Overview

Looking at the whole picture, this repo consists of the left half of this diagram:

![Whaler Technical Overview](https://user-images.githubusercontent.com/12732454/126923083-f529c8ba-a43b-49d5-976d-745047c1a230.png)

Whaler-api is currently hosted on Heroku. Here's a look at the high level project structure:
<img width="641" alt="Mind Map 4" src="https://user-images.githubusercontent.com/12732454/126924006-7cc46a3a-7de1-4f5f-9cfb-2c2ae750f044.png">

Let's step through each node to get an idea of how it all works. 

### Middleware

![IMG_24CA09751158-1](https://user-images.githubusercontent.com/12732454/126927906-9163a0cd-5d38-45f9-94cf-3fc8bf5c6707.jpeg)

In `jwtAuth.go`, a simple bit of middleware authenticates the jwt token attached to incoming requests. If the token is valid, then the userId will be extracted from the token and added to the request context.

### Controllers

Whaler-api only uses REST endpoints for authentication related requests. In `authController.go`, you'll find the handlers for authenticating and refreshing the token. These functions parse the token and then call into the `Token` model to do more work.

### Graph

Whaler-api uses GraphQL for all data operations. The full schema is available in this directory as `schema.graphqls`. Whaler-api uses a Go module called `gqlgen` to generate the GraphQL-related boilerplate. This schema-first approach lets us focus on defining our API needs and then leveraging the strong typing of GraphQL and Go to catch any inconsistencies as the API changes. 

In order to generate the boilerplate with the current schema, use `go run github.com/99designs/gqlgen generate`. The generated code will be placed in `generated.go` and includes Go interfaces defining the schema and Go functions implementing the actual GraphQL operations themselves. You don't need to worry much about the contents of this file; `generate` will place empty resolver methods in `schema.resolvers.go`. Each query and mutations in the schema will map to a resolver function added to this file, and all you have to do then is fill in those resolver functions. This effectively abstracts almost all of the implementation details of GraphQL and lets us focus on the parts that are specific to this project. Amazing!

`generate` will create Go structs for each type in the schema, but you can also map them to your own structs if you wish. Mappings can be added to the top-level `gqlgen.yml` file. For the sake of having more control over struct definitions, most Whaler models are mapped.

A schema explorer called `playground` runs on the `/schema` route accessable on the browser. This playground introspects your schema and presents an explorer alongside a terminal where you can execute queries and mutations. This is very useful for using and testing resolvers. **Note**: These requests do still require a token, so you'll need to generate one from the REST endpoint and add it as the `Authorization` header in the playground interface.

### Models

This package contains all of the Whaler-api database models. `gorm` is the ORM that facilitates communication with the connected Postgres database. Models include `account`, `contact`, `accountAssignmentEntry`, `note`, `task`, `token`, and more. Each class has functions that provide CRUD operations on the data and are usually called from GraphQL resolvers.

Access tokens, defined in `token.go`, are of the JWT format, and encode the userId in addition to the standard JWT claims. Refresh tokens, also defined in `token.go`, contain the userId, an expiration date, and a randomly generated 256 byte hash. The current expiration time is 90 days for refresh tokens. 

### Websocket

### OT


### Roadmap
Ideally, the Salesforce integration would be moved to the backend. The frontend would still be responsible for initiating the integration process, but once a Salesforce token is obtained, it would be sent to the backend, where the integration would then be managed. Moving all the Salesforce data management to the backend would free up the app to focus on being a great frontend. As a part of that migration, I would abstract the few direct Salesforce references (e.g. SalesforceId) behind a generic CRM integration interface. Then, adding support for HubSpot, Pipedrive, or any other CRM would be much simpler. 

Additionally, this project could use a bit more delineation. Over time, the `models` package would continue to expand and take on too much responsibility. I'd like to further reduce the scope of the existing packages by creating new ones and splitting them up better. Because this project started as a REST api and later switched to GraphQL, there are some leftovers from that transition that could use an update. 

## Running Locally

Make sure you have [Go](http://golang.org/doc/install) version 1.12 or newer and the [Heroku Toolbelt](https://toolbelt.heroku.com/) installed.

```sh
$ git clone https://github.com/heroku/whaler-api.git
$ cd whaler-api
$ go build -o bin/whaler-api -v . # or `go build -o bin/whaler-api.exe -v .` in git bash
github.com/mattn/go-colorable
gopkg.in/bluesuncorp/validator.v5
golang.org/x/net/context
github.com/heroku/x/hmetrics
github.com/gin-gonic/gin/render
github.com/manucorporat/sse
github.com/heroku/x/hmetrics/onload
github.com/gin-gonic/gin/binding
github.com/gin-gonic/gin
github.com/heroku/whaler-api
$ heroku local
```

Your app should now be running on [localhost:5000](http://localhost:5000/).

## Deploying to Heroku

```sh
$ heroku create
$ git push heroku master
$ heroku open
```

or

[![Deploy](https://www.herokucdn.com/deploy/button.png)](https://heroku.com/deploy)


## Documentation

For more information about using Go on Heroku, see these Dev Center articles:

- [Go on Heroku](https://devcenter.heroku.com/categories/go)
