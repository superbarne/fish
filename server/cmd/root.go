package cmd

import (
	"runtime/debug"

	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "fish",
		Run: func(cmd *cobra.Command, args []string) {
			serve()
		},
	}

	rootCmd.AddCommand(NewServeCmd())

	return rootCmd
}

func Execute() {
	NewRootCmd().Execute()
}

func gitCommit() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				return setting.Value
			}
		}
	}
	return "local"
}
