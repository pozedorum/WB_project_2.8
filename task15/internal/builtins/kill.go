package builtins

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"syscall"

	"task15/internal/core"
)

type KillUtil struct{}

func NewKillUtil() *KillUtil {
	return &KillUtil{}
}

func (kilu KillUtil) Name() string {
	return "kill"
}

func (kilu *KillUtil) Execute(args []string, env core.Environment, stdin io.Reader, stdout io.Writer) error {
	if len(args) < 1 {
		return errors.New("usage: kill [-s SIGNAL] PID")
	}

	var (
		pid int
		sig = syscall.SIGTERM
		err error
	)

	if args[0] == "-s" || args[0] == "--signal" {
		if len(args) < 3 {
			return errors.New("missing signal or pid")
		}
		sig, err = parseSignal(args[1])
		if err != nil {
			return err
		}
		pid, err = strconv.Atoi(args[2])
	} else {
		pid, err = strconv.Atoi(args[0])
	}

	if err != nil && pid <= 0 {
		return errors.New("invalid PID")
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("process %d not found", pid)
	}

	if err := process.Signal(sig); err != nil {
		return fmt.Errorf("failed to kill %d: %v", pid, err)
	}

	fmt.Fprintf(stdout, "Sent signal %s to process %d\n", sig, pid)
	return nil
}

func parseSignal(sigStr string) (syscall.Signal, error) {
	// Поддержка числовых сигналов
	if sigNum, err := strconv.Atoi(sigStr); err == nil {
		return syscall.Signal(sigNum), nil
	}

	// Поддержка символьных имён сигналов
	switch sigStr {
	case "SIGTERM", "TERM":
		return syscall.SIGTERM, nil
	case "SIGKILL", "KILL":
		return syscall.SIGKILL, nil
	case "SIGINT", "INT":
		return syscall.SIGINT, nil
	case "SIGHUP", "HUP":
		return syscall.SIGHUP, nil
	default:
		return 0, fmt.Errorf("unknown signal: %s", sigStr)
	}
}
