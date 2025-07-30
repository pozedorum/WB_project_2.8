package shell

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"task15/internal/executor"
	"task15/internal/parser"
	"time"
)

type Shell struct {
	exec    *executor.Executor
	running bool
}

func NewShell() *Shell {
	return &Shell{
		exec:    executor.NewDefaultExecutor(),
		running: true,
	}
}

func (s *Shell) Run() {
	// Инициализация обработчиков
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigChan)

	reader := bufio.NewReader(os.Stdin)
	// Главный цикл
	for {
		select {
		case sig := <-sigChan:
			switch sig {
			case syscall.SIGINT:
				fmt.Fprintln(os.Stderr, "^C")
				s.exec.Interrupt()
			case syscall.SIGTERM:
				fmt.Println("\nTerminating shell...")
				return
			}
		default:
			if !s.readCommand(reader) {
				return // Завершаем shell при EOF
			}
		}
	}
}

func (s *Shell) readCommand(reader *bufio.Reader) bool {
	fmt.Print("> ")
	line, err := reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			fmt.Println("\nExiting...")
			return false
		}
		fmt.Fprintf(os.Stderr, "Read error: %v\n", err)
		return true
	}

	cmd, err := parser.ParseLine(line)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Parse error: %v\n", err)
		return true
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.exec.Execute(ctx, cmd); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			fmt.Println("Command timed out")
		}
		fmt.Fprintf(os.Stderr, "Execution error: %v\n", err)
	}
	return true
}
