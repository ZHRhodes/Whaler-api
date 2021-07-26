
# Whaler-api

Whaler-api is the backend to [Whaler](https://github.com/ZHRhodes/Whaler). 

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

Whaler-api uses GraphQL for all data operations. The full schema is available in this directory as `schema.graphqls`. Whaler-api uses a Go module called `gqlgen` to generate the GraphQL-related boilerplate. This lets us focus on updating the schema first and foremost.


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
