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

func (echu *EchoUtil) Execute(args []string, env core.Environment, stdin io.Reader, stdout io.Writer) error {
	return nil
}
