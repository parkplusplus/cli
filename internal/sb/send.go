package sb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/parkplusplus/cli/internal/auth"
	"github.com/spf13/cobra"
)

type sendArgs struct {
	queueOrTopic string
	isBody       bool
	timeout      time.Duration
	auth         *auth.Namespace
	cmd          *cobra.Command
}

func newSendCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send [queue/topic] [flags]",
		Short: "Sends a message to a queue or topic",
	}

	sendArgs := &sendArgs{
		auth: auth.NewNamespace(cmd.PersistentFlags(), defaultConnectionStringVar),
		cmd:  cmd,
	}

	cmd.Args = cobra.ExactArgs(1) // <queue or topic>
	cmd.PersistentFlags().BoolVar(&sendArgs.isBody, "body", true, "Causes the entire contents of stdin to be used as the message Body.")
	cmd.PersistentFlags().DurationVar(&sendArgs.timeout, "timeout", time.Minute, "Time to wait for a send to complete.")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		sendArgs.queueOrTopic = args[0]
		return sendCommand(sendArgs)
	}

	return cmd
}

func sendCommand(args *sendArgs) error {
	if args.queueOrTopic == "" {
		args.cmd.Flags().PrintDefaults()
		return fmt.Errorf("you need to specify a queue or topic to send the message to")
	}

	client, err := newClient(args.auth)

	if err != nil {
		return fmt.Errorf("failed to create a Service Bus client: %w", err)
	}

	sender, err := client.NewSender(args.queueOrTopic, nil)

	if err != nil {
		return fmt.Errorf("failed to create a Service Bus Sender: %w", err)
	}

	fmt.Fprintf(os.Stderr, "Reading message from stdin (ctrl+d to send)\n")

	// read the message in from stdin
	bytes, err := io.ReadAll(os.Stdin)

	if err != nil {
		return fmt.Errorf("failed to read message bytes from stdin: %w", err)
	}

	var sendableMessage *azservicebus.Message

	if args.isBody {
		// assume the payload is just what we read from stdin
		sendableMessage = &azservicebus.Message{
			Body: bytes,
		}
	} else {
		// try deserializing their message as a azservicebus.Message in azservicebus - this lets
		// them control more of the message, including things like the ApplicationProperties
		if err := json.Unmarshal(bytes, &sendableMessage); err != nil {
			return fmt.Errorf("failed to deserialize stdin as an azservicebus.Message (pass -body if you intended to send this as the message body, not as a full azservicebus.Message): %w", err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), args.timeout)
	defer cancel()

	fmt.Fprintf(os.Stderr, "Sending message to %s, waiting up to %s\n", args.queueOrTopic, args.timeout)

	err = sender.SendMessage(ctx, &azservicebus.Message{
		Body: bytes,
	}, nil)

	if errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("timed out sending message (timeout duration was %s): %w", args.timeout, err)
	}

	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
