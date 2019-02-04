## Tagallery

[![Build Status](https://travis-ci.org/catyphram/tagallery.svg?branch=master)](https://travis-ci.org/catyphram/tagallery)

Tagallery is an automated image tagging gallery where one can categorize images. To speed up this handy work, a neural network will try to figure out the categories of an image by comparing it with already categorized ones and proposes these cateogries for validation.

**This project is Work In Progress.**

### Inspiration

The main goal of this pet project is to try out new programming languages, frameworks and technologies in general. Particular interest lies in [Go](https://golang.org/), [MongoDB](https://www.mongodb.com/), [Tensorflow](https://www.tensorflow.org/), [Python](https://www.python.org/), [Vue.js](https://vuejs.org/), as well as test automation, versioning, releasing and the setup of a Continuous Integration/Delivery pipeline utilizing [Docker](https://www.docker.com/) container.

### Setup

The application can easily be started with docker-compose. Simply run `docker-compose up -d` to start, or `docker-compose down` to stop the application. Please visit the [documentation](https://docs.docker.com/compose/install/) for more information and instructions on how to install Docker and docker-compose.  
For more detailed instruction on how to separately start the services, see section [API](#api) and [Client](#client) below.

### API

The REST-API is written in [Go](https://golang.org/) and connects to a [MongoDB](https://www.mongodb.com/) database.

#### Configuration

A MongoDB server has to be started manually beforehand.  
Copy `api/src/config_example.json` and configure your settings.  
The default location of the config file is at `<workdir>/config.json` - `<workdir>` being the folder in which the binary is executed - and can be overridden via the command line argument `-c` or `--config`.  
The config params can also be set via environment variables.

#### Compilation

Go into the `api/` folder and run `go get -d ./...` to download the dependencies, followed by `go build` to compile the executable.  
This project is using [Go Modules](https://github.com/golang/go/wiki/Modules), which were introduced with Go v1.11. So, make sure that the project code is outside your `$GOPATH/src`, or set the environment variable `GO111MODULE=on` when running `go get|build|...`.

#### Execution

Once built the executable can be started with `./api -c=/path/to/config.json`, or, if using env variables, `UNPROCESSED_IMAGES=/path/to/image/dir PORT=3333 DATABASE=tagallery DATABASE_HOST=localhost:27017 ./api`.

#### Testing

Run `UNPROCESSED_IMAGES=/path/to/image/dir PORT=3333 DATABASE=tagallery DATABASE_HOST=localhost:27017 go test ./...` to test all the packages.
  
When testing multiple packages at once the config options HAVE to be set via environment variables.  
A config file in the same directory won't work since `go test` changes the working directory for each test to the directory of the to-be-tested package, and specifying a path via `--config /path/to/config.json` will cause [problems](https://stackoverflow.com/a/49927684) due to the flag not being supported by all packages.

### Client

The client is written using [Vue.js](https://vuejs.org/) with [TypeScript](https://www.typescriptlang.org/), [SCSS](https://sass-lang.com/), [Jest](https://jestjs.io/) and [Cypress](https://www.cypress.io/). [Vuex](https://vuex.vuejs.org/) is used for state management and [Vue Material](https://vuematerial.io/) - which is based on [Material Design](https://material.io/design/) - for layout and styling.

Seek the [Vue Cli ducumentation](https://cli.vuejs.org/guide/cli-service.html) for more information on how to serve, build and test the application via the command line or the Vue Cli UI.

As a quick reference: Start the Vue UI with `vue ui`, or use the scripts `yarn run serve|build|lint|test:e2e|test:unit`.