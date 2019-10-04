# gcov2lcov

Convert golang test coverage to lcov format (which can be uploaded to
coveralls).

See [gcov2lcov-action](https://github.com/jandelgado/gcov2lcov-action) 
for an github action which uses this tool. 

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

```
$ go test -coverprofile=coverage.out && \
gcov2lcov -inputfile=coverage.out -outfile=coverage.lcov
```

## Build and Test

Run `make test` or `make build`.

## Author

Jan Delgado

## License 

MIT
