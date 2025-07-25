package builtins

import (
	"io"

	"task15/internal/core"
)

type KillUtil struct{}

func NewKillUtil() *KillUtil {
	return &KillUtil{}
}

func (kilu KillUtil) Name() string {
	return "kill"
}

func (kilu *KillUtil) Execute(args []string, env core.Environment, w io.Writer) error {
	return nil
}
