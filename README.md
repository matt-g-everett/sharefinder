# sharefinder

Command line tool that reads a JSON doc containing funds and shares and extracts the shares for a given fund

When run, sharefinder always reads a file from testdata/example.json and outputs the shares that were found
in a JSON array to stdout.

## Structure

The structure of the application packages: -

- main: the application entry point
- api: defines and handles the input JSON format
- model: defines and creates the DAG (Directed Acyclic Graph) of holdings (funds or shares)
- finder: provides functionality for finding shares in the model

## Running

The application can be run with: -

    go run main.go

## Testing

The solution contains unit tests that can be run with

    go test
