// Package main is the main application.
package main

import (
	"log/slog"
	"time"

	"github.com/landru29/dump1090/internal/application"
	"github.com/landru29/dump1090/internal/database"
	"github.com/landru29/dump1090/internal/logger"
	"github.com/landru29/dump1090/internal/model"
	"github.com/landru29/dump1090/internal/processor"
	"github.com/landru29/dump1090/internal/processor/decoder"
	"github.com/landru29/dump1090/internal/serialize"
	"github.com/landru29/dump1090/internal/serialize/nmea"
	"github.com/landru29/dump1090/internal/transport/net"
	"github.com/spf13/cobra"
)

const (
	defaultNMEAmid                        = 226
	defaultFrequency                      = 1090000000
	defaultDatabaseLifetime time.Duration = time.Minute
)

func rootCommand() *cobra.Command { //nolint: funlen
	var (
		app                  *application.App
		config               application.Config
		httpConf             httpConfig
		transportScreen      string
		transportFile        string
		nmeaMid              uint16
		nmeaVessel           vessel = nmea.VesselTypeAircraft
		availableSerializers []serialize.Serializer
		loop                 bool
	)

	udpConf := net.NewProtocol("udp")
	tcpConf := net.NewProtocol("tcp")

	rootCommand := &cobra.Command{
		Use:   "dump1090",
		Short: "dump1090",
		Long:  "dump1090 main command",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error

			log := slog.New(slog.NewTextHandler(cmd.OutOrStdout(), nil))

			ctx := logger.WithLogger(cmd.Context(), log)

			cmd.SetContext(ctx)

			aircraftDB := database.NewElementStorage[model.ICAOAddr, model.Aircraft](
				ctx,
				database.ElementWithLifetime[model.ICAOAddr, model.Aircraft](config.DatabaseLifetime),
				database.ElementWithCleanCycle[model.ICAOAddr, model.Aircraft](config.DatabaseLifetime),
			)
			messageDB := database.NewChainedStorage[model.ICAOAddr, model.Squitter](
				ctx,
				database.ChainedWithLifetime[model.ICAOAddr, model.Squitter](config.DatabaseLifetime),
				database.ChainedWithCleanCycle[model.ICAOAddr, model.Squitter](config.DatabaseLifetime),
			)

			var serializers map[string]serialize.Serializer

			serializers, availableSerializers = provideSerializers(log, nmea.VesselType(nmeaVessel), nmeaMid)

			transporters, err := provideTransporters(
				ctx,
				log,
				availableSerializers,
				serializers,
				aircraftDB,
				httpConf,
				udpConf,
				tcpConf,
				transportScreen,
				transportFile,
			)
			if err != nil {
				return err
			}

			decoderCfg := []decoder.Configurator{
				decoder.WithDatabaseLifetime(config.DatabaseLifetime),
				// decoder.WithChecksumCheck(),
			}
			for _, transporter := range transporters {
				log.Info("loading transporter", "name", transporter.String())

				decoderCfg = append(decoderCfg, decoder.WithTransporter(transporter))
			}

			app, err = application.New(
				log,
				&config,
				[]processor.Processer{
					decoder.New(
						ctx,
						log,
						decoderCfg...,
					),
					// raw.New(log),
				},
				aircraftDB,
				messageDB,
			)

			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := app.Start(ctx); err != nil {
				return err
			}

			<-ctx.Done()

			return nil
		},
	}

	rootCommand.PersistentFlags().StringVarP(
		&config.FixturesFilename,
		"fixture-file",
		"",
		"",
		"Filename of the fixture data file",
	)

	rootCommand.PersistentFlags().Uint32VarP(
		&config.DeviceIndex,
		"device",
		"d",
		0,
		"Device index",
	)

	rootCommand.PersistentFlags().BoolVarP(
		&config.EnableAGC,
		"enable-agc",
		"a",
		false,
		"Enable AGC",
	)

	rootCommand.PersistentFlags().Uint32VarP(
		&config.Frequency,
		"frequency",
		"f",
		defaultFrequency,
		"frequency in Hz",
	)

	rootCommand.PersistentFlags().DurationVarP(
		&config.DatabaseLifetime,
		"db-lifetime",
		"",
		defaultDatabaseLifetime,
		"lifetime of elements in the AC database",
	)

	rootCommand.PersistentFlags().Float64VarP(
		&config.Gain,
		"gain",
		"g",
		0,
		"gain valid values are: 1.5, 4, 6.5, 9, 11.5, 14, 16.5, 19, 21.5, 24, 29, 34, 42, 43, 45, 47, 49",
	)

	rootCommand.PersistentFlags().VarP(
		&udpConf,
		"udp",
		"",
		"transmit data over udp (syntax: 'direction>format@host:port'; ie: --udp dial>json@192.168.1.10:8000)",
	)

	rootCommand.PersistentFlags().VarP(
		&tcpConf,
		"tcp",
		"",
		"transmit data over tcp (syntax: 'direction>format@host:port'; ie: --tcp bind>json@192.168.1.10:8000)",
	)

	rootCommand.PersistentFlags().VarP(
		&httpConf,
		"http",
		"",
		"transmit data over http (syntax: 'host:port/path'; ie: --http 0.0.0.0:8080/api)",
	)

	rootCommand.PersistentFlags().StringVarP(
		&transportScreen,
		"screen",
		"",
		"",
		"format to display output on the screen (json|nmea|text|none)",
	)

	rootCommand.PersistentFlags().VarP(
		&nmeaVessel,
		"nmea-vessel",
		"",
		"MMSI vessel (aircraft|helicopter)",
	)

	rootCommand.PersistentFlags().Uint16VarP(
		&nmeaMid,
		"nmea-mid",
		"",
		defaultNMEAmid,
		"MID (command 'mid' to list)",
	)

	rootCommand.PersistentFlags().BoolVarP(
		&loop,
		"loop",
		"",
		false,
		"With --fixture-file, read the same file in a loop",
	)

	rootCommand.PersistentFlags().StringVarP(
		&transportFile,
		"out-file",
		"",
		"",
		"format to display output on a file; ie --out-file nmea@/tmp/foo.txt",
	)

	rootCommand.AddCommand(&cobra.Command{
		Use: "mid",
		Run: func(cmd *cobra.Command, args []string) {
			for _, elt := range nmea.MidList {
				cmd.Printf("%d: %v %s\n", elt.MID, elt.Code, elt.Loc)
			}
		},
	})

	rootCommand.AddCommand(&cobra.Command{
		Use: "serializers",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("Available serializers:")

			for _, serializer := range availableSerializers {
				cmd.Printf(" - %s\n", serializer)
			}
		},
	})

	return rootCommand
}
