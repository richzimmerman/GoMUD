package output

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestANSIFormatter(t *testing.T) {
	in := "<B>Some colored</B> Text"
	expectedOut := []byte("\n\u001B[37m\u001B[30mSome colored\u001B[0m\u001B[37m Text\n\n")

	out, err := ANSIFormatter(in)
	assert.Equal(t, expectedOut, out)
	assert.Nil(t, err)
}

func TestANSIFormatterFail(t *testing.T) {
	in := "<W>Mismatched ANSI coloring"

	out, err := ANSIFormatter(in)
	assert.Nil(t, out)
	assert.NotNil(t, err)
}
