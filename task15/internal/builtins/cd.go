package builtins

import (
	"fmt"
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

func (cdu *CdUtil) Execute(args []string, env core.Environment, _ io.Reader, _ io.Writer) error {
	switch len(args) {
	case 0:
		home, err := env.GetHomeDir()
		if err != nil {
			return err
		}
		return env.Chdir(home)
	case 1:
		return env.Chdir(args[0])
	default:
		return fmt.Errorf("cd: string not in pwd: %s", args[0])
	}
}
