package sb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/spf13/cobra"
)

type receiveArgs struct {
	oneLine bool
	timeout time.Duration
	count   int
	auth    *Auth

	queue string
	// or
	topic        string
	subscription string

	cmd *cobra.Command
}

type mode string

const (
	receiveAndDelete = "ReceiveAndDelete"
)

var modes = []string{
	receiveAndDelete,
}

func newReceiveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "receive [queue|topic subscription] [flags]",
		Short: "Receive messages from a queue or a topic and subscription.",
	}

	receiveArgs := &receiveArgs{
		auth: AddAuth(cmd.Flags()),
		cmd:  cmd,
	}

	cmd.Args = cobra.RangeArgs(1, 2)
	cmd.Flags().BoolVar(&receiveArgs.oneLine, "oneline", true, "Print each message as a single line.")
	cmd.Flags().DurationVarP(&receiveArgs.timeout, "timeout", "t", time.Minute, "Maximum time to wait for a single message to arrive.")
	cmd.Flags().IntVarP(&receiveArgs.count, "count", "c", 1, "Maximum number of messages to wait for.")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 {
			receiveArgs.queue = args[0]
		} else if len(args) == 2 {
			receiveArgs.topic, receiveArgs.subscription = args[0], args[1]
		}

		return receiveCommand(receiveArgs)
	}

	return cmd
}

func receiveCommand(args *receiveArgs) error {
	client, err := args.auth.NewClient()

	if err != nil {
		return fmt.Errorf("failed to create a Service Bus client: %w", err)
	}

	defer client.Close(context.Background())

	var receiver *azservicebus.Receiver

	if args.queue != "" {
		receiver, err = client.NewReceiverForQueue(args.queue, &azservicebus.ReceiverOptions{
			ReceiveMode: azservicebus.ReceiveModeReceiveAndDelete,
		})
	} else {
		receiver, err = client.NewReceiverForSubscription(args.topic, args.subscription, &azservicebus.ReceiverOptions{
			ReceiveMode: azservicebus.ReceiveModeReceiveAndDelete,
		})
	}

	defer receiver.Close(context.Background())

	ctx, cancel := context.WithTimeout(context.Background(), args.timeout)
	defer cancel()

	messages, err := receiver.ReceiveMessages(ctx, args.count, nil)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("no messages arrived within %s", args.timeout)
		}

		return fmt.Errorf("failed receiving messages: %w", err)
	}

	for _, m := range messages {
		var bytes []byte
		var err error

		if args.oneLine {
			bytes, err = json.Marshal(m)
		} else {
			bytes, err = json.MarshalIndent(m, "  ", "  ")
		}

		if err != nil {
			return fmt.Errorf("failed to JSON.Marshal message with MessageID %s: %w", m.MessageID, err)
		}

		fmt.Printf("%s\n", string(bytes))
	}

	return nil
}
