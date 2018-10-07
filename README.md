# wraperr
[![Build Status](https://travis-ci.com/srvc/wraperr.svg?branch=master)](https://travis-ci.com/srvc/wraperr)
[![codecov](https://codecov.io/gh/srvc/wraperr/branch/master/graph/badge.svg)](https://codecov.io/gh/srvc/wraperr)
[![Go project version](https://badge.fury.io/go/github.com%2Fsrvc%2Fwraperr.svg)](https://badge.fury.io/go/github.com%2Fsrvc%2Fwraperr)
[![Go Report Card](https://goreportcard.com/badge/github.com/srvc/wraperr)](https://goreportcard.com/report/github.com/srvc/wraperr)
[![license](https://img.shields.io/github/license/srvc/wraperr.svg)](./LICENSE)

Check that error return value are wrapped

## Install

```
go get -u github.com/srvc/wraperr/cmd/wraperr
```


## Usage

To check all packages beneath the current directory:

```
wraperr ./...
```


## Inspired

- [errcheck](https://github.com/kisielk/errcheck)
