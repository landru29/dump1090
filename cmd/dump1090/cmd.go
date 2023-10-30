package main

import (
	"fmt"
	"strings"

	"github.com/landru29/dump1090/internal/application"
	nmeaencoder "github.com/landru29/dump1090/internal/nmea"
	"github.com/landru29/dump1090/internal/serialize"
	"github.com/landru29/dump1090/internal/serialize/basestation"
	"github.com/landru29/dump1090/internal/serialize/json"
	"github.com/landru29/dump1090/internal/serialize/nmea"
	"github.com/landru29/dump1090/internal/serialize/none"
	"github.com/landru29/dump1090/internal/serialize/text"
	"github.com/landru29/dump1090/internal/transport"
	"github.com/landru29/dump1090/internal/transport/file"
	"github.com/landru29/dump1090/internal/transport/http"
	"github.com/landru29/dump1090/internal/transport/net"
	"github.com/landru29/dump1090/internal/transport/screen"
	"github.com/spf13/cobra"
)

func rootCommand() *cobra.Command {
	var (
		app                  *application.App
		config               application.Config
		httpConf             httpConfig
		transportScreen      string
		transportFile        string
		nmeaMid              uint16
		nmeaVessel           string
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
			var (
				err error
			)

			serializers := map[string]serialize.Serializer{}

			vesselType, ok := map[string]nmea.VesselType{
				"aircraft":   nmea.VesselTypeAircraft,
				"helicopter": nmea.VesselTypeHelicopter,
			}[nmeaVessel]
			if !ok {
				return fmt.Errorf("unknow vessel type %s", nmeaVessel)
			}

			availableSerializers = []serialize.Serializer{
				none.Serializer{},
				json.Serializer{},
				text.Serializer{},
				basestation.Serializer{},
				nmea.New(vesselType, nmeaMid),
			}

			for _, serializer := range availableSerializers {
				serializers[serializer.String()] = serializer
			}

			transporters := []transport.Transporter{}

			if httpConf.addr != "" {
				httpTransport, err := http.New(cmd.Context(), httpConf.addr, httpConf.apiPath, availableSerializers)
				if err != nil {
					return err
				}
				transporters = append(transporters, httpTransport)

				cmd.Printf("API on http://%s%s\n", httpConf.addr, httpConf.apiPath)
			}

			if udpConf.IsValid() {
				udpTransport, err := net.New(cmd.Context(), serializers, udpConf, cmd.OutOrStdout())
				if err != nil {
					return err
				}
				transporters = append(transporters, udpTransport)
			}

			if tcpConf.IsValid() {
				tcpTransport, err := net.New(cmd.Context(), serializers, tcpConf, cmd.OutOrStdout())
				if err != nil {
					return err
				}
				transporters = append(transporters, tcpTransport)
			}

			if transportScreen != "" {
				screenTransport, err := screen.New(serializers[transportScreen])
				if err != nil {
					return err
				}

				transporters = append(transporters, screenTransport)
			}

			if transportFile != "" {
				splitter := strings.Split(transportFile, "@")
				if len(splitter) > 1 {
					fileTransport, err := file.New(cmd.Context(), strings.Join(splitter[1:], "@"), serializers[splitter[0]])
					if err != nil {
						return err
					}

					transporters = append(transporters, fileTransport)
				}
			}

			if len(transporters) == 0 {
				transporters = append(transporters, screen.Transporter{})
			}

			app, err = application.New(&config, transporters)

			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.Start(cmd.Context(), loop); err != nil {
				return err
			}

			<-cmd.Context().Done()

			return nil
		},
	}

	rootCommand.PersistentFlags().StringVarP(&config.FixturesFilename, "fixture-file", "", "", "Filename of the fixture data file")
	rootCommand.PersistentFlags().Uint32VarP(&config.DeviceIndex, "device", "d", 0, "Device index")
	rootCommand.PersistentFlags().BoolVarP(&config.EnableAGC, "enable-agc", "a", false, "Enable AGC")
	rootCommand.PersistentFlags().Uint32VarP(&config.Frequency, "frequency", "f", 1090000000, "frequency in Hz")
	rootCommand.PersistentFlags().IntVarP(&config.Gain, "gain", "g", 0, "gain")
	rootCommand.PersistentFlags().VarP(&udpConf, "udp", "", "transmit data over udp (syntax: 'direction>format@host:port'; ie: --udp dial>json@192.168.1.10:8000)")
	rootCommand.PersistentFlags().VarP(&tcpConf, "tcp", "", "transmit data over tcp (syntax: 'direction>format@host:port'; ie: --tcp bind>json@192.168.1.10:8000)")
	rootCommand.PersistentFlags().VarP(&httpConf, "http", "", "transmit data over http (syntax: 'host:port/path'; ie: --http 0.0.0.0:8080/api)")
	rootCommand.PersistentFlags().StringVarP(&transportScreen, "screen", "", "", "format to display output on the screen (json|nmea|text|none)")
	rootCommand.PersistentFlags().StringVarP(&nmeaVessel, "nmea-vessel", "", "aircraft", "MMSI vessel (aircraft|helicopter)")
	rootCommand.PersistentFlags().Uint16VarP(&nmeaMid, "nmea-mid", "", 226, "MID (command 'mid' to list)")
	rootCommand.PersistentFlags().BoolVarP(&loop, "loop", "", false, "With --fixture-file, read the same file in a loop")
	rootCommand.PersistentFlags().StringVarP(&transportFile, "out-file", "", "", "format to display output on a file; ie --out-file nmea@/tmp/foo.txt")

	rootCommand.AddCommand(&cobra.Command{
		Use: "mid",
		Run: func(cmd *cobra.Command, args []string) {
			for _, elt := range nmeaencoder.MidList {
				fmt.Printf("%d: %v %s\n", elt.MID, elt.Code, elt.Loc)
			}
		},
	})

	rootCommand.AddCommand(&cobra.Command{
		Use: "serializers",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Available serializers:")
			for _, serializer := range availableSerializers {
				fmt.Printf(" - %s\n", serializer)
			}
		},
	})

	return rootCommand
}
