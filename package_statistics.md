# Package Statistics

| package | type | direct deps | all deps | users | max score | min score |
| :- | :-: | -: | -: | -: | -: | -: |
| [/](#root-package) | [ \[G\] ](#legend) | [9](#direct-dependencies-imports-of-root-package) | [9](#all-including-transitive-dependencies-imports-of-root-package) | 0 | 0 | 0 |
| [deps](#package-deps) | [ \[S\] ](#legend) | [3](#direct-dependencies-imports-of-package-deps) | [3](#all-including-transitive-dependencies-imports-of-package-deps) | [1](#packages-using-importing-package-deps) | 0 | 0 |
| [doc](#package-doc) | [ \[S\] ](#legend) | [1](#direct-dependencies-imports-of-package-doc) | [1](#all-including-transitive-dependencies-imports-of-package-doc) | [1](#packages-using-importing-package-doc) | 0 | 0 |
| [parse](#package-parse) | [ \[S\] ](#legend) | [1](#direct-dependencies-imports-of-package-parse) | [1](#all-including-transitive-dependencies-imports-of-package-parse) | [1](#packages-using-importing-package-parse) | 0 | 0 |
| [size](#package-size) | [ \[S\] ](#legend) | [1](#direct-dependencies-imports-of-package-size) | [1](#all-including-transitive-dependencies-imports-of-package-size) | [1](#packages-using-importing-package-size) | 0 | 0 |
| [stat](#package-stat) | [ \[S\] ](#legend) | [1](#direct-dependencies-imports-of-package-stat) | [1](#all-including-transitive-dependencies-imports-of-package-stat) | [1](#packages-using-importing-package-stat) | 0 | 0 |
| [x/config](#package-xconfig) | [ \[T\] ](#legend) | [1](#direct-dependencies-imports-of-package-xconfig) | [1](#all-including-transitive-dependencies-imports-of-package-xconfig) | [2](#packages-using-importing-package-xconfig) | 0 | 0 |

### Legend

* package - name of the internal package without the part common to all packages.
* type - type of the package:
  * [G] - God package (can use all packages)
  * [D] - Database package (can only use tool and other database packages)
  * [T] - Tool package (foundational, no dependencies)
  * [S] - Standard package (can only use tool and database packages)
* direct deps - number of internal packages directly imported by this one.
* all deps - number of transitive internal packages imported by this package.
* users - number of internal packages that import this one.
* max score - sum of the numbers of packages hidden from user packages.
* min score - number of packages hidden from all user packages combined.


### Root Package


#### Direct Dependencies (Imports) Of Root Package
`data`, [deps](#package-deps), [doc](#package-doc), [parse](#package-parse), [size](#package-size), [stat](#package-stat), [x/config](#package-xconfig), `x/dirs`, `x/pkgs`

#### All (Including Transitive) Dependencies (Imports) Of Root Package
`data`, [deps](#package-deps), [doc](#package-doc), [parse](#package-parse), [size](#package-size), [stat](#package-stat), [x/config](#package-xconfig), `x/dirs`, `x/pkgs`

### Package deps


#### Direct Dependencies (Imports) Of Package deps
`data`, [x/config](#package-xconfig), `x/pkgs`

#### All (Including Transitive) Dependencies (Imports) Of Package deps
`data`, [x/config](#package-xconfig), `x/pkgs`

#### Packages Using (Importing) Package deps
[root](#root-package)

### Package doc


#### Direct Dependencies (Imports) Of Package doc
`data`

#### All (Including Transitive) Dependencies (Imports) Of Package doc
`data`

#### Packages Using (Importing) Package doc
[root](#root-package)

### Package parse


#### Direct Dependencies (Imports) Of Package parse
`x/pkgs`

#### All (Including Transitive) Dependencies (Imports) Of Package parse
`x/pkgs`

#### Packages Using (Importing) Package parse
[root](#root-package)

### Package size


#### Direct Dependencies (Imports) Of Package size
`x/pkgs`

#### All (Including Transitive) Dependencies (Imports) Of Package size
`x/pkgs`

#### Packages Using (Importing) Package size
[root](#root-package)

### Package stat


#### Direct Dependencies (Imports) Of Package stat
`data`

#### All (Including Transitive) Dependencies (Imports) Of Package stat
`data`

#### Packages Using (Importing) Package stat
[root](#root-package)

### Package x/config


#### Direct Dependencies (Imports) Of Package x/config
`data`

#### All (Including Transitive) Dependencies (Imports) Of Package x/config
`data`

#### Packages Using (Importing) Package x/config
[root](#root-package), [deps](#package-deps)
