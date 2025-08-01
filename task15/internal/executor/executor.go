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

func NewDefaultExecutor() *Executor {
	return NewExecutor(
		builtins.NewRegistryWithDefaults(),
		&core.DefaultEnvironment{},
		os.Stdin,
		os.Stdout,
	)
}

func (e *Executor) Execute(ctx context.Context, cmd *core.Command) error {
	if cmd == nil || cmd.Name == "" {
		return errors.New("error: No command")
	}
	var err error
	// Установка редиректов
	if err = e.setupRedirects(ctx, cmd); err != nil {
		return err
	}

	// Обработка команд с управляющими символами
	switch {
	case cmd.PipeTo != nil:
		return e.handlePipe(ctx, cmd)
	case cmd.AndNext != nil:
		if err = e.runCommand(ctx, cmd); err == nil { // Только если первая команда успешна
			return e.Execute(ctx, cmd.AndNext)
		}
		return err
	case cmd.OrNext != nil:
		if err := e.runCommand(ctx, cmd); err != nil { // Только если первая команда неуспешна
			return e.Execute(ctx, cmd.OrNext)
		}
		return nil
	default:
		return e.runCommand(ctx, cmd)
	}
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
			return fmt.Errorf("error: unsupported redirect type: %s", redirect.Type)
		}
	}
	return nil
}

func (e *Executor) handlePipe(ctx context.Context, cmd *core.Command) error {
	pr, pw, err := os.Pipe()
	if err != nil {
		return fmt.Errorf("pipe error: %w", err)
	}

	errCh := make(chan error, 2)

	// Левый процесс (пишет в пайп)
	go func() {
		defer pw.Close() // Важно: закрываем pipe-writer после завершения команды
		tempExec := NewExecutor(e.builtins, e.env, e.stdin, pw)
		errCh <- tempExec.runCommand(ctx, cmd)
	}()

	// Правый процесс (читает из пайпа)
	go func() {
		defer pr.Close() // Закрываем pipe-reader
		tempExec := NewExecutor(e.builtins, e.env, pr, e.stdout)
		errCh <- tempExec.runCommand(ctx, cmd.PipeTo)
	}()

	// Ожидаем завершения
	var firstErr, secondErr error
	for i := 0; i < 2; i++ {
		if err := <-errCh; err != nil && firstErr == nil {
			firstErr = err
		} else if err != nil {
			secondErr = err
		}
	}

	if firstErr != nil {
		return fmt.Errorf("pipe command failed: %w", firstErr)
	}
	if secondErr != nil {
		return fmt.Errorf("pipe command failed: %w", secondErr)
	}
	return nil
}

func (e *Executor) runCommand(ctx context.Context, cmd *core.Command) error {
	if cmd.IsEmpty() {
		return nil
	}

	// Для встроенных команд
	if bCmd, ok := e.builtins.GetCommand(cmd.Name); ok {
		return bCmd.Execute(cmd.Args, e.env, e.stdin, e.stdout)
	}

	// Для внешних команд
	proc := exec.CommandContext(ctx, cmd.Name, cmd.Args...)
	proc.Stdin = e.stdin
	proc.Stdout = e.stdout
	proc.Stderr = e.stderr

	// Особые флаги для ps в пайплайнах
	if cmd.Name == "ps" && cmd.PipeTo != nil {
		proc.Args = append(proc.Args, "-x") // Для macOS
	}

	// Сохраняем процесс для прерывания
	e.procMutex.Lock()
	if err := proc.Start(); err != nil {
		e.procMutex.Unlock()
		return fmt.Errorf("start failed: %w", err)
	}
	e.currentProc = proc.Process
	e.procMutex.Unlock()

	defer func() {
		e.procMutex.Lock()
		e.currentProc = nil
		e.procMutex.Unlock()
	}()

	err := proc.Wait()
	if exitErr, ok := err.(*exec.ExitError); ok {
		// Для внешних команд возвращаем ошибку с кодом возврата
		return fmt.Errorf("exit status %d", exitErr.ExitCode())
	}

	return err
}

func (e *Executor) Interrupt() {
	e.procMutex.Lock()
	defer e.procMutex.Unlock()

	if e.currentProc != nil {
		if err := e.currentProc.Signal(os.Interrupt); err != nil {
			panic(err)
		}
	}
}

func (e *Executor) Close() error {
	e.procMutex.Lock()
	defer e.procMutex.Unlock()

	var errs []error
	for _, closer := range e.closers {
		if err := closer.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	e.closers = nil

	if len(errs) > 0 {
		return fmt.Errorf("errors while closing: %v", errs)
	}
	return nil
}
