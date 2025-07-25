package builtins

import "io"

type CdUtil struct{}

func NewCdUtil() *CdUtil {
	return &CdUtil{}
}

func (cdu CdUtil) Name() string {
	return "cd"
}

func (cdu *CdUtil) Execute(args []string, env Environment, w io.Writer) error {
	return nil
}
