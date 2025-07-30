package executor

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"

	"task15/internal/builtins"
	"task15/internal/core"
)

var testMode = false

type Executor struct {
	builtins *builtins.Registry
	env      core.Environment
	stdin    io.Reader
	stdout   io.Writer
	stderr   io.Writer

	closers []io.Closer

	procMutex   sync.Mutex
	currentProc *os.Process
}

// NewExecutor — Базовый конструктор
func NewExecutor(
	builtins *builtins.Registry,
	env core.Environment,
	stdin io.Reader,
	stdout io.Writer,
) *Executor {
	return &Executor{
		builtins: builtins,
		env:      env,
		stdin:    stdin,
		stdout:   stdout,
		stderr:   os.Stderr,
	}
}

// NewDefaultExecutor — Упрощённый конструктор
func NewDefaultExecutor() *Executor {
	return NewExecutor(
		builtins.NewRegistryWithDefaults(),
		&core.DefaultEnvironment{},
		os.Stdin,
		os.Stdout,
	)
}

func (e *Executor) Execute(ctx context.Context, cmd *core.Command) error {
	if cmd == nil {
		return errors.New("error: No command")
	}
	// Установка редиректов
	if err := e.setupRedirects(ctx, cmd); err != nil {
		return err
	}

	// Обработка pipe/условий
	if cmd.PipeTo != nil {
		return e.runPipeline(ctx, cmd)
	}

	return e.runCommand(ctx, cmd)
}

func (e *Executor) setupRedirects(ctx context.Context, cmd *core.Command) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	e.Close()
	e.stdin, e.stdout = os.Stdin, os.Stdout

	for _, redirect := range cmd.Redirects {
		if err := ctx.Err(); err != nil {
			return err
		}

		switch redirect.Type {
		case "<":
			file, err := os.Open(redirect.File)
			if err != nil {
				return fmt.Errorf("error: cannot open input file: %w", err)
			}
			e.closers = append(e.closers, file)
			e.stdin = file
		case ">":
			file, err := os.OpenFile(redirect.File, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
			if err != nil {
				return fmt.Errorf("error: cannot open output file: %w", err)
			}
			e.closers = append(e.closers, file)
			e.stdout = file
		default:
			return fmt.Errorf("error: this type of redirect is not supported: %s", redirect.Type)
		}
	}
	return nil
}

func (e *Executor) runPipeline(ctx context.Context, cmd *core.Command) error {
	pr, pw, err := os.Pipe()
	if err != nil {
		return err
	}
	defer pr.Close()
	// Первая команда пишет в пайп
	firstExec := NewExecutor(e.builtins, e.env, e.stdin, pw)
	errCh := make(chan error, 1)

	go func() {
		defer pw.Close()
		errCh <- firstExec.runCommand(ctx, cmd)
	}()

	// Вторая команда читает из пайпа
	secondExec := NewExecutor(e.builtins, e.env, pr, e.stdout)
	err = secondExec.runCommand(ctx, cmd.PipeTo)

	// Дожидаемся завершения первой команды
	if firstErr := <-errCh; firstErr != nil {
		return firstErr
	}

	return err
}

func (e *Executor) runCommand(ctx context.Context, cmd *core.Command) error {
	//log.Printf("runCommand: %s %v\n", cmd.Name, cmd.Args)
	//log.Printf("  stdin: %T, stdout: %T\n", e.stdin, e.stdout)

	// Для встроенных команд
	if bCmd, ok := e.builtins.GetCommand(cmd.Name); ok {
		//log.Println("Executing builtin command")
		return bCmd.Execute(cmd.Args, e.env, e.stdin, e.stdout)
	}

	// Проверка доступности команды
	if !testMode {
		//log.Println("Checking command availability...")
		if _, err := exec.LookPath(cmd.Name); err != nil {
			//log.Printf("Command not found: %v\n", err)
			return fmt.Errorf("command not found: %w", err)
		}
	}

	proc := exec.CommandContext(ctx, cmd.Name, cmd.Args...)
	proc.Stdin = e.stdin
	proc.Stdout = e.stdout
	proc.Stderr = e.stderr

	//log.Println("Starting external command...")
	if err := proc.Start(); err != nil {
		//log.Printf("Command start failed: %v\n", err)
		return fmt.Errorf("start failed: %w", err)
	}

	//log.Println("Waiting for command to complete...")
	err := proc.Wait()
	//log.Printf("Command completed with error: %v\n", err)
	return err
}

func (e *Executor) Interrupt() {
	e.procMutex.Lock()
	defer e.procMutex.Unlock()

	if e.currentProc != nil {
		e.currentProc.Signal(os.Interrupt)
	}
}

func (e *Executor) Close() error {
	e.procMutex.Lock()
	defer e.procMutex.Unlock()

	for _, closer := range e.closers {
		if err := closer.Close(); err != nil {
			return err
		}
	}
	e.closers = nil
	return nil
}
