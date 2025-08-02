package executor

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"task15/internal/core"
)

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		panic(err)
	}

	return buf.String()
}

func TestEchoCommand(t *testing.T) {
	ctx := context.Background()
	e := NewDefaultExecutor()

	output := captureOutput(func() {
		cmd := &core.Command{
			Name: "echo",
			Args: []string{"hello", "world"},
		}
		if err := e.Execute(ctx, cmd); err != nil {
			t.Fatalf("Execute failed: %v", err)
		}
	})

	expected := "hello world\n"
	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}
}

func TestCdCommand(t *testing.T) {
	ctx := context.Background()
	originalDir, _ := os.Getwd()
	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			panic(err)
		}

	}()
	testDir, err := os.MkdirTemp("", "testcd")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(testDir)

	e := NewDefaultExecutor()

	cmd := &core.Command{
		Name: "cd",
		Args: []string{testDir},
	}

	err = e.Execute(ctx, cmd)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	currentDir, _ := os.Getwd()
	if !strings.HasSuffix(currentDir, filepath.Base(testDir)) {
		t.Errorf("Expected dir to end with %q, got %q", filepath.Base(testDir), currentDir)
	}
}

func TestPwdCommand(t *testing.T) {
	ctx := context.Background()
	e := NewDefaultExecutor()

	output := captureOutput(func() {
		cmd := &core.Command{Name: "pwd"}
		if err := e.Execute(ctx, cmd); err != nil {
			t.Fatalf("Execute failed: %v", err)
		}
	})

	expected, _ := os.Getwd()
	if !strings.Contains(output, expected) {
		t.Errorf("Expected output to contain %q, got %q", expected, output)
	}
}

func TestPsCommand(t *testing.T) {
	ctx := context.Background()
	e := NewDefaultExecutor()

	output := captureOutput(func() {
		cmd := &core.Command{Name: "ps"}
		if err := e.Execute(ctx, cmd); err != nil {
			t.Fatalf("Execute failed: %v", err)
		}
	})

	if len(output) == 0 {
		t.Error("Expected process list, got empty output")
	}
	if !strings.Contains(output, "PID") {
		t.Error("Expected PS header, got:", output)
	}
}

func TestKillCommand(t *testing.T) {
	ctx := context.Background()
	e := NewDefaultExecutor()

	testCmd := exec.Command("sleep", "1")
	if err := testCmd.Start(); err != nil {
		t.Skipf("Cannot start test process: %v", err)
	}
	defer func() {
		if err := testCmd.Process.Kill(); err != nil {
			panic(err)
		}
	}()
	output := captureOutput(func() {
		cmd := &core.Command{
			Name: "kill",
			Args: []string{strconv.Itoa(testCmd.Process.Pid)},
		}
		if err := e.Execute(ctx, cmd); err != nil {
			t.Fatalf("Execute failed: %v", err)
		}
	})

	if !strings.Contains(output, strconv.Itoa(testCmd.Process.Pid)) {
		t.Errorf("Expected output to contain PID %d, got %q", testCmd.Process.Pid, output)
	}
}

func TestExecutorClosesResources(t *testing.T) {
	e := NewDefaultExecutor()

	tmpFile, err := os.CreateTemp("", "test_cleanup")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	e.closers = append(e.closers, tmpFile)

	output := captureOutput(func() {
		if err = e.Close(); err != nil {
			t.Errorf("Close() failed: %v", err)
		}
	})

	_, err = tmpFile.Write([]byte("test"))
	if err == nil {
		t.Error("File should be closed after executor cleanup")
	}
	if output != "" {
		t.Errorf("Expected no output from Close(), got %q", output)
	}
}

func TestCommandWithTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Millisecond)
	defer cancel()

	e := NewDefaultExecutor()

	output := captureOutput(func() {

		cmd := &core.Command{
			Name: "sleep",
			Args: []string{"5"},
		}

		start := time.Now()
		err := e.Execute(ctx, cmd)
		duration := time.Since(start)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		if duration >= 5*time.Second {
			t.Errorf("Expected command to be interrupted, but it ran for full duration: %v", duration)
		}

		// Допустимые варианты ошибки:
		expectedErrors := []string{
			"context deadline exceeded", // Ожидаемая ошибка контекста
			"signal: killed",            // Альтернативный вариант (процесс убит по сигналу)
		}

		found := false
		for _, expected := range expectedErrors {
			if strings.Contains(err.Error(), expected) {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("Expected error to contain one of %v, got: %v", expectedErrors, err)
		}
	})

	if output != "" {
		t.Errorf("Expected no output for timeout test, got %q", output)
	}
}

func TestExecuteWithCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	e := NewDefaultExecutor()

	output := captureOutput(func() {
		cmd := &core.Command{Name: "sleep", Args: []string{"1"}}

		err := e.Execute(ctx, cmd)
		if err == nil {
			t.Fatal("Expected error for cancelled context, got nil")
		}

		if !errors.Is(err, context.Canceled) &&
			!strings.Contains(err.Error(), "context canceled") {
			t.Errorf("Expected context canceled error, got: %v", err)
		}
	})

	if output != "" {
		t.Errorf("Expected no output for cancelled context test, got %q", output)
	}
}

func TestExecuteInvalidCommand(t *testing.T) {
	e := NewDefaultExecutor()

	output := captureOutput(func() {
		cmd := &core.Command{Name: "nonexistent_command_123", Args: []string{}}

		err := e.Execute(context.Background(), cmd)
		if err == nil {
			t.Fatal("Expected error for invalid command, got nil")
		}

		if !strings.Contains(err.Error(), "not found") {
			t.Errorf("Expected 'not found' error, got: %v", err)
		}
	})

	if output != "" {
		t.Errorf("Expected no output for invalid command test, got %q", output)
	}
}

func TestExecuteEmptyCommand(t *testing.T) {
	e := NewDefaultExecutor()

	tests := []struct {
		name    string
		cmd     *core.Command
		wantErr bool
	}{
		{"Empty name", &core.Command{Name: ""}, true},
		{"Nil command", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				err := e.Execute(context.Background(), tt.cmd)
				if (err != nil) != tt.wantErr {
					t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				}
			})

			if output != "" {
				t.Errorf("Expected no output for empty command test, got %q", output)
			}
		})
	}
}

func TestRedirectOutputToFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_redirect_out")
	if err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	e := NewDefaultExecutor()

	output := captureOutput(func() {
		cmd := &core.Command{
			Name: "echo",
			Args: []string{"hello redirect"},
			Redirects: []core.Redirect{
				{Type: ">", File: tmpFile.Name()},
			},
		}

		err = e.Execute(context.Background(), cmd)
		if err != nil {
			t.Fatalf("Execute failed: %v", err)
		}
	})

	if output != "" {
		t.Errorf("Expected no output when redirecting to file, got %q", output)
	}

	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	expected := "hello redirect\n"
	if string(content) != expected {
		t.Errorf("Expected %q in file, got %q", expected, string(content))
	}
}

func TestRedirectInputFromFile(t *testing.T) {

	tmpFile, err := os.CreateTemp("", "test_input")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString("file content"); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	e := NewDefaultExecutor()

	output := captureOutput(func() {
		cmd := &core.Command{
			Name: "cat",
			Redirects: []core.Redirect{
				{Type: "<", File: tmpFile.Name()},
			},
		}

		if err := e.Execute(context.Background(), cmd); err != nil {
			t.Fatalf("Execute failed: %v", err)
		}
	})

	if output != "file content" {
		t.Errorf("Expected %q, got %q", "file content", output)
	}
}

func TestPipelineTwoCommands(t *testing.T) {
	e := NewDefaultExecutor()

	output := captureOutput(func() {
		cmd := &core.Command{
			Name: "echo",
			Args: []string{"hello pipeline"},
			PipeTo: &core.Command{
				Name: "tr",
				Args: []string{"a-z", "A-Z"},
			},
		}

		if err := e.Execute(context.Background(), cmd); err != nil {
			t.Fatalf("Execute failed: %v", err)
		}
	})

	expected := "HELLO PIPELINE\n"
	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}
}

func TestRedirectErrors(t *testing.T) {
	e := NewDefaultExecutor()

	tests := []struct {
		name     string
		redirect core.Redirect
		wantErr  string
	}{
		{
			"Invalid input file",
			core.Redirect{Type: "<", File: "/nonexistent/file"},
			"cannot open input file",
		},
		{
			"Invalid output file",
			core.Redirect{Type: ">", File: "/invalid/path/to/file"},
			"cannot open output file",
		},
		{
			"Invalid redirect type",
			core.Redirect{Type: ">>", File: "test.txt"},
			"error: unsupported redirect type: >>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				cmd := &core.Command{
					Name:      "echo",
					Args:      []string{"test"},
					Redirects: []core.Redirect{tt.redirect},
				}

				err := e.Execute(context.Background(), cmd)
				if err == nil {
					t.Fatal("Expected error, got nil")
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("Expected error containing %q, got %q", tt.wantErr, err.Error())
				}
			})

			if output != "" {
				t.Errorf("Expected no output for redirect error test, got %q", output)
			}
		})
	}
}
