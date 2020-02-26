package telnet

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGmcpNegotiation(t *testing.T) {
	telnet := Telnet{
		Connection:         nil,
		InputSteam:         nil,
		Data:           make(chan []byte, 1),
		negState:       0,
		subnegOffset:   0,
		subnegotiation: make([]byte, 2048),
	}

	testGmcp := []byte{iac, sb, optGmcp, 116, 101, 115, 116, 103, 109, 99, 112, 32, 123, 34, 100, 97, 116, 97, 34, 58, 32, 34, 102, 111, 111, 34, 125, iac, se}
	outBuffer := make([]byte, 2048)
	size, err := telnet.negotiate(outBuffer, testGmcp)
	if err != nil {
		t.Errorf("failed to subnegotiate: %v", err)
	}
	gmcpData := <-telnet.Data

	assert.Equal(t, -1, size)
	// Indicates an "empty" buffer
	assert.Equal(t, byte(0), outBuffer[0])
	assert.Equal(t, "testgmcp {\"data\": \"foo\"}", string(gmcpData))
}

func TestNormalNegotiation(t *testing.T) {
	telnet := Telnet{
		Connection:         nil,
		InputSteam:         nil,
		Data:           make(chan []byte, 1),
		negState:       0,
		subnegOffset:   0,
		subnegotiation: make([]byte, 2048),
	}

	testData := []byte{102, 111, 111, 32, 98, 97, 114, 32, 98, 97, 122}
	outBuffer := make([]byte, 2048)
	size, err := telnet.negotiate(outBuffer, testData)
	if err != nil {
		t.Errorf("failed to negotiate: %v", err)
	}

	// Trim NUL bytes for the sake of the test, this is done in the main function
	outBuffer = bytes.Trim(outBuffer, "\x00")
	assert.Equal(t, 11, size)
	assert.Equal(t, "foo bar baz", string(outBuffer))
}
