package builtins

import "io"

type PwdUtil struct{}

func NewPwdUtil() *PwdUtil {
	return &PwdUtil{}
}

func (pwdu PwdUtil) Name() string {
	return "pwd"
}

func (pwdu *PwdUtil) Execute(args []string, env Environment, w io.Writer) error {
	return nil
}
