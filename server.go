package gudcat

import (
	"net"
	"os"
	"time"
)

// Server contains options to listen UDP data and write it to stdout.
type Server struct {
	Addr    *net.UDPAddr
	Timeout time.Duration
	started bool
	err     error
}

// NewServer returns a Server that listens on addr
func NewServer(addr *net.UDPAddr, timeout time.Duration) *Server {
	return &Server{
		Addr:    addr,
		Timeout: timeout,
		started: false,
	}
}

// Run starts the server
func (server *Server) Run() error {
	// open listener
	conn, err := net.ListenUDP("udp", server.Addr)
	// remember to close the connection on exit
	defer conn.Close()
	if err != nil {
		return err
	}

	buffer := make([]byte, 65535)
	hasTimeout := server.Timeout > time.Duration(0)

	for {
		// read data from conn
		n, _, err := conn.ReadFromUDP(buffer)

		// check err
		if server.isErr(err) {
			break
		}

		// write data to Stdout
		os.Stdout.Write(buffer[0:n])

		if hasTimeout {
			// set the deadline if neccessary
			conn.SetReadDeadline(time.Now().Add(server.Timeout))
		}
	}

	return server.err
}

func (server *Server) isErr(err error) bool {
	if err != nil {
		// we HAVE an error...
		if ne, ok := err.(*net.OpError); ok {
			if !ne.Timeout() {
				// ... and it's not a timeout!
				server.err = err
			}
		} else {
			// its something else
			server.err = err
		}

		return true
	}

	return false
}
