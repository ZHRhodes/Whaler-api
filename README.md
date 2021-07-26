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
<img width="703" alt="Mind Map 8" src="https://user-images.githubusercontent.com/12732454/126958935-0667384b-ae53-42d7-b803-f96d9b3c2547.png">

Let's step through each node to get an idea of how it all works. 

### Middleware

<img width="500" alt="IMG_DA5BD9A32359-1 copy" src="https://user-images.githubusercontent.com/12732454/126959995-926a5c54-27b5-41bb-93e1-999ae9fe0300.png">

In `jwtAuth.go`, a simple bit of middleware authenticates the jwt token attached to incoming requests. If the token is valid, then the userId will be extracted from the token and added to the request context. This middleware is added to the responder chain during the initialization of the router in `main.go`. 

### Controllers

<img width="500" alt="IMG_05FDC30FEF94-1 copy" src="https://user-images.githubusercontent.com/12732454/126959700-376594b1-0952-47e5-aa4c-de0f8cdd2f94.png">


Whaler-api only uses REST endpoints for authentication related requests. In `authController.go`, you'll find the handlers for authenticating and refreshing the token. These functions parse the token and then call into the `Token` model to do more work.

### Graph

<img width="500" alt="IMG_8DA394044771-1 copy" src="https://user-images.githubusercontent.com/12732454/126960018-58e6fe71-9ee5-45ae-a2cb-876776f6001c.png">

Whaler-api uses GraphQL for all data operations. The full schema is available in this directory as `schema.graphqls`. Whaler-api uses a Go module called `gqlgen` to generate the GraphQL-related boilerplate. This schema-first approach lets us focus on defining our API needs and then leveraging the strong typing of GraphQL and Go to catch any inconsistencies as the API changes. 

In order to generate the boilerplate with the current schema, use `go run github.com/99designs/gqlgen generate`. The generated code will be placed in `generated.go` and includes Go interfaces defining the schema and Go functions implementing the actual GraphQL operations themselves. You don't need to worry much about the contents of this file; `generate` will place empty resolver methods in `schema.resolvers.go`. Each query and mutations in the schema will map to a resolver function added to this file, and all you have to do then is fill in those resolver functions. This effectively abstracts almost all of the implementation details of GraphQL and lets us focus on the parts that are specific to this project. Amazing!

`generate` will create Go structs for each type in the schema, but you can also map them to your own structs if you wish. Mappings can be added to the top-level `gqlgen.yml` file. For the sake of having more control over struct definitions, most Whaler models are mapped.

A schema explorer called `playground` runs on the `/schema` route accessable on the browser. This playground introspects your schema and presents an explorer alongside a terminal where you can execute queries and mutations. This is very useful for using and testing resolvers. **Note**: These requests do still require a token, so you'll need to generate one from the REST endpoint and add it as the `Authorization` header in the playground interface.

### Models

<img width="500" alt="IMG_C9B2799EDAFE-1 copy" src="https://user-images.githubusercontent.com/12732454/126960043-25023299-94e8-499a-97e9-82237a44ab1b.png">

This package contains all of the Whaler-api database models. `gorm` is the ORM that facilitates communication with the connected Postgres database. Models include `account`, `contact`, `accountAssignmentEntry`, `note`, `task`, `token`, and more. Each class has functions that provide CRUD operations on the data and are usually called from GraphQL resolvers.

Access tokens, defined in `token.go`, are of the JWT format, and encode the userId in addition to the standard JWT claims. Refresh tokens, also defined in `token.go`, contain the userId, an expiration date, and a randomly generated 256 byte hash. The current expiration time is 90 days for refresh tokens. 

### Websocket

<img width="500" alt="IMG_F5B32B4D93E1-1 copy" src="https://user-images.githubusercontent.com/12732454/126959549-2d8357c0-7c1a-44b9-b702-caa5d5e01d22.png">

Websocket connections are represented as `Clients` defined in `clients.go`. When a new connection request is received (`websocket.go`), a `Client` is created and added to the `Pool` associated with the `resourceId` passed with the request. If this is the first connection seen for that `resourceId`, then a new `Pool` is created. When a `Client` receives a message, it does a bit of pre-processing in `client.go` before sending the message to the `Process` function defined in `process.go`. From here, the message type is used to decode the message data into the appropriate struct, which is then sent into the correct handler function. 

Messages can easily be broadcast to an entire `Pool` by sending a `SocketMessage` into the `Pool.broadcast` channel. The `Pool` will iterate over each of its clients and send it the new message. If you're trying to replicate an incoming message across other members of the pool, then you can leverage the `message.SenderId` property to ensure the data will not be sent back to the originating `Client`. 

New connections arrive with a grouping identifier, typically the resource id of the object being viewed or collaborated on. In the case of a websocket client wanting updates for the accounts they're tracking, this id will be the user's `organizationId`. The first iteration of real time model updates simply fires a message containing the resource id of the updated model to each client registered with the organizationId. While this dragnet approach may be a bit inefficient, it effectively suits the needs of the MVP. This `ResourceUpdate` message and all other websocket messages are defined in `messages.go`.

The models package defines a `ChangeConsumer` interface that is implemented by `websocket.ChangeConsumer`. The websocket implementation is set as the consumer for the models package. When models are changed, these changes are forwarded to the `ChangeConsumer` for further handling. In this case, that results in a message being broadcast to the `Pool` for the resourceId indicating that the model has changed. Currently, this triggers a reload on the client side, but later this message could be expanded to contain the new model itself.

### OT

<img width="500" alt="IMG_110634F00130-1 copy" src="https://user-images.githubusercontent.com/12732454/126960078-17ba9cae-5c90-4ea7-9532-8bc7e5dfb528.png">

Whaler-api provides real-time document editing à la Google docs. Currently, this is used for the notes section of the account details page. It works via an operational transform implementation, where Whaler-api functions as the `server` in the algorithm, and the clients each run a corresponding `client` when they wish to engage in collaborative editing.

The full details of the operational transform algorithm is complex and outside the scope of this README. There are plenty of resources online covering this topic; this one in particular helped me finally understand how it works: https://www.aha.io/blog/text-editor.

Since that's a pretty long article, I'll try to summarize things at a high level. Whaler-api is the OT `server` (`server.go`), which means it maintains the source of truth doc (`doc.go`). This doc is held in memory as a 2D array of `runes`. `Clients` also hold their own copy of the doc, and the goal is to keep all these copies in sync by conveying changes as a series of transformations. Each time a change is made to the doc by a client, that change is packaged up into a series of `Ops` (`ops.go`). These `Ops` define the entire set of changes describing the whole document. For example, if the doc was the text `Hello orld` and a client added the missing `w`, that change would be described by this series of ops: 

1. `retain 6` (H e l l o [space]) - 6 chars
2. `insert "w"`
3. `retain 4` (o r l d) - 4 chars

A set of ops must fully describe the document or they will be rejected by the server. There is one more Op type: delete. If you wanted to delete the space in the last phrase, you would do: 

1. `retain 5`
2. `delete 1`
3. `retain 5` (don't forget we added the "w")

See the `Op` struct defined in `ops.go` for more details on how these operations are represented using an int `N` and a string `S`. 



### Future 👀
#1. The Salesforce integration would be moved to the backend. The frontend would still be responsible for initiating the integration process, but once a Salesforce token is obtained, it would be sent to the backend, where the integration would then be managed. Moving all the Salesforce data management to the backend would free up the app to focus on being a great frontend. As a part of that migration, I would abstract the few direct Salesforce references (e.g. SalesforceId) behind a generic CRM integration interface. Then, adding support for HubSpot, Pipedrive, or any other CRM would be much simpler. 

#2. The current reliance on storing items in memory will present challenges when scaling later. A longer term solution would be to use Redis to store these items. For example, if we store the docs being actively collaborated on in Redis, then we could scale the api to multiple Heroku dynos.

#3. We need pre-production environments, namely a develop and staging environment. This would need to include seperate DBs, as testing migrations is a key reason for the seperation.

#4. This project could use a bit more delineation. Over time, the `models` package would continue to expand and take on too much responsibility. I'd like to further reduce the scope of the existing packages by creating new ones and splitting them up better. Because this project started as a REST api and later switched to GraphQL, there are some leftovers from that transition that could use an update. 

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
