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

Thankfully in the Go world circular dependencies between packages are already prevented by the compiler.
So this tool has to care only about additional undesired dependencies.

## Installation

Of course you can just head over to the
[latest release](https://github.com/flowdev/spaghetti-cutter/releases/latest)
and grab a pre-built binary and change the extension for your OS.
But that is difficult to keep in sync when collaborating with others in a team.

A much better approach for teams goes this way:

First include the latest version in your `go.mod` file, e.g.:
```Go
require (
	github.com/flowdev/spaghetti-cutter v0.9
)
```

Now add a file like the following to your main package.

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
  care. `main` is the only default god package used if no explicit package is
  given. You should only rarely add more.  You can switch `main` to a standard
  package with the `no-god` switch. This makes sense if you have got multiple
  `main` packages with different dependencies.

Any of these rules can be overwritten with an explicit `allow` directive.


### Configuration

It is strongly recommended to use a JSON configuration file
`.spaghetti-cutter.json` in the root directory of your project.
This serves multiple purposes:
- It helps the `spaghetti-cutter` to find the root directory of your project.
- It saves you from retyping command line options again and again.
- It documents the structure within the project.

The configuration can have the following elements:
- `tool`, `db` and `god` for tool, database and god packages as discussed above.
- `allow`: for allowing additional dependencies.
- `size`: the maximum allowed size/complexity of a package. Default is `2048`.
- `no-god`: `main` won't be god package.
- `ignore-vendor`: ignore vendor directories when searching for the project root
- `root`: explicit project root. Should be given by the position of the config file instead.
  (only makes sense as a command line argument).

The size configuration key prevents a clever developer from just thowing all of
the spaghetti code into a single package.
With the `spaghetti-cutter` such things will become obvious and you can put
them as technical dept into your back log.

This is a simple example configuration file:
```json
{
	"tool": "x/*"
}
```
All packages directly under `x` are tool packages that can be used everywhere else in the project.

A slightly different variant is:
```json
{
	"tool": "x/**"
}
```
All packages under `x` are tool packages that can be used everywhere else in the project.
So the `**` makes all sub-packages tool packages, too.
In most cases one level is enough.

Multiple values are possible for a single key.
So this is another valid configuration file:
```json
{
	"tool": ["x/*", "parse"]
}
```

`*`, `**` and multiple values are allowed for the `tool`, `db`, `god` and `allow` keys.

So a rather complex example looks like this:
```json
{
	"tool": "pkg/x/*",
	"db": ["pkg/model", "pkg/postgres"],
	"allow": ["pkg/shopping pkg/catalogue", "pkg/shopping pkg/cart"],
	"god": "cmd/**",
	"size": 1024
}
```
The `god` line shouldn't be necessary as all packages under `cmd/` should be `main` packages.

The case with multiple executables with different dependencies is interesting, too:
```json
{
	"tool": "pkg/x/*",
	"db": ["pkg/model", "pkg/postgres"],
	"allow": [
		"cmd/front-end pkg/shopping",
		"cmd/back-end pkg/catalogue",
		"pkg/shopping pkg/catalogue",
		"pkg/shopping pkg/cart"
	],
	"no-god": true,
	"size": 1024
}
```
Here we have got a front-end application for the shopping experience and a
back-end application for updating the catalogue.

### Command line options
Generally it is possible to override any configuration key with command line options.

It is most useful to use these options on the command line:
- `--ignore-vendor`: ignore vendor directories when searching for the project root.
- `--root`: explicit project root. Should be given by the position of the config file instead.
  So it only makes sense if you don't have got any configuration file at all or
  you have to override a misplaced `go.mod` file.

Unfortunately it is currently impossible to use the `--root` option to find the
correct configuration file because the config file is searched before command
line options are read.
Instead the directory tree is crawled upward starting at the current working
directory in order to find the `.spaghetti-cutter.json` file.
I would be willing to change this if there is demand for it.

If no `--root` option or `root` config key is given the root directory is found
by crawling up the directory tree starting at the current working directory.
The first directory that contains
- the configuration file `.spaghetti-cutter.json` or
- the Go module file `go.mod` or
- a vendor directory `vendor` (if not prevented by `--ignore-vendor`)

will be taken as project root.

The possible command line options are:
```
Usage of spaghetti-cutter:
  -allow value
     allowed package dependency (e.g. 'pkg/a/uses pkg/x/util')
  -db value
     common domain/database package (can only depend on tools) (e.g. 'pkg/*/db'; '*' matches anything except for a '/')
  -god value
     god package that can see everything (default: 'main')
  -ignore-vendor
     ignore any 'vendor' directory when searching the project root
  -no-god
     override default: 'main' won't be implicit god package
  -root string
     project root directory
  -size uint
     maximum size of a package in "lines" (default 2048)
  -tool value
     tool package (leave package) (e.g. 'pkg/x/**'; '**' matches anything including a '/')
```
They work exactly like the similar configuration keys and overwrite them.


## Best Practices

For web APIs it is useful to split into independent business packages at router level.
Router itself can be in a central (god) package. The split can be done in two ways:
1. The central router calls handlers that reside in the business packages.
1. The central router composes itself from subrouters in business packages.

The second option minimizes the API surface of the business package and helps
to ensure that all routes handled by a business package share a common URL path root.
But it also adds the concept of subrouters that isn't used so widely and
increasing cognitive load this way.
Plus it makes it harder to find the implementation for a route.

So I recommend to start with the first option and switch to the second if the
central router becomes too big to handle.


### Criteria For When To Split A Service

A common reason to split a service is when different parts of the service have
to scale very differently.
A shop front-end that has to serve many thousand customers needs to scale much
more than the shop back-end that only has to serve a few employees.
Of course this isn't useful as long as you have got only a single instance of
the shop front-end running.
Please remember that Go is often used to consolidate many servers written in
some script language. Often replacing ten script servers with a single Go
instance saving a lot of hosting costs and operational work.

Another good reason to split a service is when the data the different parts of
the service work on is very or even completely different.
Please remember that overlaping data will lead to redundancies and you have to
ensure consistency on your own.
After such a split the overall system is usually only eventual consistent.

The last and weakest indicator is that the service is growing unbounded like cancer.
It is completely normal that a service is growing.
When the tests run for too long it is better to find a tool for handling
monorepos that helps you to run only the necessary tests.
Unfortunately I can't point you to one. But I know that this is a problem that
has been solved multiple times.
Such additional tools can go a long way before it makes sense to split a service.


### Recommendation How To Split A Service If Really Useful

I recommend to split a service if it is sure to be really useful by first
looking at the package structure in `.spaghetti-cutter.json`.
You wouldn't want to separate packages into own services if one package depends
on the other as per `allow` directive.
`tool` packages are a bit simpler since they tend change less often.
They should just serve a single purpose well.
So they can be easily extracted into an external library or they can be copied
if the reusage isn't driven by businss needs but more accidental.
If some `tool` packages won't be used by all the services after the split you
should take advantage of that.

Next it is important to look at the DB usage. Packages that share a lot of data
or methods to access data should not be split into separate services.
Now it is time to find the weakest link between future services.
You should consider all three types of links:
- `tool` packages (least important),
- database usage (quite important) and
- `allow` directives (most important).

When the weakes links are found it is great if you can even minimize these links.
Over time they tend to accumulate and some aren't really necessary anymore.
It is often great to get a perspective from the business side about this.

Now it is time to replace the remaining internal calls between packages that
will become separate services with external calls.
RESTful HTTP requests and gRPC are used most often for this.
Messaging between services can give you more scalability and decoupling but is
harder to debug and some additional technology to master.
Often you already have got a company wide standard for communication between services.

Creating multiple main packages for the different services should be rather simple.
Each should just be a subset of the old main package.
If you have got more god packages than just `main` you should split them now of course.
Now you already have multiple separately deployable services.

Finally you can do the split.

You can minimize the necessary work a lot by always watching dependencies grow
and minimizing links as soon as possible.
The `spaghetti-cutter` can be your companion on the way.
