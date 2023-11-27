[![build](https://github.com/tkitsunai/messenger/actions/workflows/go-build.yml/badge.svg?branch=master)](https://github.com/tkitsunai/messenger/actions/workflows/go-build.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/tkitsunai/messenger)](https://goreportcard.com/report/github.com/tkitsunai/messenger)

# Messenger

Messenger is a simple pub/sub based on the idea that messages are tagged.

The server and client exchange messages using a TCP connection and a JSON encoder/decoder.

## Install

```bash
go install github.com/tkitsunai/messenger
```

## Short-Term Roadmap

- Support Subject Matcher.
  - Provides Publishing/Subscription for subject lines using pattern matching, not just simple strings.
- Support for Administrators.
  - Provide performance and server metrics.
  - Add authorization role and features.

## Medium-Term Roadmap

- Support server clustering.
- Support for message persistent.

## Long-Term Roadmap

- ideal
- ideal
- ideal
