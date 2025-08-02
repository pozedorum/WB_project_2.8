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

func (pwdu *PwdUtil) Execute(_ []string, env core.Environment, _ io.Reader, stdout io.Writer) error {
	wd, err := env.Getwd()
	if err != nil {
		return err
	}
	_, err = stdout.Write([]byte(wd + "\n"))
	return err
}
