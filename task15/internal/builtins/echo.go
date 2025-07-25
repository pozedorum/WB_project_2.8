package builtins

import "io"

type EchoUtil struct{}

func NewEchoUtil() *EchoUtil {
	return &EchoUtil{}
}

func (echu EchoUtil) Name() string {
	return "echo"
}

func (echu *EchoUtil) Execute(args []string, env Environment, w io.Writer) error {
	return nil
}
