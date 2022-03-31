# go-mixcloud-search

A simple cobra-based cli for searching mixcloud.

<p align="left">
<img src="https://github.com/darren-reddick/go-mixcloud-search/actions/workflows/cicd.yml/badge.svg?branch=main">
<img src="https://img.shields.io/badge/License-Apache%202.0-blue.svg">
</p>

## :musical_keyboard: Overview

[Mixcloud](https://www.mixcloud.com/) is an amazing service for discovering and listening to radio shows, DJ mixes and podcasts.

The mixcloud web interface can be cumbersome when trying to search for mixes with specific terms etc. This project was created as a way to search mixcloud for mix data with a quick and simple search.

The cli supports searching for a term on mixcloud and then filtering the results (client-side) using arrays of exclude and include terms.

Currently the results are output as json to file **test.json** (early days!).

## :factory: Installing

### :page_facing_up: Prerequisites

- Go (>=1.17)

### :wrench: Building the cli

Select the appropriate build command from below. This will output the cli binary into the current directory.

```
# Mac M1
GOOS=darwin GOARCH=arm64 go build -o gmc
# Mac Intel
GOOS=darwin GOARCH=amd64 go build -o gmc
# Linux
GOOS=linux go build -o gmc
```


## :notebook_with_decorative_cover: Usage

The cli supports the following subcommands

### :eyeglasses: search

```
Search for mixes on mixcloud by term.

Usage:
  gmc search [flags]

Flags:
  -d, --debug             Enable debug
  -e, --exclude strings   Filter to exclude entry
  -h, --help              help for search
  -i, --include strings   Filter to include entry
  -l, --limit int         Limit number of results
  -t, --term string       Search term
```

#### :stars: Example usage

```
gmc search --term truncate --include berghain --exclude 2016
```



### :scroll: history

```
Search listen history for a user.

Usage:
  gmc history [flags]

Flags:
  -d, --debug             Enable debug
  -e, --exclude strings   Filter to exclude entry
  -h, --help              help for history
  -i, --include strings   Filter to include entry
  -l, --limit int         Limit number of results
  -u, --user string       User name to search
```

#### :stars: Example usage

```
gmc history --user akumad --include party
```

