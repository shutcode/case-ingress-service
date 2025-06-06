package proto

import (
	"context"
	"encoding/json"
	"log"

	"github.com/cloudwego/hertz/pkg/app"
	// "github.com/cloudwego/hertz/pkg/protocol/consts" // Was unused
	"github.com/yourusername/yourrepository/biz/model/proto" // Generated proto models
	"github.com/yourusername/yourrepository/pkg/config"     // Our config package
	"github.com/yourusername/yourrepository/pkg/messagequeue" // Our message queue package
)

// CaseImportServiceService is the generated struct for the service.
type CaseImportServiceService struct {
	// Base service struct if any, or empty
}

// PublishCaseEvent implements the PublishCaseEvent method of the CaseImportService.
func (s *CaseImportServiceService) PublishCaseEvent(ctx context.Context, c *app.RequestContext, req *proto.ImportEventRequest) (*proto.ImportEventResponse, error) {
	log.Printf("Received PublishCaseEvent request: %+v", req)

	if req.Event == nil {
		log.Println("Error: Request event is nil")
		return &proto.ImportEventResponse{
			Success: false,
			Message: "Error: request event is nil",
		}, nil
	}

	// 1. Instantiate KafkaProducer using configuration
	// Ensure config.LoadConfig() has been called at application startup
	if len(config.Cfg.Kafka.Brokers) == 0 || config.Cfg.Kafka.Topic == "" {
		log.Println("Error: Kafka configuration (brokers or topic) is missing. Ensure LoadConfig() was called.")
		// Attempt to load config here as a fallback for safety, though it should be done in main.
		// This is not ideal for production but helps during development if main isn't updated yet.
		if err := config.LoadConfig(); err != nil {
			log.Printf("Fallback LoadConfig failed: %v", err)
		}
		if len(config.Cfg.Kafka.Brokers) == 0 || config.Cfg.Kafka.Topic == "" {
			return &proto.ImportEventResponse{
				Success: false,
				Message: "Error: Kafka configuration is missing",
			}, nil
		}
	}

	mqProducer, err := messagequeue.NewKafkaProducer(config.Cfg.Kafka.Brokers)
	if err != nil {
		log.Printf("Error creating Kafka producer: %v", err)
		return &proto.ImportEventResponse{
			Success: false,
			Message: "Error initializing message queue: " + err.Error(),
		}, nil
	}
	defer mqProducer.Close()

	// 2. Serialize CaseEvent to JSON
	eventData, err := json.Marshal(req.Event)
	if err != nil {
		log.Printf("Error marshalling event data to JSON: %v", err)
		return &proto.ImportEventResponse{
			Success: false,
			Message: "Error serializing event data: " + err.Error(),
		}, nil
	}

	// 3. Publish to Kafka using configured topic
	log.Printf("Publishing event to Kafka topic '%s': %s", config.Cfg.Kafka.Topic, string(eventData))
	err = mqProducer.Publish(config.Cfg.Kafka.Topic, eventData)
	if err != nil {
		log.Printf("Error publishing message to Kafka: %v", err)
		return &proto.ImportEventResponse{
			Success: false,
			Message: "Error publishing event to message queue: " + err.Error(),
		}, nil
	}

	log.Printf("Event published successfully to Kafka topic '%s'", config.Cfg.Kafka.Topic)
	return &proto.ImportEventResponse{
		Success: true,
		Message: "Event published successfully",
		EventId: req.Event.CaseId,
	}, nil
}
