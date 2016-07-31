# TLDR;

A simple Makefile with some useful rules for go 1.6+ development.

# Long read

It's for 1.6+ because takes into account that `vendor/` is the dependencies directory, and tries to not execute tests and linters of vendors.

Also has two types of rules for tests and linters: the normal ones and the `-ci` ones, which are for running within Continuous Integration environements, where we want to stop executing the build as the first thing fails. The latter use [fgt](https://github.com/GeertJohan/fgt) is a little go utility needed because not all the linters return the same status code and output when finding issues ([source](https://github.com/golang/lint/issues/65)).

# Rule types:

- testing: launch tests excluding vendor packages.
- linters
- code coverage, with [go-carpet](https://github.com/msoap/go-carpet)
- dependency management with [godep](https://github.com/tools/godep)

# Linters

- fmt
- imports
- vet, without alerting when not using [composite literals](https://golang.org/cmd/vet/#hdr-Unkeyed_composite_literals).
- lint
- errcheck, without alerting if `Close` returned errors are not checked. Otherwise you would need to wrap the `Close()` inside a closure when deferring it ([source](https://github.com/kisielk/errcheck/issues/101)).

If you need more linters check out [gometalinter](https://github.com/alecthomas/gometalinter).
