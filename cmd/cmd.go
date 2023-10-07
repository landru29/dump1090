package cmd

import (
	"github.com/landru29/dump1090/internal/application"
	"github.com/spf13/cobra"
)

func RootCommand() *cobra.Command {
	var (
		app    *application.App
		config application.Config
	)

	rootCommand := &cobra.Command{
		Use:   "dump1090",
		Short: "dump1090",
		Long:  "dump1090 main command",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			var err error

			app, err = application.New(&config)

			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.Start(cmd.Context())
		},
	}

	rootCommand.PersistentFlags().StringVarP(&config.FixturesFilename, "fixture-file", "", "", "Filename of the fixture data file")
	rootCommand.PersistentFlags().Uint32VarP(&config.DeviceIndex, "device", "d", 0, "Device index")
	rootCommand.PersistentFlags().BoolVarP(&config.EnableAGC, "enable-agc", "a", false, "Enable AGC")
	rootCommand.PersistentFlags().Uint32VarP(&config.Frequency, "frequency", "f", 1090000000, "frequency in Hz")
	rootCommand.PersistentFlags().IntVarP(&config.Gain, "gain", "g", 0, "gain")

	return rootCommand
}
