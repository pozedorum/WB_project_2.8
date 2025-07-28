package executor

import (
	"bytes"
	"io"
	"os"
	"testing"

	"task15/internal/builtins"
	"task15/internal/core"
)

// TestCommand реализует BuiltinCommand для тестирования
type TestCommand struct {
	name    string
	execute func(args []string, stdin io.Reader, stdout io.Writer) error
}

func (t *TestCommand) Name() string { return t.name }
func (t *TestCommand) Execute(args []string, _ core.Environment, stdin io.Reader, stdout io.Writer) error {
	return t.execute(args, stdin, stdout)
}

func TestSimpleCommand(t *testing.T) {
	registry := builtins.NewRegistry()
	registry.Register(&TestCommand{
		name: "test",
		execute: func(args []string, _ io.Reader, stdout io.Writer) error {
			_, err := stdout.Write([]byte("test output\n"))
			return err
		},
	})

	e := NewExecutor(registry, &core.DefaultEnvironment{}, os.Stdin, os.Stdout)

	var buf bytes.Buffer
	e.stdout = &buf

	cmd := &core.Command{Name: "test"}
	err := e.Execute(cmd)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if buf.String() != "test output\n" {
		t.Errorf("Expected %q, got %q", "test output\n", buf.String())
	}
}

func TestOutputRedirect(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_output")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	registry := builtins.NewRegistry()
	registry.Register(&TestCommand{
		name: "test",
		execute: func(args []string, _ io.Reader, stdout io.Writer) error {
			_, err = stdout.Write([]byte("file content\n"))
			return err
		},
	})

	e := NewExecutor(registry, &core.DefaultEnvironment{}, os.Stdin, os.Stdout)

	cmd := &core.Command{
		Name: "test",
		Redirects: []core.Redirect{
			{Type: ">", File: tmpFile.Name()},
		},
	}

	err = e.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	if string(content) != "file content\n" {
		t.Errorf("Expected %q, got %q", "file content\n", content)
	}
}

func TestInputRedirect(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_input")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	testData := "test input data\n"
	err = os.WriteFile(tmpFile.Name(), []byte(testData), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	registry := builtins.NewRegistry()
	registry.Register(&TestCommand{
		name: "cat",
		execute: func(_ []string, stdin io.Reader, stdout io.Writer) error {
			_, err = io.Copy(stdout, stdin)
			return err
		},
	})

	e := NewExecutor(registry, &core.DefaultEnvironment{}, os.Stdin, os.Stdout)

	var buf bytes.Buffer
	e.stdout = &buf

	cmd := &core.Command{
		Name: "cat",
		Redirects: []core.Redirect{
			{Type: "<", File: tmpFile.Name()},
		},
	}

	err = e.Execute(cmd)
	if err != nil {
		t.Fatal(err)
	}

	if buf.String() != testData {
		t.Errorf("Expected %q, got %q", testData, buf.String())
	}
}
