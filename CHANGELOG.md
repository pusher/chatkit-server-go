# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased](https://github.com/pusher/chatkit-server-go/compare/0.2.0...HEAD)

## [1.0.0](https://github.com/pusher/chatkit-server-go/compare/0.2.0...1.0.0)

The SDK has been rewritten from the ground up, so best to assume that
everything has changed and refer to the [GoDoc][].

## [0.2.0](https://github.com/pusher/chatkit-server-go/compare/0.1.0...0.2.0) - 2018-04-24

### Changes

- `TokenManager` renamed to `Authenticator`

### Removals

- `NewChatkitUserToken` has been removed
- `NewChatkitSUToken` has been removed

### Additions

- `NewChatkitToken` has been added and essentially replaces `NewChatkitSUToken` and `NewChatkitUserToken`
- `Authenticate` added to `Client`

[GoDoc]: http://godoc.org/github.com/pusher/chatkit-server-go
