// Package main is the UDP application test.
package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

const (
	defaultUDPport    = 2000
	defaultBufferSize = 1024
)

func main() {
	var port uint32

	rootCommand := &cobra.Command{
		Use:   "udp",
		Short: "udp",
		Long:  "udp main command",
	}

	rootCommand.AddCommand(&cobra.Command{
		Use:   "bind",
		Short: "bind",
		Long:  "bind port",
		RunE: func(cmd *cobra.Command, args []string) error {
			udpServer, err := net.ListenPacket("udp", fmt.Sprintf(":%d", port))
			if err != nil {
				return err
			}
			defer func() {
				_ = udpServer.Close()
			}()

			cmd.Printf("Listening on port %d\n", port)

			go func() {
				for {
					buf := make([]byte, defaultBufferSize)
					length, _, err := udpServer.ReadFrom(buf)
					if err != nil {
						continue
					}

					cmd.Print(string(buf[:length]))
				}
			}()

			<-cmd.Context().Done()
			cmd.Println("Quitting")

			return nil
		},
	})

	rootCommand.PersistentFlags().Uint32VarP(&port, "port", "p", defaultUDPport, "port to bind")

	s := make(chan os.Signal, 1)

	// add any other syscalls that you want to be notified with
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()

	go func() {
		<-s

		cancel()
	}()

	if err := rootCommand.ExecuteContext(ctx); err != nil {
		fmt.Println(err) //nolint: forbidigo

		cancel()

		os.Exit(1) //nolint: gocritic
	}
}
