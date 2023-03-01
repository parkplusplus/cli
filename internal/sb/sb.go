package sb

import (
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/parkplusplus/cli/internal/auth"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "sb",
		Short: "Tools for working with Azure Service Bus",
	}

	rootCmd.AddCommand(newSendCommand())
	rootCmd.AddCommand(newReceiveCmd())

	return rootCmd
}

func newClient(ns *auth.Namespace) (*azservicebus.Client, error) {

	if ns.DAC != nil {
		return azservicebus.NewClient(ns.Namespace, ns.DAC, nil)
	}

	return azservicebus.NewClientFromConnectionString(ns.ConnectionString, nil)
}

const defaultConnectionStringVar = "SERVICEBUS_CONNECTION_STRING"
