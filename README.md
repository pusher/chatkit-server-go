# chatkit-server-go [![godoc-badge][]][GoDoc]

package chatkit is the Golang server SDK for [Pusher Chatkit][].

This package provides the Client type for managing Chatkit users and
interacting with roles and permissions of those users. It also contains some
helper functions for creating your own JWT tokens for authentication with the
Chatkit service.

Please report any bugs or feature requests via a GitHub issue on this repo.

## Installation

    $ go get github.com/pusher/chatkit-server-go

## Getting Started

Please refer to the [`/examples`][] directory or the [GoDoc][].

## Tests

    $ go test -v -cover

## Documentation

Refer to [GoDoc][].

## License

This code is free to use under the terms of the MIT license. Please refer to
LICENSE.md for more information.

[GoDoc]: http://godoc.org/github.com/pusher/chatkit-server-go
[Pusher Chatkit]: https://pusher.com/chatkit
[`/auth_example`]: auth_example/
[`/examples`]: examples/
[godoc-badge]: https://godoc.org/github.com/pusher/chatkit-server-go?status.svg
