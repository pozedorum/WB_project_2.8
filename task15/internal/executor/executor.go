package executor

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"task15/internal/builtins"
	"task15/internal/core"
)

var testMode = false

type Executor struct {
	builtins  *builtins.Registry
	env       core.Environment
	processes []core.ProcessInfo // Для управления
	stdin     io.Reader
	stdout    io.Writer
	stderr    io.Writer

	closers []io.Closer
}

// NewExecutor — Базовый конструктор
func NewExecutor(
	builtins *builtins.Registry,
	env core.Environment,
	stdin io.Reader,
	stdout io.Writer,
) *Executor {
	return &Executor{
		builtins:  builtins,
		env:       env,
		processes: make([]core.ProcessInfo, 0),
		stdin:     stdin,
		stdout:    stdout,
		stderr:    os.Stderr,
	}
}

// NewDefaultExecutor — Упрощённый конструктор (для большинства случаев)
func NewDefaultExecutor() *Executor {
	return NewExecutor(
		builtins.NewRegistryWithDefaults(),
		&core.DefaultEnvironment{},
		os.Stdin,
		os.Stdout,
	)
}

func (e *Executor) Execute(cmd *core.Command) error {
	// Установка редиректов
	if err := e.setupRedirects(cmd); err != nil {
		return err
	}

	// Обработка pipe/условий
	if cmd.PipeTo != nil {
		return e.runPipeline(cmd)
	}

	// 2. Запуск процесса

	return e.runCommand(cmd)
}

func (e *Executor) setupRedirects(cmd *core.Command) error {
	e.Close()
	e.stdin, e.stdout = os.Stdin, os.Stdout

	for _, redirect := range cmd.Redirects {
		switch redirect.Type {
		case "<":
			file, err := os.Open(redirect.File)
			if err != nil {
				return fmt.Errorf("error: cannot open input file: %w", err)
			}
			defer file.Close()
			e.stdin = file
		case ">":
			file, err := os.OpenFile(redirect.File, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
			if err != nil {
				return fmt.Errorf("error: cannot open output file: %w", err)
			}
			defer file.Close()
			e.stdout = file
		default:
			return fmt.Errorf("error: this type of rederect is not processing: %s", redirect.Type)
		}
	}
	return nil
}

func (e *Executor) runPipeline(cmd *core.Command) error {
	// Создать pipe, запустить команды цепочкой

	pipeReader, pipeWriter, err := os.Pipe()
	if err != nil {
		return fmt.Errorf("error: cannot create a pipe: %w", err)
	}

	e.closers = append(e.closers, pipeReader)

	originalStdout := e.stdout

	// Если stdout не перенаправлен в файл - используем pipe
	outputRedirected := false
	for _, r := range cmd.Redirects {
		if r.Type == ">" || r.Type == ">>" {
			outputRedirected = true
			break
		}
	}

	if !outputRedirected {
		e.stdout = pipeWriter // Перенаправляем вывод в pipe
	} else {
		pipeWriter.Close()
	}

	errCh := make(chan error, 1)
	go func() {
		defer pipeWriter.Close()

		nextExecutor := NewExecutor(e.builtins, e.env, pipeReader, originalStdout)

		errCh <- nextExecutor.Execute(cmd.PipeTo) // передаём ошибку в канал
	}()

	// Запускаем текущую команду
	cmdErr := e.runCommand(cmd)

	// Ждём завершения следующей команды
	pipelineErr := <-errCh

	// Приоритет ошибки текущей команды
	if cmdErr != nil {
		return fmt.Errorf("pipeline stage failed: %w", cmdErr)
	}
	if pipelineErr != nil {
		return fmt.Errorf("pipeline command failed: %w", pipelineErr)
	}

	return nil
}

func (e *Executor) runCommand(cmd *core.Command) error {
	defer e.Close()

	if bCom, ok := e.builtins.GetCommand(cmd.Name); ok {
		return bCom.Execute(cmd.Args, e.env, e.stdin, e.stdout)
	}

	if !testMode {
		if _, err := exec.LookPath(cmd.Name); err != nil {
			return fmt.Errorf("command %s not found", cmd.Name)
		}
	}

	proc := exec.Command(cmd.Name, cmd.Args...)
	proc.Stdin = e.stdin
	proc.Stdout = e.stdout
	proc.Stderr = e.stderr
	proc.Env = append(os.Environ(), e.env.Environ()...)

	if err := proc.Start(); err != nil {
		return fmt.Errorf("failed to start a proc %s: %w", cmd.Name, err)
	}

	e.processes = append(e.processes, core.ProcessInfo{
		Pid:     proc.Process.Pid,
		Command: cmd,
	})

	return proc.Wait()
}

func (e *Executor) Close() {
	for _, closer := range e.closers {
		closer.Close()
	}
	e.closers = nil
}
