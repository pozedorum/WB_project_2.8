package builtins

import "io"

type PsUtil struct{}

func NewPsUtil() *PsUtil {
	return &PsUtil{}
}

func (psu PsUtil) Name() string {
	return "ps"
}

func (psu *PsUtil) Execute(args []string, env Environment, w io.Writer) error {
	return nil
}
