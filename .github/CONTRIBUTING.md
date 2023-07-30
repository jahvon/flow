# Contributing

By participating in this project, you agree to abide our
[code of conduct](https://github.com/jahvon/pilotcli/blob/main/.github/CODE_OF_CONDUCT.md).

## Set up your machine

`pilotcli` is written in [Go](https://golang.org/).

Prerequisites:

- [Go 1.19+](https://golang.org/doc/install)

Other things you might need to run the tests:

- [Docker](https://www.docker.com/)

Clone `pilotcli` anywhere:

```sh
git clone git@github.com:thazelart/pilotcli.git
```

`cd` into the directory and install the dependencies:

```sh
make local/deps
```

## Test your change

You can create a branch for your changes and try to build from the source as you go:

```sh
make go/build
```

When you are satisfied with the changes, we suggest you run:

```sh
make go/test
```

Before you commit the changes, we also suggest you run:

```sh
make pre-commit
```

## Create a commit

Commit messages should be well formatted, and to make that "standardized", we
are using Conventional Commits.

You can follow the documentation on
[their website](https://www.conventionalcommits.org).

## Submit a pull request

Push your branch to your `pilotcli` fork and open a pull request against the main branch.

## Credit

This CONTRIBUTING guideline is very inspired by the [goreleaser](https://github.com/goreleaser/goreleaser/blob/main/CONTRIBUTING.md). Thanks `goreleaser`.
