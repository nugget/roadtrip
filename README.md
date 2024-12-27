[![Go](https://github.com/nugget/roadtrip/actions/workflows/go.yml/badge.svg)](https://github.com/nugget/roadtrip/actions/workflows/go.yml) [![Go Reference](https://pkg.go.dev/badge/github.com/nugget/roadtrip.svg)](https://pkg.go.dev/github.com/nugget/roadtrip) [![Go Report Card](https://goreportcard.com/badge/github.com/nugget/roadtrip)](https://goreportcard.com/report/github.com/nugget/roadtrip)

[Road Trip](https://darrensoft.ca/roadtrip/) is an iOS application written by
Darren Stone. This Go package provides methods and functions for reading and
parsing the backup files created by Road Trip so that you can work with this
data in your Go applications. Where possible it transforms the underlying 
Road Trip data into Go native data types and structures.

Road Trip itself supports native syncing of data between iOS devices via iCloud 
or Dropbox sync folders and the most convenient use of this package is to reference
a local, live updating copy of this sync directory on your device/host.

The roadtrip package is strictly read-only and does not allow for the creation of
new records to be pushed into the Road Trip app's data. It's safe to run against your
production/live sync files without harm.

This package was created by David "nugget" McNett and is not official or supported by
Darren Stone. Please don't bother the app developer with questions or feedback about this
package.


## Known Issues

### Only supports language "en"

The top section of each Road Trip CSV data file contains a block that
advertises the file's version, language, and CSV delimeter values. This package
has only been tested on english language data files. It relies on section
header names in the data file to distinguish data blocks and it's feasible that
these header strings change if you're running Road Trip in a different base
language. Additionally, no attempt is made to parse and utilize the delimeter
characters in the file info block.

<details>
<summary>Example data file header section</summary>
ROAD TRIP CSV ",."
Version,Language
1500,en
</details>


## Installation

`go get -u github.com/nugget/roadtrip`


## Links

- [Road Trip MPG iOS App](https://darrensoft.ca/roadtrip/)
- [Package Source](https://github.com/nugget/roadtrip)
