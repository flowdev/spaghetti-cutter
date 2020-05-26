![spaghetti cutter](./spaghetti-cutter.jpg "spaghetti cutter")

# spaghetti-cutter - Win The Fight Against Spaghetti Code

![CircleCI](https://img.shields.io/circleci/build/github/flowdev/spaghetti-cutter/master)
[![Test Coverage](https://api.codeclimate.com/v1/badges/91d98c13ac5390ba6116/test_coverage)](https://codeclimate.com/github/flowdev/spaghetti-cutter/test_coverage)
[![Maintainability](https://api.codeclimate.com/v1/badges/91d98c13ac5390ba6116/maintainability)](https://codeclimate.com/github/flowdev/spaghetti-cutter/maintainability)
[![Go Report Card](https://goreportcard.com/badge/github.com/flowdev/spaghetti-cutter)](https://goreportcard.com/report/github.com/flowdev/spaghetti-cutter)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/flowdev/spaghetti-cutter)
![Twitter URL](https://img.shields.io/twitter/url?style=social&url=https%3A%2F%2Fgithub.com%2Fflowdev%2Fspaghetti-cutter)

`spaghetti-cutter` is a command line tool for CI/CD pipelines (and dev machines)
that helps to cut Go spaghetti code (a.k.a. big ball of mud) into manageable pieces
and keep it that way.

## Installation

First include the latest version in your `go.mod` file, e.g.:
```Go
require (
	github.com/flowdev/spaghetti-cutter v0.9
)
```

Now add a file like the following to a main package.

```Go
//+build tools

package main

import (
    _ "github.com/flowdev/spaghetti-cutter"
)
```

Or add the import line to an existing file with similar build comment.
This ensures that the package is indeed fetched and built but not included in
the main or test executables. This is the
[canonical workaround](https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module)
to keep everything in sync and lean.
Here is a [talk by Robert Radestock](https://www.youtube.com/watch?v=PhBhwgYFuw0)
about this topic.

Finally you can run `go mod vendor` if that is what you like.


## Usage

You can simply call it with `go run github.com/flowdev/spaghetti-cutter`
from anywhere inside your project.
This will most likely give you some error messages and an exit code bigger than
zero because you didn't configure the `spaghetti-cutter` yet.


### Standard Use Case: Web API

This tool was especially created with Web APIs in mind as that is what about
95% of all Gophers do according to my own totally unscientifical research.

So it offers special handling for the following cases:
- Tools: Tool packages are allowed to be used everywhere else except in other
  tool packages. But subpackages of a tool package are allowed to be used by
  the parent tool package.
- Database: DB packages are allowed to be used in other DB packages and
  standard (business) packages. Of course they can use tool packages.
  Domain data structures can be either `db` or `tool` packages.
- God: A god package can see and use everything. You should use this with great
  care. `main` is the only default god package. You should only rarely add more.
  You can switch `main` to a standard package. This makes sense if you have got
  multiple `main` packages.

Any of these rules can be overwritten with an explicit `allow` directive.


### Configuration

It is strongly recommended to use a JSON configuration file
`.spaghetti-cutter.json` in the root directory of your project.
This serves multiple purposes:
- It helps the `spaghetti-cutter` to find the root directory of your project.
- It saves you from retyping command line options again and again.
- It documents the structure within the project.

The configuration file is explained in detail below.


Configuration file: syntax with examples


### Command line options

Details:
- How the project root is found


## Best Practices

- Split into independent business packages at router level
  1. Router itself can be in central (god) package with
     handlers called by the router in the business packages.
  1. You can use subrouters in business packages and
     compose them in the central router.


### Criteria For When To Split A Service

- When different parts of the service have to scale very differently
  (e.g. front-end vs. back-end of a shop).
- The data the different parts of the service work on is very or even completely different.
- Last and weakest indicator: A service is growing unbounded like cancer.

### Recommendation How To Split A Service If Really Useful

1. Look at the structure (allowed dependencies)
1. Look at DB usage
1. Find spot of "weakest link"
1. Try to minimize links (but not artificially)
1. Replace remaining internal calls with external (e.g. HTTP) calls or messages.
1. Split.
