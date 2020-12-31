## Tagallery

[![Build Status](https://travis-ci.org/catyphram/tagallery.svg?branch=master)](https://travis-ci.org/catyphram/tagallery)
[![codecov](https://codecov.io/gh/catyphram/tagallery/branch/master/graph/badge.svg)](https://codecov.io/gh/catyphram/tagallery)
[![Commitizen friendly](https://img.shields.io/badge/commitizen-friendly-brightgreen.svg)](http://commitizen.github.io/cz-cli/)


Tagallery is an automated image tagging gallery where one can categorize images. To speed up this handy work, a neural network will try to figure out the categories of an image by comparing it with already categorized ones and proposes these categories for validation.

**This project is Work In Progress.**

### Inspiration

The main goal of this pet project is to try out new programming languages, frameworks and technologies in general. Particular interest lies in [Go](https://golang.org/), [MongoDB](https://www.mongodb.com/), [Tensorflow](https://www.tensorflow.org/), [Python](https://www.python.org/), [Nuxt.js](https://nuxtjs.org/), [Vue.js](https://vuejs.org/), as well as test automation, versioning, releasing and the setup of a Continuous Integration/Delivery pipeline utilizing [Docker](https://www.docker.com/) container.

### Setup

The application can easily be started with docker-compose. Simply run `docker-compose up -d` to start, or `docker-compose down` to stop the application. While the `./image` directory is automatically created you may need to change it's permissions/owner since the docker container uses the root user, e. g. `sudo chown -R $(id -u):$(id -g) images`. Please visit the [documentation](https://docs.docker.com/compose/install/) for more information and instructions on how to install Docker and docker-compose.  
You can now fill your `./images/unprocessed` directory with your images and start categorizing them.
For more detailed instruction on how to separately start the services, see section [API](#api) and [Client](#client) below.

### API

The REST-API is written in [Go](https://golang.org/) and connects to a [MongoDB](https://www.mongodb.com/) database.

#### Configuration

No configuration is needed.  
You can however override the default options via environment variables:
- `DATABASE_HOST=localhost:27017`
- `DATABASE=tagallery`
- `DEBUG=false`
- `PORT=3333`
- `IMAGES=./images`

#### Compilation

Go into the `api/` folder and run `go get -d ./...` to download the dependencies, followed by `go build` to compile the executable.
This project is using [Go Modules](https://github.com/golang/go/wiki/Modules), which were introduced with Go v1.11.

#### Execution

Once built the executable can be started with `./api`, or, if using env variables, `PORT=3333 DATABASE=tagallery DATABASE_HOST=localhost:27017 ./api`.

#### Testing

Run `go tool vet .` to lint the code and `go test ./...` to test all the packages.
  
### Client

The client is a PWA rendered server side by using [Nuxt.js](https://nuxtjs.org/) - which is built upon [Vue.js](https://vuejs.org/) - with [TypeScript](https://www.typescriptlang.org/), [SCSS](https://sass-lang.com/) and [Jest](https://jestjs.io/) and. [Vuex](https://vuex.vuejs.org/) is used for state management and [Vuetify](https://vuetifyjs.com/) - which is based on [Material Design](https://material.io/design/) - for layout and styling.

Seek the [Nuxt Install ducumentation](https://nuxtjs.org/docs/2.x/get-started/installation) for more information on how to serve, build and test the application via the command line.

TLDR: Use `yarn dev` to start a dev server, `yarn lint|test` for linting/testing and `yarn build && yarn start` to serve a production build.

### Contribution

Refer to the [API](#testing) and [Client](#client) section for testing and linting instructions.

Commit messages follow the [Conventional Commits specification](https://www.conventionalcommits.org/) and are enforeced by [Commitlint](https://conventional-changelog.github.io/commitlint/#/) using the [conventional config](https://github.com/conventional-changelog/commitlint/tree/master/%40commitlint/config-conventional#type-enum).  
Keywords may be added to [close a ticket](https://help.github.com/articles/closing-issues-using-keywords/) or otherwise state the progress, e. g. `progresses #1`.  
The repository supports [Commitizen](http://commitizen.github.io/cz-cli/), so you may use `yarn cz` as an alternativ to `git commit` if you commit on the command line. Use `\n` to separate a multi-line body message.

We are following a [feature branching model](https://guides.github.com/introduction/flow/). Changes to the `master` branch have to be made through PRs and need to pass the [CI pipeline](https://travis-ci.com/), which runs linters and tests, before it can be merged via Github.
