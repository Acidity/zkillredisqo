zkillredisqo
=========

[![Build Status](https://travis-ci.org/MorpheusXAUT/zkillredisqo.svg?branch=master)](https://travis-ci.org/MorpheusXAUT/zkillredisqo) [![GoDoc](https://godoc.org/github.com/MorpheusXAUT/zkillredisqo?status.svg)](https://godoc.org/github.com/MorpheusXAUT/zkillredisqo)

A small library for subscribing to zKillboard's RedisQ kill feed, written in Golang. The name originates from me being *incredibly* uncreative and combining "zKillboard", "RedisQ" and "Go".

Package zkillredisqo provides a lightweight library for interfacing with zKillboard's RedisQ service for receiving killmails in realtime.

The library allows for applications to wait for new kills using a channel. Errors are handed off to the calling app via a separate error channel.

Installation
------

```bash
go get -u github.com/morpheusxaut/zkillredisqo
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

Attribution
------

### zKillboard
zKillboard's RedisQ is courtesy of [zKillboard](https://zkillboard.com), legal information can be found [here](https://zkillboard.com/information/legal/). Great thanks to [Squizz Caphinator](https://zkillboard.com/character/1633218082/) for preventing this service.

### CCP
EVE Online and the EVE logo are the registered trademarks of CCP hf. All rights are reserved worldwide. All other trademarks are the property of their respective owners. EVE Online, the EVE logo, EVE and all associated logos and designs are the intellectual property of CCP hf. All artwork, screenshots, characters, vehicles, storylines, world facts or other recognizable features of the intellectual property relating to these trademarks are likewise the intellectual property of CCP hf. CCP hf. has granted permission to MorpheusXAUT to use EVE Online and all associated logos and designs for promotional and information purposes on its website but does not endorse, and is not in any way affiliated with, MorpheusXAUT. CCP is in no way responsible for the content on or functioning of this website, nor can it be liable for any damage arising from the use of this website.

License
------

[MIT License](https://opensource.org/licenses/mit-license.php)