package integration

import (
	"fmt"
	"github.com/google/uuid"
)

// EventMessage represents a message that contains the event to be sent or received (outbox pattern).
type EventMessage struct {
	// ID unique id of each message.
	ID string `json:"id"`

	// AggregateType the type of the aggregate root to which a given event is related.
	AggregateType string `json:"aggregateType"`

	// AggregateID the id of the aggregate root that is affected by a given event.
	AggregateID string `json:"aggregateId"`

	// Type of event, e.g. "Order Created" or "Order Line Canceled".
	Type string `json:"eventType"`

	// Payload of the event.
	Payload []byte `json:"payload"`
}

func (e EventMessage) String() string {
	return fmt.Sprintf(`{"id": "%s", "aggregateType": "%s", "aggregateId": "%s", "eventType": "%s"}`,
		e.ID, e.AggregateType, e.AggregateType, e.Type)
}

// ExportedEvent interface must be implemented by structures that represent events (outbox pattern).
type ExportedEvent interface {
	// AggregateID the id of the aggregate root that is affected by a given event; this could for instance be the id
	// of a purchase order or a customer id; Similar to the aggregate type, events pertaining to a sub-entity contained
	// within an aggregate should use the id of the containing aggregate root, e.g. the purchase order id for an order
	// line cancelation event. This id will be used as the key for Kafka messages later on. That way, all events
	// pertaining to one aggregate root or any of its contained sub-entities will go into the same partition of that
	// Kafka topic, which ensures that consumers of that topic will consume all the events related to one and the same
	// aggregate in the exact order as they were produced.
	AggregateID() string

	// AggregateType the type of the aggregate root to which a given event is related; the idea being, leaning on the
	// same concept of domain-driven design, that exported events should refer to an aggregate ("a cluster of domain
	// objects that can be treated as a single unit"), where the aggregate root provides the sole entry point for
	// accessing any of the entities within the aggregate. This could for instance be "purchase order" or "customer".
	//
	// This value will be used to route events to corresponding topics in Kafka, so there’d be a topic for all events
	// related to purchase orders, one topic for all customer-related events etc. Note that also events pertaining to a
	// child entity contained within one such aggregate should use that same type. So e.g. an event representing the
	// cancelation of an individual order line (which is part of the purchase order aggregate) should also use the type
	// of its aggregate root, "order", ensuring that also this event will go into the "order" Kafka topic.
	AggregateType() string

	// Type the type of event, e.g. "Order Created" or "Order Line Canceled". Allows consumers to trigger suitable
	// event handlers.
	Type() string

	// Payload a JSON structure with the actual event contents, e.g. containing a purchase order, information about
	// the purchaser, contained order lines, their price etc.
	Payload() []byte
}

// BuildEventMessage takes ExportedEvent as a parameter and returns an EventMessage.
func BuildEventMessage(event ExportedEvent) EventMessage {
	return EventMessage{
		ID:            uuid.New().String(),
		AggregateType: event.AggregateType(),
		AggregateID:   event.AggregateID(),
		Type:          event.Type(),
		Payload:       event.Payload(),
	}
}
