package main

import (
	"github.com/landru29/dump1090/internal/dump"

	"github.com/spf13/cobra"
)

func main() {
	rootCommand := &cobra.Command{
		Use:   "dump1090",
		Short: "dump1090",
		Long:  "dump1090 main command",
		RunE: func(cmd *cobra.Command, args []string) error {
			return dump.Start(
				0,
				0,
				0,
				true,
				"internal/dump/testdata/modes1.bin",
				func(msg *dump.Message) {},
				func(ac *dump.Aircraft) {
					cmd.Println(ac)
				},
			)
		},
	}

	if err := rootCommand.Execute(); err != nil {
		panic(err)
	}
}
