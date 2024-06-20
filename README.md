![spaghetti cutter](./spaghetti-cutter.jpg "spaghetti cutter")

# spaghetti-cutter - Win The Fight Against Spaghetti Code

![CircleCI](https://img.shields.io/circleci/build/github/flowdev/spaghetti-cutter/master)
[![Test Coverage](https://api.codeclimate.com/v1/badges/91d98c13ac5390ba6116/test_coverage)](https://codeclimate.com/github/flowdev/spaghetti-cutter/test_coverage)
[![Maintainability](https://api.codeclimate.com/v1/badges/91d98c13ac5390ba6116/maintainability)](https://codeclimate.com/github/flowdev/spaghetti-cutter/maintainability)
[![Go Report Card](https://goreportcard.com/badge/github.com/flowdev/spaghetti-cutter)](https://goreportcard.com/report/github.com/flowdev/spaghetti-cutter)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/flowdev/spaghetti-cutter)
![Twitter URL](https://img.shields.io/twitter/url?style=social&url=https%3A%2F%2Fgithub.com%2Fflowdev%2Fspaghetti-cutter)


## Overview

`spaghetti-cutter` is a command line tool for CI/CD pipelines (and dev machines)
that helps to prevent Go spaghetti code (a.k.a. big ball of mud).

Thankfully in the Go world circular dependencies between packages are already prevented by the compiler.
So this tool has to care only about additional undesired dependencies.

I gave a talk that includes the motivation for this tool and some (old) usage examples:
[![Microservices - The End of Software Design](https://img.youtube.com/vi/ev0dD12bxmg/0.jpg)](https://www.youtube.com/watch?v=ev0dD12bxmg "Microservices - The End of Software Design")

Additionally this tool documents the structure of a project in it's
[configuration](./.spaghetti-cutter.hjson).


## Usage

You can simply call it with `go run github.com/flowdev/spaghetti-cutter@latest`
from anywhere inside your project.

The possible command line options are:
```
Usage of spaghetti-cutter:
  -e    don't report errors and don't exit with an error (shorthand)
  -noerror
        don't report errors and don't exit with an error
  -r string
        root directory of the project (shorthand) (default ".")
  -root string
        root directory of the project (default ".")
```

If no `--root` option is given the root directory is found
by crawling up the directory tree starting at the current working directory.
The first directory that contains the configuration file `.spaghetti-cutter.hjson`
will be taken as project root.

The output looks like:
```
2020/09/10 09:37:08 INFO - configuration 'allowOnlyIn': `github.com/hjson/**`: `x/config` ; `golang.org/x/tools**`: `parse*`, `x/pkgs*`
2020/09/10 09:37:08 INFO - configuration 'allowAdditionally': `*_test`: `parse`
2020/09/10 09:37:08 INFO - configuration 'god': `main`
2020/09/10 09:37:08 INFO - configuration 'tool': `x/*`
2020/09/10 09:37:08 INFO - configuration 'db': ...
2020/09/10 09:37:08 INFO - configuration 'size': 1024
2020/09/10 09:37:08 INFO - configuration 'noGod': false
2020/09/10 09:37:08 INFO - root package: github.com/flowdev/spaghetti-cutter
2020/09/10 09:37:08 INFO - Size of package 'x/config': 699
2020/09/10 09:37:08 INFO - Size of package 'x/pkgs': 134
2020/09/10 09:37:08 INFO - Size of package 'deps': 401
2020/09/10 09:37:08 INFO - Size of package 'parse': 109
2020/09/10 09:37:08 INFO - Size of package 'size': 838
2020/09/10 09:37:08 INFO - Size of package 'x/dirs': 86
2020/09/10 09:37:08 INFO - Size of package '/': 202
2020/09/10 09:37:08 INFO - No errors found.
```

First the configuration values and the root package are reported.
So you can easily ensure that the correct configuration file is taken.

All package sizes are reported and last but not least any violations found.
Since no error was found the return code is 0.

A typical error message would be:
```
2020/09/10 10:31:14 ERROR - domain package 'pkg/shopping' isn't allowed to import package 'pkg/cart'
```

The return code is 1.
From the output you can see that
- the package `pkg/shopping` is recognized as standard domain package,
- it imports the `pkg/cart` package and
- there is no `allowAdditionally` configuration to allow this.

You can fix that by adding a bit of configuration.

Other non-zero return codes are possible for technical problems (unparsable code: 6, ...).
If used properly in the build pipeline a non-zero return code will stop the
build and the problem has to be fixed first.
So undesired imports (spaghetti relationships) are prevented.


## Standard Use Case: Web API

This tool was especially created with Web APIs in mind as that is what about
95% of all Gophers do according to my own completely unscientifical research.

So it offers special handling for the following cases:
- Tools: Tool packages are allowed to be used everywhere else except in other
  tool packages. But they aren't allowed to import any other internal packages.
- Tool sub-packages: Sub-packages of tool packages aren't allowed to import any
  other internal package like tool packages. Additionally they aren't allowed
  to be used anywhere else in the project. So you should use explicit
  configuration with explanations as comments (what the sub-packages contain
  and why they exist at all).
- Database: DB packages are allowed to be used in standard (business) packages.
  Of course they can use tool packages but nothing else.  Domain data
  structures should be in a tool package.
- Database sub-packages: Sub-packages of DB packages are allowed to only import
  tool packages like DB packages. Additionally they aren't allowed to be used
  anywhere else in the project. So you should use explicit configuration with
  explanations as comments (what the sub-packages contain and why they exist at
  all).
- God: A god package can see and use everything. You should use this with great
  care. `main` is the only default god package used if no explicit
  configuration is given. You should only rarely add more.  You can switch
  `main` to a standard package with the `noGod` configuration key. This makes
  sense if you have got multiple `main` packages with different dependencies.

These cases needn't be used and can be overridden with explicit configuration.


## Configuration

It is mandatory to use a HJSON configuration file `.spaghetti-cutter.hjson` in
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
```hjson
{
	"tool": ["x/*"]
}
```
All packages directly under `x` are tool packages that can be used everywhere else in the project.

A slightly different variant is:
```hjson
{
	"tool": ["x/**"]
}
```
All packages under `x` are tool packages that can be used everywhere else in the project.
So the `**` makes all sub-packages tool packages, too.
In most cases one level is enough.

Multiple values are possible for a single key.
So this is another valid configuration file:
```hjson
{
	"tool": ["x/*", "parse"]
}
```

`*`, `**` and multiple values are allowed for the `tool`, `db`, `god`,
`allowOnlyIn` and `allowAdditionally` values.
`*` and `**` are supported for `allowOnlyIn` and `allowAdditionally` keys, too.

So a full example looks like this:
```hjson
{
	"allowOnlyIn": {
		"github.com/lib/pq": ["main"]
		"github.com/jmoiron/sqlx": ["pkg/model", "pkg/postgres"]
	},
	"allowAdditonally": {"pkg/shopping": ["pkg/catalogue", "pkg/cart"]},
	"tool": ["pkg/x/*"],
	"db": ["pkg/model", "pkg/postgres"],
	"god": ["cmd/**"],
	"size": 1024
}
```
The `god` line shouldn't be necessary as all packages under `cmd/` should be `main` packages.

The case with multiple executables with different dependencies is interesting, too:
```hjson
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

Finally you can use variables in the key/value maps:
```hjson
{
	"allowAdditonally": {"pkg/$*/db": ["pkg/$1/model"]},
}
```
In this example there are big "modules" that each have their own database,
model and tool packages.
Cross "module" access can be easily prevented with only one line of
configuration as the `$1` in the value has to be the same as the `$*` in the
key.
You can use `$**` similarly.
Multiple variables are possible and they can be used as `$1` ... `$9` in the
values (`$1` refering to the first `$*` and `$9` to the ninth `$**` in the key).
The maximum of `9` should be big enough for really complex projects and
helps to find configuration errors.


## Installation

Of course you can just head over to the
[latest release](https://github.com/flowdev/spaghetti-cutter/releases/latest)
and grab a pre-built binary for your OS.
But that is difficult to keep in sync when collaborating with others in a team.

A much better way is to just do: `go run github.com/flowdev/spaghetti-cutter@latest`
If you use an explicit version and use it in multiple places the next approach is better.

A great approach for big projects goes this way:

First include the latest version in your `go.mod` file, e.g.:
```
require (
	github.com/flowdev/spaghetti-cutter v0.9.9
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

This ensures that the package is indeed fetched and built but not included in
the main or test executables. This is the
[canonical workaround](https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module)
to keep everything in sync and lean.
Here is a [talk by Robert Radestock](https://www.youtube.com/watch?v=PhBhwgYFuw0)
about this topic.

Finally you can run `go mod vendor` if that is what you like.
