//go:build tools

package internal

// this file exists to forcefully vendor dependencies used by go:generate

import (
	_ "github.com/dmarkham/enumer"
	_ "go.uber.org/mock/mockgen"
)
