package util

import (
	"net"
	"os"
	"strings"
)

// NewListener checks the given address and returns the corresponding
// net.Listener for either the "tcp" or "unix" network domains. addr
// may either be a filesystem path, or a "address:port" string.
func NewListener(addr string) (net.Listener, error) {
	if strings.Contains(addr, ":") {
		return net.Listen("tcp", addr)
	}

	if IsUnixSocket(addr) {
		os.Remove(addr) //nolint(gosec)
	}

	return net.Listen("unix", addr)
}

// IsUnixSocket returns true if the path is a unix domain socket.
func IsUnixSocket(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}

	return (s.Mode() & os.ModeSocket) == os.ModeSocket
}
