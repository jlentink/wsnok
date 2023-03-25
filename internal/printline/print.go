package printline

import (
	"fmt"
	"os"
)

var Debug = false
var Out = os.Stdout

// Print prints a line to stdout
func Print(isDebug bool, s string) {
	if isDebug && Debug || !isDebug {
		fmt.Fprint(Out, s)
	}
}

// Printf prints a line to stdout
func Printf(isDebug bool, s string, args ...interface{}) {
	if isDebug && Debug || !isDebug {
		fmt.Fprintf(Out, s, args...)
	}
}
