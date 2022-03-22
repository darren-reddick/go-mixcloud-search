# go-mixcloud-search

A simple cobra-based cli for searching mixcloud.

<p align="left">
<img src="https://github.com/darren-reddick/go-mixcloud-search/actions/workflows/cicd.yml/badge.svg?branch=main">
<img src="https://img.shields.io/badge/License-Apache%202.0-blue.svg">
</p>

## :musical_keyboard: Overview

The mixcloud web interface can be cumbersome when trying to search for mixes with specific terms etc. This project was created as a way to search mixcloud for mix data with a quick and simple search.

## :notebook_with_decorative_cover: Usage

The cli supports the following subcommands

### :eyeglasses: search

```
Search for mixes on mixcloud by term.

Usage:
  gmc search [flags]

Flags:
  -e, --exclude strings   Filter to exclude entry
  -h, --help              help for search
  -i, --include strings   Filter to include entry
  -l, --limit int         Limit number of results
  -t, --term string       Search term
```

### :scroll: history

```
Search listen history for a user.

Usage:
  gmc history [flags]

Flags:
  -e, --exclude strings   Filter to exclude entry
  -h, --help              help for history
  -i, --include strings   Filter to include entry
  -l, --limit int         Limit number of results
  -u, --user string       User name to search
```

