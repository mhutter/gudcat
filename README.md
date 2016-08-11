# gudcat

[![GoDoc](https://godoc.org/github.com/mhutter/gudcat?status.svg)](https://godoc.org/github.com/mhutter/gudcat)

gudcat lets you send data via UDP in a more or less controlled manner.

## Synopsis

```bash
# on the server
gudcat server :3000 > target.file

# on the client
gudcat client -size 512 -delay 1ms < source.file
```

## Installation

    go get -u github.com/mhutter/gudcat/gudcat

(yep, thats `gudcat` twice)

## Usage

    Usage:
        gudcat client [options] address
        gudcat server address

### Client

    Usage:
        gudcat client [options] address

    Read data from stdin and send it to <address> via UDP.

    Usage of flags:
      -delay duration
        	delay between packages
      -size int
        	package size in bytes (default 64000)

    Examples:
        gudcat client -delay 10ms -size 510 localhost:3388 < input.file
        gudcat client '[fe80::1]:3388'


### Server

    Usage:
        gudcat server address

    Listen on <address> for data and print it to stdout

    Examples:
        gudcat server :3388
        gudcat server [::]:3388
        gudcat server 127.0.0.1:3388
