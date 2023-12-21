// Package main is the TCP application test.
package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/spf13/cobra"
)

const (
	defaultTCPport = 2000
	bufferSize     = 2048
	maxRetries     = 300
)

func main() { //nolint: funlen,gocognit
	var (
		port    uint32
		address string
	)

	rootCommand := &cobra.Command{
		Use:   "tcp",
		Short: "tcp",
		Long:  "tcp main command",
	}

	rootCommand.AddCommand(&cobra.Command{
		Use:   "bind",
		Short: "bind",
		Long:  "bind tcp addr",
		RunE: func(cmd *cobra.Command, args []string) error {
			tcpServer, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
			if err != nil {
				return err
			}
			defer func() {
				_ = tcpServer.Close()
			}()

			cmd.Printf("Listening on port %d\n", port)

			go func() {
				for {
					// Wait for a connection.
					conn, err := tcpServer.Accept()
					if err != nil {
						cmd.PrintErr(err)

						return
					}

					defer func(c net.Conn) {
						_ = c.Close()
					}(conn)

					go func(c net.Conn) {
						for {
							_, _ = c.Write([]byte("foo"))
							time.Sleep(time.Second)
						}
					}(conn)
				}
			}()

			<-cmd.Context().Done()
			cmd.Println("Quitting")

			return nil
		},
	})

	rootCommand.AddCommand(&cobra.Command{
		Use:   "dial",
		Short: "dial",
		Long:  "dial tcp addr",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				dialer net.Dialer
				conn   net.Conn
			)

			cmd.Printf("trying to connect to %s ...\n", address)

			bckoff := backoff.WithMaxRetries(backoff.NewConstantBackOff(1*time.Second), maxRetries)
			err := backoff.Retry(func() error {
				var err error
				conn, err = dialer.DialContext(cmd.Context(), "tcp", address)
				if err != nil {
					return err
				}

				return nil
			}, bckoff)
			if err != nil {
				return err
			}

			defer func() {
				_ = conn.Close()
			}()

			cmd.Printf("connected to %s\n", address)

			errChan := make(chan error)

			go func() {
				var cnt int
				packet := make([]byte, bufferSize)
				for {
					cnt, err = conn.Read(packet)
					if err != nil {
						cmd.PrintErr(err)

						errChan <- err

						return
					}

					cmd.Printf("%d =>%s\n", cnt, string(packet[:cnt]))
				}
			}()

			for {
				select {
				case <-cmd.Context().Done():
					cmd.Println("Quitting")

					return nil
				case err := <-errChan:
					cmd.Println("error occurred")

					return err
				}
			}
		},
	})

	rootCommand.PersistentFlags().Uint32VarP(&port, "port", "p", defaultTCPport, "port to bind")
	rootCommand.PersistentFlags().StringVarP(&address, "addr", "a", "127.0.0.1:3000", "address to dial")

	osSignal := make(chan os.Signal, 1)

	// add any other syscalls that you want to be notified with
	signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()

	go func() {
		<-osSignal

		cancel()
	}()

	if err := rootCommand.ExecuteContext(ctx); err != nil {
		fmt.Println(err) //nolint: forbidigo

		os.Exit(1) //nolint: gocritic
	}
}
