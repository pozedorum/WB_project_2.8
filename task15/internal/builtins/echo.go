package builtins

import (
	"io"

	"task15/internal/core"
)

type EchoUtil struct{}

func NewEchoUtil() *EchoUtil {
	return &EchoUtil{}
}

func (echu EchoUtil) Name() string {
	return "echo"
}

func (echu *EchoUtil) Execute(args []string, env core.Environment, w io.Writer) error {
	return nil
}
