package sb

import (
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/joho/godotenv"
	"github.com/spf13/pflag"
)

type Auth struct {
	namespace  string
	cs         string
	useEnvFile bool
}

func AddAuth(fs *pflag.FlagSet) *Auth {
	auth := &Auth{}

	fs.BoolVar(&auth.useEnvFile, "env", false, "Load an .env file from the current directory.")
	fs.StringVar(&auth.namespace, "namespace", "", "Namespace (assumes DefaultAzureCredential)")
	fs.StringVar(&auth.cs, "connection-string-name", "SERVICEBUS_CONNECTION_STRING", "Environment variable that contains a connection string")

	return auth
}

func (a *Auth) NewClient() (*azservicebus.Client, error) {
	if a.useEnvFile {
		if err := godotenv.Load(); err != nil {
			return nil, fmt.Errorf("failed to load .env file: %w", err)
		}
	}

	var client *azservicebus.Client

	if a.namespace != "" {
		dac, err := azidentity.NewDefaultAzureCredential(nil)

		if err != nil {
			return nil, err
		}

		client, err = azservicebus.NewClient(a.namespace, dac, nil)

		if err != nil {
			return nil, err
		}
	} else {
		cs := os.Getenv(a.cs)

		if cs == "" {
			return nil, fmt.Errorf("no connection string in environment variable %s", cs)
		}

		var err error
		client, err = azservicebus.NewClientFromConnectionString(cs, nil)

		if err != nil {
			return nil, err
		}
	}

	return client, nil
}
