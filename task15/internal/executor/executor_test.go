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

	"task15/internal/builtins"
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
	io.Copy(&buf, r)
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
	defer os.Chdir(originalDir)

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
	defer testCmd.Process.Kill()

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

	// Добавляем файл в closers
	e.closers = append(e.closers, tmpFile)

	if err = e.Close(); err != nil {
		t.Errorf("Close() failed: %v", err)
	}

	// Проверяем что файл действительно закрыт
	_, err = tmpFile.Write([]byte("test"))
	if err == nil {
		t.Error("File should be closed after executor cleanup")
	}
}

func TestCommandWithTimeout(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Millisecond)
	defer cancel()

	e := NewExecutor(builtins.NewRegistryWithDefaults(), &core.DefaultEnvironment{}, os.Stdin, os.Stdout)

	// Используем команду, которая точно не завершится мгновенно
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

	// Проверяем, что команда завершилась ДО истечения 5 секунд (значит, прервана)
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
}

func TestExecuteWithCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Отменяем контекст сразу

	e := NewExecutor(builtins.NewRegistryWithDefaults(), &core.DefaultEnvironment{}, os.Stdin, os.Stdout)
	cmd := &core.Command{Name: "sleep", Args: []string{"1"}}

	err := e.Execute(ctx, cmd)
	if err == nil {
		t.Fatal("Expected error for cancelled context, got nil")
	}

	if !errors.Is(err, context.Canceled) &&
		!strings.Contains(err.Error(), "context canceled") {
		t.Errorf("Expected context canceled error, got: %v", err)
	}
}

func TestExecuteInvalidCommand(t *testing.T) {
	e := NewExecutor(builtins.NewRegistryWithDefaults(), &core.DefaultEnvironment{}, os.Stdin, os.Stdout)
	cmd := &core.Command{Name: "nonexistent_command_123", Args: []string{}}

	err := e.Execute(context.Background(), cmd)
	if err == nil {
		t.Fatal("Expected error for invalid command, got nil")
	}

	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("Expected 'not found' error, got: %v", err)
	}
}

func TestExecuteEmptyCommand(t *testing.T) {
	e := NewExecutor(builtins.NewRegistryWithDefaults(), &core.DefaultEnvironment{}, os.Stdin, os.Stdout)

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
			err := e.Execute(context.Background(), tt.cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRedirectOutputToFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_redirect_out")
	if err != nil {
		t.Fatal(err)
	}
	tmpFile.Close() // Закрываем, чтобы executor мог переоткрыть
	defer os.Remove(tmpFile.Name())

	e := NewExecutor(builtins.NewRegistryWithDefaults(), &core.DefaultEnvironment{}, os.Stdin, os.Stdout)

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
	// Создаем файл с данными
	tmpFile, err := os.CreateTemp("", "test_input")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	// Явно записываем и закрываем файл
	if _, err := tmpFile.WriteString("file content"); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	// Буфер для вывода
	output := &bytes.Buffer{}
	e := NewExecutor(builtins.NewRegistryWithDefaults(), &core.DefaultEnvironment{}, nil, output)

	cmd := &core.Command{
		Name: "cat",
		Redirects: []core.Redirect{
			{Type: "<", File: tmpFile.Name()},
		},
	}

	if err := e.Execute(context.Background(), cmd); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if output.String() != "file content" {
		t.Errorf("Expected %q, got %q", "file content", output.String())
	}
}

func TestPipelineTwoCommands(t *testing.T) {
	output := &bytes.Buffer{}
	e := NewExecutor(builtins.NewRegistryWithDefaults(), &core.DefaultEnvironment{}, nil, output)

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

	expected := "HELLO PIPELINE\n"
	if output.String() != expected {
		t.Errorf("Expected %q, got %q", expected, output.String())
	}
}

func TestRedirectErrors(t *testing.T) {
	e := NewExecutor(builtins.NewRegistryWithDefaults(), &core.DefaultEnvironment{}, os.Stdin, os.Stdout)

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
			"not supported", // Обновляем ожидаемую ошибку
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
	}
}

// func TestComplexPipelineWithRedirects(t *testing.T) {
// 	// Проверяем доступность команд
// 	for _, cmd := range []string{"printf", "grep", "wc"} {
// 		if _, err := exec.LookPath(cmd); err != nil {
// 			t.Skipf("Command %q not available: %v", cmd, err)
// 		}
// 	}
// 	// Создаем временный выходной файл
// 	tmpOut, err := os.CreateTemp("", "test_output_*")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	tmpOut.Close()
// 	defer os.Remove(tmpOut.Name())
// 	e := NewExecutor(builtins.NewRegistryWithDefaults(), &core.DefaultEnvironment{}, os.Stdin, os.Stdout)
// 	// Используем printf вместо echo для корректной обработки переносов строк
// 	cmd := &core.Command{
// 		Name: "printf",
// 		Args: []string{"%s\\n", "apple", "banana", "orange", "pear", "kiwi"},
// 		PipeTo: &core.Command{
// 			Name: "grep",
// 			Args: []string{"-E", "a$"}, // -E для расширенных regexp в BSD grep
// 			PipeTo: &core.Command{
// 				Name: "wc",
// 				Args: []string{"-l"},
// 				Redirects: []core.Redirect{
// 					{Type: ">", File: tmpOut.Name()},
// 				},
// 			},
// 		},
// 	}
// 	if err := e.Execute(context.Background(), cmd); err != nil {
// 		t.Fatalf("Pipeline failed: %v", err)
// 	}
// 	content, err := os.ReadFile(tmpOut.Name())
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	actual := strings.TrimSpace(string(content))
// 	expected := "3" // banana, orange, pear
// 	if actual != expected {
// 		t.Errorf("Expected wc -l output %q, got %q", expected, actual)
// 		// Дополнительная диагностика
// 		diagnostic := exec.Command("sh", "-c", "printf '%s\\n' apple banana orange pear kiwi | grep -E 'a$' | wc -l")
// 		out, _ := diagnostic.CombinedOutput()
// 		t.Logf("Diagnostic command output: %q", out)
// 	}
// }
