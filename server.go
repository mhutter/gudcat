package gudcat

import (
	"io"
	"net"
	"os"
)

// Server contains options to listen UDP data and write it to stdout.
type Server struct {
	Addr *net.UDPAddr
}

// NewServer returns a Server that listens on addr
func NewServer(addr *net.UDPAddr) *Server {
	return &Server{Addr: addr}
}

// Run starts the server
func (server *Server) Run() error {
	conn, err := net.ListenUDP("udp", server.Addr)
	defer conn.Close()
	if err != nil {
		return err
	}

	for {
		// copy data directly to stdout

		if _, err := io.Copy(os.Stdout, conn); err != nil {
			return err
		}
	}
}
