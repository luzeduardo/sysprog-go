package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestMainProgram(t *testing.T) {
	var stdoutBuf, stderrBuf bytes.Buffer
	config, err := NewCliConfig(WithOutStream(&stdoutBuf), WithErrStream(&stderrBuf))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	app([]string{"main", "game", "golang", "error"}, config)
	output := stdoutBuf.String()
	if len(output) == 0 {
		t.Fatal("Expected output got nothing")
	}
	if !strings.Contains(output, "word game is even") {
		t.Fatalf("Expected output does not contain 'word game is even'")
	}

	if !strings.Contains(output, "word golang is even") {
		t.Fatal("Expected output does not contain 'word golang is even")
	}

	errors := stderrBuf.String()
	if len(errors) == 0 {
		t.Fatal("Expected errors got nothing")
	}

	if !strings.Contains(errors, "word error is odd") {
		t.Fatal("Expected errors does not contain 'word error is odd'")
	}

}
