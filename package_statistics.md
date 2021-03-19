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

#### Direct Dependencies (Imports) For Root Package

`data`, [deps](#package-deps), [doc](#package-doc), [parse](#package-parse), [size](#package-size), [stat](#package-stat), `x/config`, `x/dirs`, `x/pkgs`

#### All Dependencies (Imports) Including Transitive Dependencies

#### Packages Using (Importing) This Package

#### Packages Not Imported By Users

#### Packages Not Imported By All Users Combined


### Package deps

#### Direct Dependencies (Imports) For Package deps

`data`, `x/config`, `x/pkgs`

#### All Dependencies (Imports) Including Transitive Dependencies

#### Packages Using (Importing) This Package

#### Packages Not Imported By Users

#### Packages Not Imported By All Users Combined


### Package doc

#### Direct Dependencies (Imports) For Package doc

`data`

#### All Dependencies (Imports) Including Transitive Dependencies

#### Packages Using (Importing) This Package

#### Packages Not Imported By Users

#### Packages Not Imported By All Users Combined


### Package parse

#### Direct Dependencies (Imports) For Package parse

`x/pkgs`

#### All Dependencies (Imports) Including Transitive Dependencies

#### Packages Using (Importing) This Package

#### Packages Not Imported By Users

#### Packages Not Imported By All Users Combined


### Package size

#### Direct Dependencies (Imports) For Package size

`x/pkgs`

#### All Dependencies (Imports) Including Transitive Dependencies

#### Packages Using (Importing) This Package

#### Packages Not Imported By Users

#### Packages Not Imported By All Users Combined


### Package stat

#### Direct Dependencies (Imports) For Package stat

`data`

#### All Dependencies (Imports) Including Transitive Dependencies

#### Packages Using (Importing) This Package

#### Packages Not Imported By Users

#### Packages Not Imported By All Users Combined
