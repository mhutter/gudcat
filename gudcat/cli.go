/*
Command gudcat either starts a server to listen for data (gudcat server), or
starts sending data itself (gudcat client).

	Usage
	    gudcat client [options] address
	    gudcat server address

	See `gudcat <command> -h` for more info


Server Usage:
    gudcat server [options] address

Listen on <address> for data and print it to stdout

Options:
    -timeout duration
        Valid time units are 'ns', 'us' (or 'µs'), 'ms', 's', 'm', 'h'.

Timeout: The timeout will only start once the server has
received some data. New data will reset the timeout.

Examples:
    gudcat server -timeout 1s :3388
    gudcat server [::]:3388
    gudcat server 127.0.0.1:3388



Client Usage:
    gudcat client [options] address

Read data from stdin and send it to <address>.

Options:
    -delay duration
        delay between packages. Valid time units are 'ns', 'us' (or 'µs'), 'ms', 's', 'm', 'h'.
    -size int
        package size in bytes (default 64000)

Examples:
    gudcat client -delay 10ms -size 510 localhost:3388 < input.file
    gudcat client '[fe80::1]:3388'
*/
package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/mhutter/gudcat"
)

type command func([]string)

var commands = map[string]command{
	"client": client,
	"server": server,
}

func main() {
	if len(os.Args) < 2 || commands[os.Args[1]] == nil {
		usage()
		os.Exit(1)
	}

	commands[os.Args[1]](os.Args[2:])
}

// start the client
func client(args []string) {
	if len(args) < 1 {
		args = []string{"-h"}
	}
	fs := flag.NewFlagSet("flags", flag.ExitOnError)
	fs.Usage = clientUsage(fs)

	// Flag definitions
	delay := fs.Duration("delay", time.Duration(0), "Delay between packages. Valid time units are 'ns', 'us' (or 'µs'), 'ms', 's', 'm', 'h'.")
	size := fs.Int64("size", 64000, "package size in bytes")

	// parse flags, print usage info and exit on failure
	err := fs.Parse(args)

	// resolve address, print usage info and quit on failure
	addr, err := net.ResolveUDPAddr("udp", fs.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving listen address: %s\n", err)
		fmt.Println()
		fs.Usage()
		os.Exit(1)
	}

	client := gudcat.NewClient(addr, *size, *delay)
	if err := client.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

// start the server
func server(args []string) {
	if len(args) < 1 {
		args = []string{"-h"}
	}

	fs := flag.NewFlagSet("flags", flag.ExitOnError)
	fs.Usage = serverUsage(fs)

	// Flag definitions
	timeout := fs.Duration("timeout", time.Duration(0), "Valid time units are 'ns', 'us' (or 'µs'), 'ms', 's', 'm', 'h'.")

	// parse flags, print usage info and exit on failure.
	err := fs.Parse(args)

	addr, err := net.ResolveUDPAddr("udp", fs.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving listen address: %s\n", err)
		fmt.Println()
		fs.Usage()
		os.Exit(1)
	}

	server := gudcat.NewServer(addr, *timeout)
	if err := server.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

// prints usage information
func usage() {
	fmt.Println("Usage:")
	fmt.Println("    \u001b[1mgudcat client\u001b[22m [\u001b[4moptions\u001b[24m] \u001b[4maddress\u001b[24m")
	fmt.Println("    \u001b[1mgudcat server\u001b[22m [\u001b[4moptions\u001b[24m] \u001b[4maddress\u001b[24m")
	fmt.Println("")
	fmt.Println("See `gudcat <command> -h` for more info")
}

// prints server usage information
func serverUsage(fs *flag.FlagSet) func() {
	return func() {
		fmt.Println("Usage:")
		fmt.Println("    \u001b[1mgudcat server\u001b[22m [\u001b[4moptions\u001b[24m] \u001b[4maddress\u001b[24m")
		fmt.Println("")
		fmt.Println("Listen on <address> for data and print it to stdout")
		fmt.Println("")
		fmt.Println("Options:")
		fs.PrintDefaults()
		fmt.Println("")
		fmt.Println("Timeout: The timeout will only start once the server has")
		fmt.Println("received some data. New data will reset the timeout.")
		fmt.Println("")
		fmt.Println("Examples:")
		fmt.Println("    gudcat server -timeout 1s :3388")
		fmt.Println("    gudcat server [::]:3388")
		fmt.Println("    gudcat server 127.0.0.1:3388")
	}
}

// prints client usage information
func clientUsage(fs *flag.FlagSet) func() {
	return func() {
		fmt.Println("Usage:")
		fmt.Println("    \u001b[1mgudcat client\u001b[22m [\u001b[4moptions\u001b[24m] \u001b[4maddress\u001b[24m")
		fmt.Println("")
		fmt.Println("Read data from stdin and send it to <address>.")
		fmt.Println("")
		fmt.Println("Options:")
		fs.PrintDefaults()
		fmt.Println("")
		fmt.Println("Examples:")
		fmt.Println("    gudcat client -delay 10ms -size 510 localhost:3388 < input.file")
		fmt.Println("    gudcat client '[fe80::1]:3388'")
	}
}
