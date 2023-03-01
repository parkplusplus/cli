package auth

import (
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/joho/godotenv"
	"github.com/spf13/pflag"
)

// Namespace holds authentication details for Service Bus and Event Hubs namespaces.
type Namespace struct {
	opts struct {
		csEnvVar   string
		useEnvFile bool
	}

	Namespace        string
	DAC              *azidentity.DefaultAzureCredential
	ConnectionString string
}

func NewNamespace(fs *pflag.FlagSet, defaultVarName string) *Namespace {
	auth := &Namespace{}

	fs.BoolVarP(&auth.opts.useEnvFile, "use-dotenv", "e", false, "Load an .env file from the current directory.")
	fs.StringVarP(&auth.Namespace, "namespace", "n", "", "Namespace (ex: name.servicebus.windows.net, only needed when using the DefaultAzureCredential.)")
	fs.StringVarP(&auth.opts.csEnvVar, "use-envvar", "v", defaultVarName, "Environment variable that contains a connection string.")

	return auth
}

func (ns *Namespace) Apply() error {
	if ns.opts.useEnvFile {
		if err := godotenv.Load(); err != nil {
			return fmt.Errorf("failed to load .env file: %w", err)
		}
	}

	if ns.Namespace != "" {
		dac, err := azidentity.NewDefaultAzureCredential(nil)

		if err != nil {
			return err
		}

		ns.DAC = dac
		return nil
	}

	// must be connection string authentication
	connectionString := os.Getenv(ns.opts.csEnvVar)

	if connectionString == "" {
		return fmt.Errorf("no connection string in environment variable %q", ns.opts.csEnvVar)
	}

	ns.ConnectionString = connectionString
	return nil
}
