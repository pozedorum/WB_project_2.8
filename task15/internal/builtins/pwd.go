package builtins

import (
	"io"

	"task15/internal/core"
)

type PwdUtil struct{}

func NewPwdUtil() *PwdUtil {
	return &PwdUtil{}
}

func (pwdu PwdUtil) Name() string {
	return "pwd"
}

func (pwdu *PwdUtil) Execute(args []string, env core.Environment, stdin io.Reader, stdout io.Writer) error {
	return nil
}
