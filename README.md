# ipod
[![Join the chat at https://gitter.im/ipod-gadget/Lobby](https://badges.gitter.im/ipod-gadget/Lobby.svg)](https://gitter.im/ipod-gadget/Lobby?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![GoDoc](https://godoc.org/github.com/oandrew/ipod?status.svg)](https://godoc.org/github.com/oandrew/ipod)
[![Build Status](https://travis-ci.org/oandrew/ipod.svg?branch=master)](https://travis-ci.org/oandrew/ipod)

ipod is a golang userspace library implementation of the ipod accessory protocol.
It includes an example client for use with https://github.com/oandrew/ipod-gadget

This is a total rewrite of what was included with the  ipod-gadget project. 
It should work as a drop-in replacement for the old app.

New features:
- Storing and replaying traces
- Detailed verbose logging for debug
- Better codebase with message type definitions
- Tests



# build and run
```
go build github.com/oandrew/ipod/cmd/ipod
# or cross compiling
GOOS=linux GOARCH=arm GOARM=6 go build github.com/oandrew/ipod/cmd/ipod

# with debug logging
./ipod -d serve /dev/iap0

# save a trace file
./ipod -d serve -w ipod.trace /dev/iap0

# simulate incoming requests from a trace file
./ipod -d replay ./ipod.trace

# view a trace file
./ipod -d view ./ipod.trace

```

Client app godoc https://godoc.org/github.com/oandrew/ipod/cmd/ipod

Refer to https://github.com/oandrew/ipod-gadget for more info on how to get the kernel part working.






