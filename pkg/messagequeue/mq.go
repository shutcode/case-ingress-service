package messagequeue

// MessageQueue defines the interface for a message queue producer.
// It provides a generic way to publish messages.
type MessageQueue interface {
	// Publish sends a message to the specified topic.
	// 'message' should be a byte slice, allowing flexibility in message format (e.g., JSON, Protobuf).
	Publish(topic string, message []byte) error

	// Close gracefully closes the connection to the message queue.
	Close() error
}
