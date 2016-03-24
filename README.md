# CatchpointSDK

An unofficial SDK and CLI for Catchpoint's Pull API

## Overview

_This package is incomplete and still actively being worked on; consider it unstable for now._

This package will allow for some basic interactions with Catchpoint's pull API; currently supported functionality includes:

*  Generating an auth token
*  Displaying products and folders (divisions are supported)
*  Displaying and adding tests (all tests options are not yet supported)

## Installation

Make sure you have a working Go environment. [See the install instructions](http://golang.org/doc/install.html).

To install, run:
```
$ go get github.com/ves/catchpointsdk
```

Make sure your `PATH` includes to the `$GOPATH/bin` directory so your commands can be easily used:
```
export PATH=$PATH:$GOPATH/bin
```

You will then need to build the binary (if you wish to use the command line utility):
```
$ go install github.com/ves/catchpointsdk/catchpoint
```

## Running the command line binary

You must enable the Pull API first (see Settings > API in Catchpoint).

You will then need to set some environment variables:

```
$ export CATCHPOINTSDK_CLIENTID="<key>"
$ export CATCHPOINTSDK_CLIENTSECRET="<secret>"
$ export CATCHPOINTSDK_ENDPOINT="https://io.catchpoint.com"
```

You can also set an environment variable for switching divisions; note that you will need to know the division_id to do this.

```
$ export CATCHPOINTSDK_DIVISION_ID="<division id>"
```
