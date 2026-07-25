package telnet

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os/signal"
	"syscall"
	"time"
)

var errEOF = errors.New("EOF (Ctrl+D pressed)")

func RunTelnetClient(host, port string, timeout time.Duration, input io.Reader, output io.Writer) error {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGTERM,
		syscall.SIGINT,
	)
	defer stop()

	dialCtx, cancelDial := context.WithTimeout(ctx, timeout)
	defer cancelDial()

	var dialer net.Dialer

	conn, err := dialer.DialContext(
		dialCtx,
		"tcp",
		net.JoinHostPort(host, port),
	)
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}

	return runSession(ctx, conn, input, output)
}

func runSession(ctx context.Context, conn net.Conn, input io.Reader, output io.Writer) error {
	errCh := make(chan error, 2)

	defer conn.Close()

	go func() {
		errCh <- sendToSocket(conn, input)
	}()

	go func() {
		errCh <- getFromSocket(conn, output)
	}()

	select {
	case <-ctx.Done():
		return nil

	case err := <-errCh:

		if err == nil || errors.Is(err, errEOF) || errors.Is(err, net.ErrClosed) {
			return nil
		}
		return err

	}
}

func getFromSocket(conn net.Conn, output io.Writer) error {
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		if _, err := fmt.Fprintln(output, scanner.Text()); err != nil {
			return fmt.Errorf("write output: %w", err)
		}
	}

	if err := scanner.Err(); err != nil {
		if errors.Is(err, net.ErrClosed) {
			return nil
		}
		return fmt.Errorf("read socket: %w", err)
	}

	return nil
}

func sendToSocket(conn net.Conn, input io.Reader) error {
	scanner := bufio.NewScanner(input)

	for scanner.Scan() {
		if _, err := fmt.Fprintln(conn, scanner.Text()); err != nil {
			return fmt.Errorf("write socket: %w", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read input: %w", err)
	}

	return errEOF
}
