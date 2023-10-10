package cmd

import (
	"fmt"

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

	parser "github.com/adrianmo/go-nmea"
	ais "github.com/andmarios/aislib"
)

func RootCommand() *cobra.Command {
	var (
		app             *application.App
		config          application.Config
		outputFormat    string
		udpAddr         string
		httpAddr        string
		transportScreen bool
	)

	output := map[string]serialize.Serializer{
		"json": json.Serializer{},
		"text": text.Serializer{},
		"nmea": nmea.New(),
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

	rootCommand.AddCommand(&cobra.Command{
		Use: "foo",
		Run: func(cmd *cobra.Command, args []string) {
			s, _ := parser.Parse("!AIVDM,1,1,,A,15RTgt0PAso;90TKcjM8h6g208CQ,0*4A")
			out := s.(parser.VDMVDO)
			fmt.Println(out)
		},
	})

	rootCommand.AddCommand(&cobra.Command{
		Use: "bar",
		Run: func(cmd *cobra.Command, args []string) {
			send := make(chan string, 1024*8)
			receive := make(chan ais.Message, 1024*8)
			failed := make(chan ais.FailedSentence, 1024*8)

			done := make(chan bool)

			go ais.Router(send, receive, failed)

			go func() {
				var message ais.Message
				var problematic ais.FailedSentence
				for {
					select {
					case message = <-receive:
						switch message.Type {
						case 1, 2, 3:
							t, _ := ais.DecodeClassAPositionReport(message.Payload)
							fmt.Println(t)
						case 255:
							done <- true
						default:
							fmt.Printf("=== Message Type %2d ===\n", message.Type)
							fmt.Printf(" Unsupported type \n\n")
						}
					case problematic = <-failed:
						fmt.Println(problematic)
					}
				}
			}()

			send <- "!AIVDM,1,1,,A,15RTgt0PAso;90TKcjM8h6g208CQ,0*4A"
			close(send)
			<-done
		},
	})

	return rootCommand
}
