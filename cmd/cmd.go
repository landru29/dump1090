package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/landru29/dump1090/internal/application"
	"github.com/landru29/dump1090/internal/serialize"
	"github.com/landru29/dump1090/internal/serialize/json"
	"github.com/landru29/dump1090/internal/serialize/nmea"
	"github.com/landru29/dump1090/internal/serialize/text"
	"github.com/landru29/dump1090/internal/transport"
	"github.com/landru29/dump1090/internal/transport/http"
	"github.com/landru29/dump1090/internal/transport/screen"
	"github.com/landru29/dump1090/internal/transport/udp"
	"github.com/spf13/cobra"
)

func RootCommand() *cobra.Command {
	var (
		app          *application.App
		config       application.Config
		outputFormat string
		transportLst []string
	)

	output := map[string]serialize.Serializer{
		"json": json.Serializer{},
		"text": text.Serializer{},
		"nmea": nmea.Serializer{},
	}

	rootCommand := &cobra.Command{
		Use:   "dump1090",
		Short: "dump1090",
		Long:  "dump1090 main command",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			var err error

			formater := output[outputFormat]
			if formater == nil {
				return fmt.Errorf("unknown format: %s", outputFormat)
			}

			transporters := []transport.Transporter{}

			for _, name := range transportLst {
				switch name {
				case "screen":
					transporters = append(transporters, screen.New(formater))
				case "udp":
					udpTrqnsport, err := udp.New(formater, 9000)
					if err != nil {
						return err
					}
					transporters = append(transporters, udpTrqnsport)
				case "http":
					httpTransport, err := http.New(cmd.Context(), formater, 8080)
					if err != nil {
						return err
					}
					transporters = append(transporters, httpTransport)
				}
			}

			if len(transporters) == 0 {
				transporters = append(transporters, screen.Transporter{})
			}

			app, err = application.New(&config, formater, transporters)

			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithCancel(cmd.Context())
			defer cancel()

			if err := app.Start(ctx); err != nil {
				return err
			}

			s := make(chan os.Signal, 1)

			signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
			<-s

			return nil
		},
	}

	rootCommand.PersistentFlags().StringVarP(&config.FixturesFilename, "fixture-file", "", "", "Filename of the fixture data file")
	rootCommand.PersistentFlags().Uint32VarP(&config.DeviceIndex, "device", "d", 0, "Device index")
	rootCommand.PersistentFlags().BoolVarP(&config.EnableAGC, "enable-agc", "a", false, "Enable AGC")
	rootCommand.PersistentFlags().Uint32VarP(&config.Frequency, "frequency", "f", 1090000000, "frequency in Hz")
	rootCommand.PersistentFlags().IntVarP(&config.Gain, "gain", "g", 0, "gain")
	rootCommand.PersistentFlags().StringVarP(&outputFormat, "format", "", "text", "format (text|json|nmea)")
	rootCommand.PersistentFlags().StringSliceVarP(&transportLst, "transport", "", []string{}, "Transporter (http/screen/udp)")

	return rootCommand
}
