package telnet

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var errEOF = errors.New("EOF (Ctrl+D pressed)")

func RunTelnetClient(host string, port string, timeout time.Duration) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	setupSignalHandling(cancel)

	var d net.Dialer
	dialctx, dialCancel := context.WithTimeout(ctx, timeout)
	defer dialCancel()
	conn, err := d.DialContext(dialctx, "tcp", host+":"+port)
	if err != nil {
		return err
	}
	tcpConn := conn.(*net.TCPConn)
	defer conn.Close()
	if err := conn.SetDeadline(time.Now().Add(1 * time.Second)); err != nil {
		return err
	}
	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := conn.SetDeadline(time.Now().Add(1 * time.Second)); err != nil { // Новый дедлайн
					log.Fatal(err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	var wg sync.WaitGroup
	wg.Add(2)
	errCh := make(chan error, 2)

	go func() {
		if err := sendToSocket(ctx, tcpConn); err != nil {
			if errors.Is(err, errEOF) {
				errCh <- nil
			} else {
				errCh <- fmt.Errorf("send error: %w", err)
			}

			cancel()
		}
		wg.Done()
	}()

	go func() {
		if err := getFromSocket(ctx, tcpConn); err != nil {
			errCh <- fmt.Errorf("read error: %w", err)
			cancel()
		}
		wg.Done()
	}()
	wg.Wait()
	close(errCh)
	return <-errCh
}

func getFromSocket(ctx context.Context, conn *net.TCPConn) error {

	scanner := bufio.NewScanner(conn)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("\nStopping socket reader...")
			return nil
		default:
			if !scanner.Scan() {
				if errors.Is(scanner.Err(), os.ErrDeadlineExceeded) {
					// Дедлайн
					return ctx.Err()
				}
				if err := scanner.Err(); err != nil {
					return err
				} else {
					return errors.New("server closed connection")
				}
			}
			text := scanner.Text()
			//fmt.Println("Server says:", text)
			fmt.Println(text)
		}
	}
}

func sendToSocket(ctx context.Context, conn net.Conn) error {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("\nStopping socket writer...")
			return nil
		default:
			if !scanner.Scan() {
				return errEOF
			}
			text := scanner.Text() + "\n"
			if _, err := conn.Write([]byte(text)); err != nil {

				return err
			}
		}

	}
}

func setupSignalHandling(cancel context.CancelFunc) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM)
	go func() {
		<-sig
		cancel()
	}()
}
