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

Include latest version in `go.mod` file.

Add a file like the following to a main package.
Or add the import line to an existing file with similar build comment.
This ensures that the package is indeed fetched and build but not included in
the main or test executables and is the
[canonical workaround](https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module)
for this kind of problem.
Here is a complete [talk by Robert Radestock](https://www.youtube.com/watch?v=PhBhwgYFuw0) about this topic.

```Go
//+build tools

package main

import (
    _ "github.com/flowdev/spaghetti-cutter"
)
```

## Usage

Description of (simple) usage

It is strongly recommended to use a JSON configuration file
`.spaghetti-cutter.json` in the root directory of your repository.
The configuration file is explained in detail below.

- For typical web API

- Hence:
  - tools
  - DB
  - god
  - allow
  - domain data structures can be either DB or tools

Best practices:
- Split into independent business packages at router level
  1. Router itself can be in central (god) package with
     handlers called by the router in the business packages.
  1. You can use subrouters in business packages and
     compose them in the central router.

### Configuration Examples

Configuration file: syntax with examples


### Command line options

Details:
- How the project root is found


## Best Practices

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
