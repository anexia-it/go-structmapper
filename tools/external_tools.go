//go:build tools
// +build tools

// This files purpose is to prevent `go mod tidy` from removing the tools as "unused" from go.mod.

package main

import (
	_ "github.com/client9/misspell/cmd/misspell"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
)
