package main

import (
	"fmt"
	"os"

	"github.com/parkplusplus/cli/mq"
	"github.com/parkplusplus/cli/sb"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{}

	rootCmd.AddCommand(
		sb.NewCommand(),
		mq.NewMetricCommand(),
		mq.NewLogCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		os.Exit(1)
	}
}
