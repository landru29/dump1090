package cmd

import (
	"fmt"

	"github.com/landru29/dump1090/internal/application"
	nmeaencoder "github.com/landru29/dump1090/internal/nmea"
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
		app             *application.App
		config          application.Config
		outputFormat    string
		udpAddr         string
		httpAddr        string
		transportScreen bool
		nmeaMid         uint16
		nmeaVessel      string
	)

	rootCommand := &cobra.Command{
		Use:   "dump1090",
		Short: "dump1090",
		Long:  "dump1090 main command",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			var (
				err      error
				formater serialize.Serializer
			)

			switch outputFormat {
			case "json":
				formater = json.Serializer{}
			case "text":
				formater = text.Serializer{}
			case "nmea":
				vesselType, ok := map[string]nmea.VesselType{
					"aircraft":   nmea.VesselTypeAircraft,
					"helicopter": nmea.VesselTypeHelicopter,
				}[nmeaVessel]
				if !ok {
					return fmt.Errorf("unknow vessel type %s", nmeaVessel)
				}

				formater = nmea.New(vesselType, nmeaMid)
			default:
				return fmt.Errorf("unknown format: %s", outputFormat)
			}

			transporters := []transport.Transporter{}

			if httpAddr != "" {
				httpTransport, err := http.New(cmd.Context(), formater, httpAddr)
				if err != nil {
					return err
				}
				transporters = append(transporters, httpTransport)
			}

			if udpAddr != "" {
				udpTrqnsport, err := udp.New(cmd.Context(), formater, udpAddr)
				if err != nil {
					return err
				}
				transporters = append(transporters, udpTrqnsport)
			}

			if transportScreen {
				transporters = append(transporters, screen.New(formater))
			}

			if len(transporters) == 0 {
				transporters = append(transporters, screen.Transporter{})
			}

			app, err = application.New(&config, formater, transporters)

			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.Start(cmd.Context()); err != nil {
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
	rootCommand.PersistentFlags().StringVarP(&outputFormat, "format", "", "text", "format (text|json|nmea)")
	rootCommand.PersistentFlags().StringVarP(&udpAddr, "udp", "", "", "transmit data over udp (syntax: 'host:port'; ie: --udp 192.168.1.10:8000)")
	rootCommand.PersistentFlags().StringVarP(&httpAddr, "http", "", "", "transmit data over http (syntax: 'host:port'; ie: --udp 0.0.0.0:8080)")
	rootCommand.PersistentFlags().BoolVarP(&transportScreen, "screen", "", false, "Display output on the screen")
	rootCommand.PersistentFlags().StringVarP(&nmeaVessel, "nmea-vessel", "", "aircraft", "MMSI vessel (aircraft|helicopter)")
	rootCommand.PersistentFlags().Uint16VarP(&nmeaMid, "nmea-mid", "", 226, "MID (command 'mid' to list)")

	rootCommand.AddCommand(&cobra.Command{
		Use: "mid",
		Run: func(cmd *cobra.Command, args []string) {
			for _, elt := range nmeaencoder.MidList {
				fmt.Printf("%d: %v %s\n", elt.MID, elt.Code, elt.Loc)
			}
		},
	})

	return rootCommand
}
