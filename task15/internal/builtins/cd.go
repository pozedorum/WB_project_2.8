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

func (cdu *CdUtil) Execute(args []string, env core.Environment, w io.Writer) error {
	return nil
}
