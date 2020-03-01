package output

import (
	"fmt"
	"regexp"
	"strings"
)

// TODO: Make ANSI color map configurable per player somehow
var ansiMap = map[string]string{
	// Black
	"<B>": "\u001B[30m",
	// Red
	"<R>": "\u001B[31m",
	// Green
	"<G>": "\u001B[32m",
	// Brown
	"<Br>": "\u001B[33m",
	// Blue
	"<Bl>": "\u001B[34m",
	// Purple
	"<P>": "\u001B[35m",
	// Cyan
	"<C>": "\u001B[36m",
	// White
	"<W>": "\u001B[37m",
	// Yellow
	"<Y>": "\u001B[33m",
	// White
	"<BW>": "\u001B[1;37m",
}

var ansiStart = regexp.MustCompile("<(B|R|G|Br|Bl|P|C|W|Y|BW)>")
var ansiEnd = regexp.MustCompile("</(B|R|G|Br|Bl|P|C|W|Y|BW)>")

var ansiReset = "\u001B[0m"
// Defaulting text color to white here so to save having to change everything to white later.
// This might not be default telnet behavior though
var defaultColor = ansiMap["<W>"]

func ANSIFormatter(output string) ([]byte, error) {
	startNum := len(ansiStart.FindAllString(output, -1))
	endNum := len(ansiEnd.FindAllString(output, -1))
	if startNum != endNum {
		return nil, fmt.Errorf("mismatch in number of opening and closing ANSI colors")
	}
	for elem, ansiCode := range ansiMap {
		output = strings.ReplaceAll(output, elem, ansiCode)
	}
	output = ansiEnd.ReplaceAllString(output, ansiReset + defaultColor)

	out := []byte(defaultColor + output + "\n")
	return out, nil
}
