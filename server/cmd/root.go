package cmd

import "github.com/spf13/cobra"

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
