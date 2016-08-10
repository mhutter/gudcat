// Package gudcat lets you send data via UDP in a more or less controlled manner
package gudcat

import (
	"io"
	"net"
	"os"
	"time"
)

// Client contains options for connecting to an address and sending data.
type Client struct {
	Addr  *net.UDPAddr
	Size  int64
	Delay time.Duration
}

// NewClient returns a Client which sends data to addr. Data is sent in Packages
// of `size` bytes, with a delay in between packages.
func NewClient(addr *net.UDPAddr, size int64, delay time.Duration) *Client {
	return &Client{Addr: addr, Size: size, Delay: delay}
}

// Run starts sending data from stdin to Addr.
func (client *Client) Run() error {
	conn, err := net.DialUDP("udp", nil, client.Addr)
	defer conn.Close()
	if err != nil {
		return err
	}

	// buf := make([]byte, client.Size)

	for {
		// send Size bytes of data
		_, err := io.CopyN(conn, os.Stdin, client.Size)
		if err == io.EOF {
			// all data is sent
			return nil
		} else if err != nil {
			return err
		}

		time.Sleep(client.Delay)
	}
}
