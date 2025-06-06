package messagequeue

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

// KafkaProducer implements the MessageQueue interface using Kafka.
type KafkaProducer struct {
	writer *kafka.Writer
}

// NewKafkaProducer creates a new Kafka producer.
// brokers is a list of Kafka broker addresses (e.g., "localhost:9092").
func NewKafkaProducer(brokers []string) (*KafkaProducer, error) {
	// Configure the Kafka writer.
	// The topic can be set here if it's fixed, or passed to Publish if dynamic.
	// For flexibility, we'll allow topic to be specified in Publish.
	writer := &kafka.Writer{
		Addr: kafka.TCP(brokers...),
		// Balancer: &kafka.LeastBytes{}, // Default is RoundRobin
		RequiredAcks: kafka.RequireAll, // Wait for all in-sync replicas to ack
		Async:        false,            // Make writes synchronous
		WriteTimeout: 10 * time.Second,
	}

	// TODO: Consider adding error checking for broker connectivity here if possible,
	// though kafka-go writer often establishes connection on first write.

	return &KafkaProducer{writer: writer}, nil
}

// Publish sends a message to the specified Kafka topic.
func (kp *KafkaProducer) Publish(topic string, message []byte) error {
	kafkaMsg := kafka.Message{
		Topic: topic,
		Value: message,
		// Key: []byte("some-key"), // Optional: for partitioning
	}

	// Using context.Background() for now, consider making context configurable
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return kp.writer.WriteMessages(ctx, kafkaMsg)
}

// Close closes the Kafka writer.
func (kp *KafkaProducer) Close() error {
	if kp.writer != nil {
		return kp.writer.Close()
	}
	return nil
}
