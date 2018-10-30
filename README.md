## Tagallery (WIP)

Tagallery is an automated image tagging gallery written as a pet project.

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

TBD
