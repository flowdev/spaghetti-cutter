# Package Statistics

Start package - /

| package | type | direct deps | all deps | users | max score | min score |
| :- | :-: | -: | -: | -: | -: | -: |
| / | [G] | 9 | 9 | 0 | 0 | 0 |
| deps | [S] | 3 | 3 | 1 | 0 | 0 |
| doc | [S] | 1 | 1 | 1 | 0 | 0 |
| parse | [S] | 1 | 1 | 1 | 0 | 0 |
| size | [S] | 1 | 1 | 1 | 0 | 0 |
| stat | [S] | 1 | 1 | 1 | 0 | 0 |

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
`data`, [deps](#package-deps), [doc](#package-doc), [parse](#package-parse), [size](#package-size), [stat](#package-stat), `x/config`, `x/dirs`, `x/pkgs`

#### All (Including Transitive) Dependencies (Imports) Of Root Package
`data`, [deps](#package-deps), [doc](#package-doc), [parse](#package-parse), [size](#package-size), [stat](#package-stat), `x/config`, `x/dirs`, `x/pkgs`

#### Packages Using (Importing) Root Package


#### Packages Not Imported By Users Of Root Package


#### Packages Not Imported By Any Users Of Root Package


### Package deps

#### Direct Dependencies (Imports) Of Package deps
`data`, `x/config`, `x/pkgs`

#### All (Including Transitive) Dependencies (Imports) Of Package deps
`data`, `x/config`, `x/pkgs`

#### Packages Using (Importing) Package deps
[/](#root-package)

#### Packages Not Imported By Users Of Package deps
* #root-package: 


#### Packages Not Imported By Any Users Of Package deps


### Package doc

#### Direct Dependencies (Imports) Of Package doc
`data`

#### All (Including Transitive) Dependencies (Imports) Of Package doc
`data`

#### Packages Using (Importing) Package doc
[/](#root-package)

#### Packages Not Imported By Users Of Package doc
* #root-package: 


#### Packages Not Imported By Any Users Of Package doc


### Package parse

#### Direct Dependencies (Imports) Of Package parse
`x/pkgs`

#### All (Including Transitive) Dependencies (Imports) Of Package parse
`x/pkgs`

#### Packages Using (Importing) Package parse
[/](#root-package)

#### Packages Not Imported By Users Of Package parse
* #root-package: 


#### Packages Not Imported By Any Users Of Package parse


### Package size

#### Direct Dependencies (Imports) Of Package size
`x/pkgs`

#### All (Including Transitive) Dependencies (Imports) Of Package size
`x/pkgs`

#### Packages Using (Importing) Package size
[/](#root-package)

#### Packages Not Imported By Users Of Package size
* #root-package: 


#### Packages Not Imported By Any Users Of Package size


### Package stat

#### Direct Dependencies (Imports) Of Package stat
`data`

#### All (Including Transitive) Dependencies (Imports) Of Package stat
`data`

#### Packages Using (Importing) Package stat
[/](#root-package)

#### Packages Not Imported By Users Of Package stat
* #root-package: 


#### Packages Not Imported By Any Users Of Package stat
