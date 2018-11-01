## Tagallery

Tagallery is an automated image tagging gallery where one can categorize images. To speed up this handy work, a neural network will try to figure out the categories of an image by comparing it with already categorized ones and proposes these cateogries for validation.

**This project is work in progress.**

### Inspiration

The main goal of this pet project is to try out new programming languages, frameworks and technologies in general. Particular interest lies in Go, MongoDb, Tensorflow, Python, Vue.js, as well as test automation, versioning and CI/CD.

### Setup

For now the API and Client have to be built and started separately. See section [API](#api) and [Client](#client) below.  
Dockerfiles and a docker-compose config file will soon be added to allow for an easier setup.

### API

#### Configuration

Run `docker-compose up -d` to start the mongo database.  
Then copy `api/src/config_example.json` and configure your settings.  
The default location of the config file is at `<workdir>/config.json` and can be overridden via command line arguments, e. g. `--config <filePath>`. Config params can also be set via environment variables.

#### Compilation

Go into the `api/src/` folder and run `go build` to compile the executable.

#### Running

Start the api via `./api -c=/path/to/config.json`.

#### Testing

Run `UNPROCESSED_IMAGES=/path/to/empty/dir PORT=3333 DATABASE=tagallery DATABASE_HOST=localhost:27017 /usr/local/go/bin/go test ./...` to test all the packages. The config options for tests can only be set via env variables, not the command line.

### Client

The client will be written using Vue.js with TypeScript, SCSS, Jest and Cypress. Seek the [Vue Cli ducumentation](https://cli.vuejs.org/) for information on how to serve, build and test the application.