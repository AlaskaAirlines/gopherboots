# gopherboots
[![Build Status](https://travis-ci.org/AlaskaAirlines/gopherboots.svg?branch=master)](https://travis-ci.org/AlaskaAirlines/gopherboots)
## Introduction 
This application was created to allow an administrator to bootstrap multiple hosts simultaneously. This functionality is not offered natively in Chef and places an unnecessary bottleneck on bootstrapping operations. Our ultimate aim is to remove as many barriers as possible from this process, enabling effective parallel bootstrapping at an enterprise level.

## Getting Started
In order to run this application, please ensure the following:
1. Go is installed and a Go environment has been created. Visit https://golang.org/doc/install for more information.
2. Install remote package for goqueue, using the following command:
```
go get github.com/Damnever/goqueue
```

## Usage
It's expected that you already have a working command line knife installation, and are able to run commands like `knife bootstrap` successfully.You'll need to set a few environment variables specific to your organization:
`SUPERUSER_NAME` should be set to a superuser account valid on the target host.
`SUPERUSER_PW` should be set to a superuser password valid on the target host.

A tsv file containing all the hosts you'd like to bootstrap should be formatted like this:
```
hostname	domain	target chef environment	runlist
```

For example:
```
test-host	example.org	linux	chef-client,base
```

To bootstrap all hosts simply run the following (where `./hosts.tsv` is the location of your tsv file):
```
./gopherboots -file=./hosts.tsv
```

If you encounter errors on any hosts during the bootstrapping process you can view them in the `./logs` directory.

## Build and Test
In order to build Gopherboots following code changes, run the following from your working directory:
```
go build
```
In order to test that your knife command generates correctly, edit the `Host` struct within `main_test.go` to reflect sample host data relevant to your environment. Subsequently, from your working directory, run the following:
```
go test
```

## Contribute
We encourage and appreciate contributions to this project form the Open Source community. In future updates we hope to include support for the following features:
- Bootstrapping Windows hosts
- Error reporting and follow-up queueing based on error group
- Re-creating `knife bootstrap` command functionality in Golang using Chef API calls, reducing Ruby resource bottleneck
- Additional unit and integration tests
