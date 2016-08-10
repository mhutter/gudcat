package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

func init() {
	// check command line args and abort
	if len(os.Args) < 4 || len(os.Args) > 5 {
		fmt.Println("Usage: gudcat <address> <buffer size> <delay> [file]")
		fmt.Println("")
		fmt.Println("    buffer size in bytes")
		fmt.Println("    delay       in ms")
		fmt.Println("")
		fmt.Println("If no file is given, data is read from stdin.")
		os.Exit(1)
	}
}

func main() {
	// check target address
	addr, err := net.ResolveUDPAddr("udp", os.Args[1])
	check(err)

	// check buffer size & allocate buffer
	i, err := strconv.Atoi(os.Args[2])
	check(err)
	buf_size := int64(i)

	// check delay & convert to ms
	i, err = strconv.Atoi(os.Args[3])
	check(err)
	delay := time.Duration(i) * time.Millisecond

	// check & open file
	var source *os.File
	if len(os.Args) > 4 {
		source, err = os.Open(os.Args[4])
		defer source.Close()
		check(err)
	} else {
		source = os.Stdin
	}

	// dial target
	conn, err := net.DialUDP("udp", nil, addr)
	defer conn.Close()
	check(err)

	for {
		_, err := io.CopyN(conn, source, buf_size)
		if err == io.EOF {
			break
		}
		check(err)
		time.Sleep(delay)
	}

}

// check if error is nil, print message and exit program otherwise.
func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
