# sharefinder

Command line tool that reads a JSON doc containing funds and shares and extracts the shares for a given fund.

When run, sharefinder always reads a file from testdata/example.json and outputs the shares that were found
in a JSON array to stdout.

Some of the more important comments in the code that describe assumptions, decisions and future possibilities
have a `NOTE` marker on them.

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

The solution contains unit tests that can be run with: -

    go test

## Benchmarking and profiling

The solution contains an extra finder algorithm that uses a memento pattern. This can be run through benchmarks with: -

    go test -bench . -cpuprofile cpu.prof -count 5

To view this in pprof, try: -

    go tool pprof -http=":3000" cpu.prof
