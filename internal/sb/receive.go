package sb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/parkplusplus/cli/internal/auth"
	"github.com/spf13/cobra"
)

type receiveCmd struct {
	oneLine bool
	timeout time.Duration
	count   int
	auth    *auth.Namespace

	queue string
	// or
	topic        string
	subscription string

	*cobra.Command
}

func newReceiveCmd() *cobra.Command {
	rc := &receiveCmd{
		Command: &cobra.Command{
			Use:   "receive [queue|topic subscription] [flags]",
			Short: "Receive messages from a queue or a topic and subscription.",
		},
	}

	rc.auth = auth.NewNamespace(rc.Command.Flags(), defaultConnectionStringVar)

	rc.Args = cobra.RangeArgs(1, 2)
	rc.Flags().BoolVar(&rc.oneLine, "oneline", true, "Print each message as a single line.")
	rc.Flags().DurationVarP(&rc.timeout, "timeout", "t", time.Minute, "Maximum time to wait for a single message to arrive.")
	rc.Flags().IntVarP(&rc.count, "count", "c", 1, "Maximum number of messages to wait for.")

	rc.RunE = func(cmd *cobra.Command, args []string) error {
		if err := rc.auth.Apply(); err != nil {
			return err
		}

		if len(args) == 1 {
			rc.queue = args[0]
		} else if len(args) == 2 {
			rc.topic, rc.subscription = args[0], args[1]
		}

		return rc.run(context.Background())
	}

	return rc.Command
}

func (cmd *receiveCmd) run(ctx context.Context) error {
	client, err := newClient(cmd.auth)

	if err != nil {
		return fmt.Errorf("failed to create a Service Bus client: %w", err)
	}

	defer client.Close(ctx)

	var receiver *azservicebus.Receiver

	if cmd.queue != "" {
		receiver, err = client.NewReceiverForQueue(cmd.queue, &azservicebus.ReceiverOptions{
			ReceiveMode: azservicebus.ReceiveModeReceiveAndDelete,
		})
	} else {
		receiver, err = client.NewReceiverForSubscription(cmd.topic, cmd.subscription, &azservicebus.ReceiverOptions{
			ReceiveMode: azservicebus.ReceiveModeReceiveAndDelete,
		})
	}

	if err != nil {
		return err
	}

	defer receiver.Close(context.Background())

	ctx, cancel := context.WithTimeout(context.Background(), cmd.timeout)
	defer cancel()

	messages, err := receiver.ReceiveMessages(ctx, cmd.count, nil)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("no messages arrived within %s", cmd.timeout)
		}

		return fmt.Errorf("failed receiving messages: %w", err)
	}

	for _, m := range messages {
		var bytes []byte
		var err error

		newM := newJSONReceivedMessage(m)

		if cmd.oneLine {
			bytes, err = json.Marshal(newM)
		} else {
			bytes, err = json.MarshalIndent(newM, "  ", "  ")
		}

		if err != nil {
			return fmt.Errorf("failed to JSON.Marshal message with MessageID %s: %w", m.MessageID, err)
		}

		fmt.Printf("%s\n", string(bytes))
	}

	return nil
}
