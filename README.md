# chatkit-server-go [![godoc-badge][]][GoDoc]

Golang server SDK for [Pusher Chatkit][].

This package provides an interface to interact with the Chatkit service. It allows
interacting with the core Chatkit service and other subservices.

Please report any bugs or feature requests via a GitHub issue on this repo.

## Installation

    $ go get github.com/pusher/chatkit-server-go

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

To run the tests, a Chatkit instance is required along with its credentials. `TEST_INSTANCE_LOCATOR` and `TEST_KEY` are required
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
