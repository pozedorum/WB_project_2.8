package builtins

import (
	"io"
	"os"
	"os/exec"

	"task15/internal/core"
)

type PsUtil struct{}

func NewPsUtil() *PsUtil {
	return &PsUtil{}
}

func (psu PsUtil) Name() string {
	return "ps"
}

func (psu *PsUtil) Execute(args []string, env core.Environment, stdin io.Reader, stdout io.Writer) error {
	cmd := exec.Command("ps", args...)
	cmd.Stdout = stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
