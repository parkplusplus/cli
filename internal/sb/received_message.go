package sb

import (
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

type JSONReceivedMessage struct {
	// ApplicationProperties can be used to store custom metadata for a message.
	ApplicationProperties map[string]interface{}

	// Body is the payload for a message.
	Body []byte

	// ContentType describes the payload of the message, with a descriptor following
	// the format of Content-Type, specified by RFC2045 (ex: "application/json").
	ContentType *string

	// CorrelationID allows an application to specify a context for the message for the purposes of
	// correlation, for example reflecting the MessageID of a message that is being
	// replied to.
	CorrelationID *string

	// DeadLetterErrorDescription is the description set when the message was dead-lettered.
	DeadLetterErrorDescription *string

	// DeadLetterReason is the reason set when the message was dead-lettered.
	DeadLetterReason *string

	// DeadLetterSource is the name of the queue or subscription this message was enqueued on
	// before it was dead-lettered.
	DeadLetterSource *string

	// DeliveryCount is number of times this message has been delivered.
	// This number is incremented when a message lock expires or if the message is explicitly abandoned
	// with Receiver.AbandonMessage.
	DeliveryCount uint32

	// EnqueuedSequenceNumber is the original sequence number assigned to a message, before it
	// was auto-forwarded.
	EnqueuedSequenceNumber *int64

	// EnqueuedTime is the UTC time when the message was accepted and stored by Service Bus.
	EnqueuedTime *time.Time

	// ExpiresAt is the time when this message will expire.
	//
	// This time is calculated by adding the TimeToLive property, set in the message that was sent, along  with the
	// EnqueuedTime of the message.
	ExpiresAt *time.Time

	// LockedUntil is the time when the lock expires for this message.
	// This can be extended by using Receiver.RenewMessageLock.
	LockedUntil *time.Time

	// LockToken is the lock token for a message received from a Receiver created with a receive mode of ReceiveModePeekLock.
	LockToken [16]byte

	// MessageID is an application-defined value that uniquely identifies
	// the message and its payload. The identifier is a free-form string.
	//
	// If enabled, the duplicate detection feature identifies and removes further submissions
	// of messages with the same MessageId.
	MessageID string

	// PartitionKey is used with a partitioned entity and enables assigning related messages
	// to the same internal partition. This ensures that the submission sequence order is correctly
	// recorded. The partition is chosen by a hash function in Service Bus and cannot be chosen
	// directly.
	//
	// For session-aware entities, the ReceivedMessage.SessionID overrides this value.
	PartitionKey *string

	// ReplyTo is an application-defined value specify a reply path to the receiver of the message. When
	// a sender expects a reply, it sets the value to the absolute or relative path of the queue or topic
	// it expects the reply to be sent to.
	ReplyTo *string

	// ReplyToSessionID augments the ReplyTo information and specifies which SessionId should
	// be set for the reply when sent to the reply entity.
	ReplyToSessionID *string

	// ScheduledEnqueueTime specifies a time when a message will be enqueued. The message is transferred
	// to the broker but will not available until the scheduled time.
	ScheduledEnqueueTime *time.Time

	// SequenceNumber is a unique number assigned to a message by Service Bus.
	SequenceNumber *int64

	// SessionID is used with session-aware entities and associates a message with an application-defined
	// session ID. Note that an empty string is a valid session identifier.
	// Messages with the same session identifier are subject to summary locking and enable
	// exact in-order processing and demultiplexing. For session-unaware entities, this value is ignored.
	SessionID *string

	// State represents the current state of the message (Active, Scheduled, Deferred).
	State azservicebus.MessageState

	// Subject enables an application to indicate the purpose of the message, similar to an email subject line.
	Subject *string

	// TimeToLive is the duration after which the message expires, starting from the instant the
	// message has been accepted and stored by the broker, found in the ReceivedMessage.EnqueuedTime
	// property.
	//
	// When not set explicitly, the assumed value is the DefaultTimeToLive for the queue or topic.
	// A message's TimeToLive cannot be longer than the entity's DefaultTimeToLive, and is silently
	// adjusted if it is.
	TimeToLive *time.Duration

	// To is reserved for future use in routing scenarios but is not currently used by Service Bus.
	// Applications can use this value to indicate the logical destination of the message.
	To *string
}

func newJSONReceivedMessage(m *azservicebus.ReceivedMessage) JSONReceivedMessage {
	return JSONReceivedMessage{
		ApplicationProperties:      m.ApplicationProperties,
		Body:                       m.Body,
		ContentType:                m.ContentType,
		CorrelationID:              m.CorrelationID,
		DeadLetterErrorDescription: m.DeadLetterErrorDescription,
		DeadLetterReason:           m.DeadLetterReason,
		DeadLetterSource:           m.DeadLetterSource,
		DeliveryCount:              m.DeliveryCount,
		EnqueuedSequenceNumber:     m.EnqueuedSequenceNumber,
		EnqueuedTime:               m.EnqueuedTime,
		ExpiresAt:                  m.ExpiresAt,
		LockedUntil:                m.LockedUntil,
		LockToken:                  m.LockToken,
		MessageID:                  m.MessageID,
		PartitionKey:               m.PartitionKey,
		ReplyTo:                    m.ReplyTo,
		ReplyToSessionID:           m.ReplyToSessionID,
		ScheduledEnqueueTime:       m.ScheduledEnqueueTime,
		SequenceNumber:             m.SequenceNumber,
		SessionID:                  m.SessionID,
		State:                      m.State,
		Subject:                    m.Subject,
		TimeToLive:                 m.TimeToLive,
		To:                         m.To,
	}
}
