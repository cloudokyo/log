//go:build !binary_log
// +build !binary_log

package log

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFuncName(t *testing.T) {
	testcases := []struct {
		input  string
		output string
	}{
		{
			input:  "Log",
			output: "Log",
		},
		{
			input:  "Debug",
			output: "Debug",
		},
		{
			input:  "(*service).Add",
			output: "Add",
		},
		{
			input:  "(*service).Add.1",
			output: "Add.1",
		},
		{
			input:  "(*service).Process",
			output: "Process",
		},
		{
			input:  "(*service).Process.2",
			output: "Process.2",
		},
		{
			input:  "",
			output: "",
		},
	}

	for _, test := range testcases {
		fmt.Println(funcname(test.input), "-->", test.output)
		assert.Equal(t, funcname(test.input), test.output)
	}
}

func TestPackageName(t *testing.T) {
	testcases := []struct {
		input  string
		output string
	}{
		{
			input:  "/go/pkg/mod/github.com/labstack/echo/v4@v4.12.0/echo.go",
			output: "echo",
		},
		{
			input:  "/go/pkg/mod/github.com/labstack/echo/v4@v4.12.0/middleware/recover.go",
			output: "middleware",
		},
		{
			input:  "/go/1.23.0/libexec/src/testing/testing.go",
			output: "testing",
		},
		{
			input:  "/go/testing.go",
			output: "go",
		},
		{
			input:  "/testing.go",
			output: "main",
		},
		{
			input:  "testing.go",
			output: "testing",
		},
		{
			input:  "",
			output: "",
		},
	}

	for _, test := range testcases {
		fmt.Println(pkgname(test.input), "-->", test.output)
		assert.Equal(t, pkgname(test.input), test.output)
	}
}
