package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAddr(t *testing.T) {
	type test struct {
		in    string
		proto string
		host  string
		port  string
	}

	cases := []test{
		{"udp://192.168.0.1:8080", "udp", "192.168.0.1", "8080"},
		{"tcp://rwlist.io:443", "tcp", "rwlist.io", "443"},
	}

	for _, c := range cases {
		proto, host, port, ok := parseAddr(c.in)
		assert.Equal(t, c.proto, proto)
		assert.Equal(t, c.host, host)
		assert.Equal(t, c.port, port)
		assert.Equal(t, true, ok)
	}
}
