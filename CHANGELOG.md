# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased](https://github.com/pusher/chatkit-server-go/compare/2.0.0...HEAD)

## [2.1.0](https://github.com/pusher/chatkit-server-go/compare/2.1.0...2.0.0)

### Additions

- Support for `PushNotificationTitleOverride` attribute in the Room model and
corresponding Update and Create structs.

## [2.0.0](https://github.com/pusher/chatkit-server-go/compare/1.2.0...2.0.0)

### Additions

- Support for user specified room IDs. Provide the `ID` parameter to the
  `CreateRoom` method.

### Changes

- The `DeleteMessage` method now *requires* a room ID parameter, `RoomID`, and
  the `ID` parameter has been renamed to `MessageId` to avoid ambiguity.

## [1.2.0](https://github.com/pusher/chatkit-server-go/compare/1.1.0...1.2.0)

- Multipart message support: `SendSimpleMessage`, `SendMultipartMessage`,
  `FetchMessagesMultipart`, and `SubscribeToRoomMultipart` deal in the
  multipart message format.

## [1.1.0](https://github.com/pusher/chatkit-server-go/compare/1.0.0...1.1.0) - 2018-11-07

### Additions

- A `CustomData` attribute for the `UpdateRoomOptions` and `CreateRoomOptions` structs.
- A `CustomData` attribute for the `Room` model.

## [1.0.0](https://github.com/pusher/chatkit-server-go/compare/0.2.0...1.0.0) - 2018-10-30

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
