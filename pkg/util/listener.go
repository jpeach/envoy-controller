package util

import (
	"net"
	"strings"
)

// NewListener checks the given address and returns the corresponding
// net.Listener for either the "tcp" or "unix" network domains. addr
// may either be a filesystem path, or a "address:port" string.
func NewListener(addr string) (net.Listener, error) {
	if strings.Contains(addr, ":") {
		return net.Listen("tcp", addr)
	}

	return net.Listen("unix", addr)
}
