package telnet

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestSendToSocket(t *testing.T) {
	clientConn, serverConn := net.Pipe()
	defer serverConn.Close()

	input := strings.NewReader("hello\nworld\n")

	errCh := make(chan error, 1)

	go func() {
		errCh <- sendToSocket(clientConn, input)
		_ = clientConn.Close()
	}()

	got, err := io.ReadAll(serverConn)
	if err != nil {
		t.Fatalf("read from pipe: %v", err)
	}

	err = <-errCh
	if !errors.Is(err, errEOF) {
		t.Fatalf("got error %v, want errEOF", err)
	}

	want := "hello\nworld\n"

	if string(got) != want {
		t.Fatalf("got %q, want %q", string(got), want)
	}
}

func TestGetFromSocket(t *testing.T) {
	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()

	var output bytes.Buffer

	errCh := make(chan error, 1)

	go func() {
		errCh <- getFromSocket(clientConn, &output)
	}()

	go func() {
		_, _ = io.WriteString(serverConn, "hello\nworld\n")
		_ = serverConn.Close()
	}()

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("getFromSocket returned error: %v", err)
		}

	case <-time.After(time.Second):
		t.Fatal("getFromSocket did not finish")
	}

	want := "hello\nworld\n"

	if output.String() != want {
		t.Fatalf("got %q, want %q", output.String(), want)
	}
}

type errorWriter struct {
	err error
}

func (w errorWriter) Write([]byte) (int, error) {
	return 0, w.err
}

func TestGetFromSocketOutputError(t *testing.T) {
	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()
	defer serverConn.Close()

	wantErr := errors.New("output failed")

	go func() {
		_, _ = io.WriteString(serverConn, "hello\n")
	}()

	err := getFromSocket(clientConn, errorWriter{err: wantErr})

	if err == nil {
		t.Fatal("expected error")
	}

	if !errors.Is(err, wantErr) {
		t.Fatalf("got error %v, want %v", err, wantErr)
	}
}

func TestRunSessionReceivesServerData(t *testing.T) {
	clientConn, serverConn := net.Pipe()
	defer serverConn.Close()

	inputReader, inputWriter := io.Pipe()
	defer inputWriter.Close()

	var output bytes.Buffer

	go func() {
		_, _ = io.WriteString(serverConn, "server response\n")
		_ = serverConn.Close()
	}()

	err := runSession(
		context.Background(),
		clientConn,
		inputReader,
		&output,
	)
	if err != nil {
		t.Fatalf("runSession returned error: %v", err)
	}

	if output.String() != "server response\n" {
		t.Fatalf(
			"got %q, want %q",
			output.String(),
			"server response\n",
		)
	}
}

func TestRunSessionStopsOnContextCancellation(t *testing.T) {
	clientConn, serverConn := net.Pipe()
	defer serverConn.Close()

	inputReader, inputWriter := io.Pipe()
	defer inputWriter.Close()

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan error, 1)

	go func() {
		done <- runSession(
			ctx,
			clientConn,
			inputReader,
			io.Discard,
		)
	}()

	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("runSession returned error: %v", err)
		}

	case <-time.After(time.Second):
		t.Fatal("runSession did not stop after context cancellation")
	}
}

func TestRunTelnetClientClosesOnEOF(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer listener.Close()

	addr := listener.Addr().(*net.TCPAddr)

	received := make(chan string, 1)
	serverDone := make(chan error, 1)

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			serverDone <- err
			return
		}
		defer conn.Close()

		scanner := bufio.NewScanner(conn)

		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				serverDone <- err
				return
			}

			serverDone <- errors.New("client closed connection before sending data")
			return
		}

		received <- scanner.Text()
		
		if scanner.Scan() {
			serverDone <- errors.New("unexpected additional data")
			return
		}

		serverDone <- scanner.Err()
	}()

	input := strings.NewReader("hello\n")
	var output bytes.Buffer

	err = RunTelnetClient(
		addr.IP.String(),
		strconv.Itoa(addr.Port),
		time.Second,
		input,
		&output,
	)
	if err != nil {
		t.Fatalf("RunTelnetClient returned error: %v", err)
	}

	select {
	case got := <-received:
		if want := "hello"; got != want {
			t.Fatalf("server received %q, want %q", got, want)
		}
	case <-time.After(time.Second):
		t.Fatal("server did not receive client data")
	}

	select {
	case err := <-serverDone:
		if err != nil {
			t.Fatalf("server error: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("client did not close connection after EOF")
	}

	if output.Len() != 0 {
		t.Fatalf("unexpected output: %q", output.String())
	}
}

func TestRunTelnetClientConnectionError(t *testing.T) {
	var output bytes.Buffer

	err := RunTelnetClient(
		"127.0.0.1",
		"1",
		50*time.Millisecond,
		strings.NewReader(""),
		&output,
	)

	if err == nil {
		t.Fatal("expected connection error")
	}
}
