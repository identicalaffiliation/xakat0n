package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRun_InvalidCommand(t *testing.T) {
	t.Parallel()

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{
		"cmd",
		"-config", "invalidpath",
		"-command", "abracadabra",
	}

	err := run()
	require.Error(t, err)
}
