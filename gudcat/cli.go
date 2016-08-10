package main

import (
	"flag"
	"fmt"
	"log"
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
	fs := flag.NewFlagSet("flags", flag.ContinueOnError)

	// Flag definitions
	delay := fs.Duration("delay", time.Duration(0), "delay between packages")
	size := fs.Int64("size", 64000, "package size in bytes")

	// parse flags, print usage info on failure
	err := fs.Parse(args)
	if err != nil || fs.NArg() < 2 {
		clientUsage()
		os.Exit(1)
	}

	// resolve address, print usage info and quit on failure
	addr, err := net.ResolveUDPAddr("udp", fs.Arg(1))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving listen address: %s\n", err)
		fmt.Println()
		clientUsage()
		os.Exit(3)
	}

	client := gudcat.NewClient(addr, *size, *delay)
	client.Run()
}

// start the server
func server(args []string) {
	if len(args) != 1 || isHelp(args[0]) {
		serverUsage()
		os.Exit(1)
	}

	addr, err := net.ResolveUDPAddr("udp", args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving listen address: %s\n", err)
		fmt.Println()
		serverUsage()
		os.Exit(3)
	}
	server := gudcat.NewServer(addr)
	err = server.Run()
	if err != nil {
		log.Fatalln(err)
	}
}

func isHelp(str string) bool {
	return (str == "-h" || str == "-help" || str == "--help")
}

// prints usage information
func usage() {
	fmt.Println("Usage:")
	fmt.Println("    \u001b[1mgudcat client\u001b[22m [\u001b[4moptions\u001b[24m] \u001b[4maddress\u001b[24m")
	fmt.Println("    \u001b[1mgudcat server\u001b[22m \u001b[4maddress\u001b[24m")
	fmt.Println("")
	fmt.Println("See `gudcat <command> -h` for more info")
}

// prints server usage information
func serverUsage() {
	fmt.Println("Usage:")
	fmt.Println("    \u001b[1mgudcat server\u001b[22m \u001b[4maddress\u001b[24m")
	fmt.Println("")
	fmt.Println("Listen on <address> for data and print it to stdout")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("    gudcat server :3388")
	fmt.Println("    gudcat server [::]:3388")
	fmt.Println("    gudcat server 127.0.0.1:3388")
}

// prints client usage information
func clientUsage() {
	fmt.Println("Usage:")
	fmt.Println("    \u001b[1mgudcat client\u001b[22m [\u001b[4moptions\u001b[24m] \u001b[4maddress\u001b[24m")
	fmt.Println("")
	fmt.Println("Read data from stdin and send it to <address>.")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("    gudcat client -delay 10ms -size 510 localhost:3388 < input.file")
	fmt.Println("    gudcat client '[fe80::1]:3388'")
}
