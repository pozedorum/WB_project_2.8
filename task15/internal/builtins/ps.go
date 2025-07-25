package builtins

import (
	"io"

	"task15/internal/core"
)

type PsUtil struct{}

func NewPsUtil() *PsUtil {
	return &PsUtil{}
}

func (psu PsUtil) Name() string {
	return "ps"
}

func (psu *PsUtil) Execute(args []string, env core.Environment, w io.Writer) error {
	return nil
}
