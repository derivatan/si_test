# si_test

This repo contains the integration tests for [si](https://github.com/derivatan/si) which is a simple ORM for golang.

The tests are put here, in its own repository, because I don't want the library itself to import packages that are only needed for the testing.

## Running the tests.

First time:

`make init`

And then to run the tests:

`make integration`
