zkillredisqo
=========

[![Build Status](https://travis-ci.org/MorpheusXAUT/zkillredisqo.svg?branch=master)](https://travis-ci.org/MorpheusXAUT/zkillredisqo) [![GoDoc](https://godoc.org/github.com/MorpheusXAUT/zkillredisqo?status.svg)](https://godoc.org/github.com/MorpheusXAUT/zkillredisqo)

A small library for subscribing to zKillboard's RedisQ kill feed, written in Golang. The name originates from me being *incredibly* uncreative and combining "zKillboard", "RedisQ" and "Go".

Package zkillredisqo provides a lightweight library for interfacing with zKillboard's RedisQ service for receiving killmails in realtime.

The library allows for applications to wait for new kills using a channel. Errors are handed off to the calling app via a separate error channel.

Installation
------

```bash
go get -u github.com/morpheusxaut/zkillboardqo
```

Usage
------

A tiny sample application demonstrating the usage of zkillredisqo is included in the `cmd/simple.go`.
In general, applications should create a `Poller` and receive messages from its `Kills` and `Errors` channels:

```go
package main

import (
	"github.com/morpheusxaut/zkillredisqo"
	"log"
)

func main() {
	poller := zkillredisqo.NewPoller(nil)

	for {
		select {
		case kill := <-poller.Kills:
			log.Printf("%+v\n", kill)
			break
		case err := <-poller.Errors:
			log.Printf("*** ERROR: %v\n", err)
			break
		}
	}
}
```

Documentation
------
see https://godoc.org/github.com/MorpheusXAUT/zkillredisqo


License
------

[MIT License](https://opensource.org/licenses/mit-license.php)