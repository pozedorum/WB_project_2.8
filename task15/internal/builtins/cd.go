package builtins

import (
	"io"

	"task15/internal/core"
)

type CdUtil struct{}

func NewCdUtil() *CdUtil {
	return &CdUtil{}
}

func (cdu CdUtil) Name() string {
	return "cd"
}

func (cdu *CdUtil) Execute(args []string, env core.Environment, stdin io.Reader, stdout io.Writer) error {
	return nil
}
