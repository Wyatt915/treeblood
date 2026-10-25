package treeblood

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

// Input that makes the parser panic is reported through the logger, not on
// standard output, so that a program using the library keeps its output.
func TestRecoveredPanicIsLogged(t *testing.T) {
	read, write, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	stdout := os.Stdout
	os.Stdout = write
	var logged bytes.Buffer
	SetLogOutput(&logged)
	defer SetLogOutput(os.Stderr)

	_, displayErr := DisplayStyle(`\hat{}`, nil)
	_, pitzErr := NewPitziil().DisplayStyle(`\hat{}`)
	NewPitziil().SemanticsOnly(`\hat{}`)

	os.Stdout = stdout
	if err := write.Close(); err != nil {
		t.Fatal(err)
	}
	printed, _ := io.ReadAll(read)
	if len(printed) > 0 {
		t.Errorf("printed to standard output:\n%s", printed)
	}
	if displayErr == nil || pitzErr == nil {
		t.Errorf("errors: %v, %v; want both set", displayErr, pitzErr)
	}
	if n := strings.Count(logged.String(), "recovered from"); n != 3 {
		t.Errorf("logged %d recovered panics, want 3:\n%s", n, logged.String())
	}
}
