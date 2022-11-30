package sb

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "sb",
	}

	rootCmd.AddCommand(newSendCommand())
	rootCmd.AddCommand(newReceiveCommand())

	return rootCmd
}
