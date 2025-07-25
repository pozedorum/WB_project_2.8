package builtins

import "io"

type KillUtil struct{}

func NewKillUtil() *KillUtil {
	return &KillUtil{}
}

func (kilu KillUtil) Name() string {
	return "kill"
}

func (kilu *KillUtil) Execute(args []string, env Environment, w io.Writer) error {
	return nil
}
