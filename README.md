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


## Standard Use Case: Web API

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
  package with the `noGod` configuration key. This makes sense if you have got
  multiple `main` packages with different dependencies.

These cases needn't be used and can be overwritten with explicit configuration.


## Usage

You can simply call it with `go run github.com/flowdev/spaghetti-cutter`
from anywhere inside your project.
This will give you an error messages and an exit code bigger than
zero because you didn't configure the `spaghetti-cutter` yet.
This is the way is is preventing spaghetti code.
So it really is more a spaghetti preventer than a spaghetti cutter but
that wouldn't give a nice title nor a nice image.


### Configuration

It is mandatory to use a JSON configuration file `.spaghetti-cutter.json` in
the root directory of your project.
This serves multiple purposes:
- It helps the `spaghetti-cutter` to find the root directory of your project.
- It saves you from retyping command line options again and again.
- It is valuable documentation especially for developers new to the project.

The configuration can have the following elements:
- `tool`, `db` and `god` for tool, database and god packages as discussed above.
- `allowOnlyIn`: for restricting a package to be used only in some packages
  (allow "key" package only in "value" packages).
- `allowAdditionally`: for allowing additional dependencies (for "key" package
  allow additionally "value" packages).
- `size`: the maximum allowed size/complexity of a package. Default is `2048`.
- `noGod`: `main` won't be god package.

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

`*`, `**` and multiple values are allowed for the `tool`, `db`, `god`,
`allowOnlyIn` and `allowAdditionally` keys.

So a full example looks like this:
```json
{
	"allowOnlyIn": {"github.com/lib/pq": ["pkg/model", "pkg/postgres"]},
	"allowAdditonally": {"pkg/shopping": ["pkg/catalogue", "pkg/cart"]},
	"tool": ["pkg/x/*"],
	"db": ["pkg/model", "pkg/postgres"],
	"god": ["cmd/**"],
	"size": 1024
}
```
The `god` line shouldn't be necessary as all packages under `cmd/` should be `main` packages.

The case with multiple executables with different dependencies is interesting, too:
```json
{
	"tool": ["pkg/x/*"],
	"db": ["pkg/model", "pkg/postgres"],
	"allowAdditionally": {
		"cmd/front-end": ["pkg/shopping"],
		"cmd/back-end": ["pkg/catalogue"],
		"pkg/shopping": ["pkg/catalogue", "pkg/cart"]
	},
	"noGod": true,
	"size": 1024
}
```
Here we have got a front-end application for the shopping experience and a
back-end application for updating the catalogue.

### Command line options

The possible command line options are:
```
Usage of spaghetti-cutter:
  -r string
        root directory of the project (shorthand) (default ".")
  -root string
        root directory of the project (default ".")
```

If no `--root` option is given the root directory is found
by crawling up the directory tree starting at the current working directory.
The first directory that contains the configuration file `.spaghetti-cutter.json`
will be taken as project root.


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
