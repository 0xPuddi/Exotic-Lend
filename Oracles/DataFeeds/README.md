# Data Feeds Server
This repo consists of an HTTP server that collects, cleans and arranges price data, to then provide endpoints with updated data, cured data, simple data and more complex informations.

## Architecture
We define struct types that can have very basic tags for Postgre database operation:

- `db`: You can define the data type to be stored inside the database
- `rel`: You can define any relation on that field
- `idx`: You can define if the field needs an index

## Usage
To test make sure to create a `.env` file and add all necessary environment variables, see `.env.example` and `Necessary testing variables`

Then run:
```sh
# To run tests
make test 
# To run test with coverage output
make test-cover
# To run test with coverage output and html view
make test-cover-view
```

To run the executable make sure to create a `.env` file and add all necessary environment variables, see `.env.example` and `Necessary environment variables`

Then run:
```sh
# Build
make build
./bin/DataFeedExec
# Build and run
make
```

### Notes
wip