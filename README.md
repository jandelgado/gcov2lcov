# gcov2lcov

Convert golang test coverage to lcov format (which can for example be uploaded
to coveralls).

See [gcov2lcov-action](https://github.com/jandelgado/gcov2lcov-action)
for a github action which uses this tool.

## Credits

This tool is based on [covfmt](https://github.com/ricallinson/covfmt) and
uses some parts of [goveralls](https://github.com/mattn/goveralls).

## Usage

```
Usage of ./gcov2lcov:
  -infile string
    	go coverage file to read, default: <stdin>
  -outfile string
    	lcov file to write, default: <stdout>
```

### Example

```sh
$ go test -coverprofile=coverage.out && \
gcov2lcov -inputfile=coverage.out -outfile=coverage.lcov
```

## Build and Test

* `make test`  to run unit tests
* `make build` to build binary in `bin/` directory
* `make inttest` to run unit integration test

## Author

Jan Delgado

## License

MIT
