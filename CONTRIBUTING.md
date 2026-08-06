## Contributing to documentation

The documentation for this SDK lives at [platform.claude.com/docs/en/api/sdks/go](https://platform.claude.com/docs/en/api/sdks/go). To suggest changes, open an issue.

## Setting up the environment

To set up the repository, run:

```sh
$ ./scripts/bootstrap
$ ./scripts/lint
```

This will install all the required dependencies and build the SDK.

You can also [install go 1.22+ manually](https://go.dev/doc/install).

## Modifying/Adding code

Most of the SDK is generated code. Modifications to code will be persisted between generations, but may
result in merge conflicts between manual patches and changes from the generator. The generator will never
modify the contents of the `lib/` and `examples/` directories.

## Spawning external programs

The agent toolset (`tools/agenttoolset`) is the only part of the SDK that starts other programs, and it never
launches one by bare name. A program is either an absolute path written out in the code (`/bin/bash`) or a bare
name resolved through the package's `lookPath` helper ([`tools/agenttoolset/exec.go`](tools/agenttoolset/exec.go)),
which returns an absolute path from `PATH` or an error and never selects a file from the current working
directory — so a binary planted in a repository the tools are pointed at cannot become a helper. When adding
code that runs a program:

- resolve a bare name with `lookPath` and pass the absolute result to `exec.Command`/`exec.CommandContext`;
  don't call `exec.LookPath` directly and don't pass a bare name straight to `exec.Command`;
- treat every lookup error, `exec.ErrDot` included, as "not installed" — never clear it and run the path anyway;
- don't hand-roll a `PATH` walk and don't build `"./" + name`.

`tools/agenttoolset/exec_test.go` and `TestExecGrepNeverRunsPlantedRipgrep` pin this behaviour; the same
guarantee is documented in the other Anthropic SDKs.

## Adding and running examples

All files in the `examples/` directory are not modified by the generator and can be freely edited or added to.

```go
# add an example to examples/<your-example>/main.go

package main

func main() {
  // ...
}
```

```sh
$ go run ./examples/<your-example>
```

## Using the repository from source

To use a local version of this library from source in another project, edit the `go.mod` with a replace
directive. This can be done through the CLI with the following:

```sh
$ go mod edit -replace github.com/anthropics/anthropic-sdk-go=/path/to/anthropic-sdk-go
```

## Running tests

Most tests require you to [set up a mock server](https://github.com/dgellow/steady) against the OpenAPI spec to run the tests.

```sh
$ ./scripts/mock
```

```sh
$ ./scripts/test
```

## Formatting

This library uses the standard gofmt code formatter:

```sh
$ ./scripts/format
```
