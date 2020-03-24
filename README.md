# Chatkit Retirement Announcement
We are sorry to say that as of April 23 2020, we will be fully retiring our
Chatkit product. We understand that this will be disappointing to customers who
have come to rely on the service, and are very sorry for the disruption that
this will cause for them. Our sales and customer support teams are available at
this time to handle enquiries and will support existing Chatkit customers as
far as they can with transition. All Chatkit billing has now ceased , and
customers will pay no more up to or beyond their usage for the remainder of the
service. You can read more about our decision to retire Chatkit here:
[https://blog.pusher.com/narrowing-our-product-focus](https://blog.pusher.com/narrowing-our-product-focus).
If you are interested in learning about how you can build chat with Pusher
Channels, check out our tutorials.

# chatkit-server-go [![godoc-badge][]][GoDoc] [![Build Status](https://travis-ci.org/pusher/chatkit-server-go.svg?branch=master)](https://travis-ci.org/pusher/chatkit-server-go)

Golang server SDK for [Pusher Chatkit][].

This package provides an interface to interact with the Chatkit service. It allows
interacting with the core Chatkit service and other subservices.

Please report any bugs or feature requests via a GitHub issue on this repo.

## Installation

    $ go get github.com/pusher/chatkit-server-go

## Go versions

This library requires Go versions >=1.9.

## Getting started

Before processing, ensure that you have created an account on the [Chatkit Dashboard](https://dash.pusher.com).
You'll then have to create a Chatkit instance and acquire credentials for it (all of which is available in the dashboard.)

A new client may be instantiated as follows

```go
import "github.com/pusher/chatkit-server-go"

client, err := chatkit.NewClient("<YOUR_INSTANCE_LOCATOR>", "<YOUR_KEY>")
if err != nil {
	return err
}

// use client to make calls to the service
```

## Deprecated versions

Versions of the library below [1.0.0](https://github.com/pusher/chatkit-server-go/releases/tag/1.0.0) are deprecated and support for them will soon be dropped.

It is highly recommended that you upgrade to the latest version if you're on an older version. To view a list of changes,
please refer to the [CHANGELOG][].


## Tests

To run the tests, a Chatkit instance is required along with its credentials. `CHATKIT_INSTANCE_LOCATOR` and `CHATKIT_INSTANCE_KEY` are required
to be set as environment variables. Note that the tests run against an actual cluster and are to be treated as integration tests.

    $ go test -v

## Documentation

Refer to [GoDoc][].

## License

This code is free to use under the terms of the MIT license. Please refer to
LICENSE.md for more information.

[GoDoc]: http://godoc.org/github.com/pusher/chatkit-server-go
[Pusher Chatkit]: https://pusher.com/chatkit
[godoc-badge]: https://godoc.org/github.com/pusher/chatkit-server-go?status.svg
[CHANGELOG]: CHANGELOG.md
