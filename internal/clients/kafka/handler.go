package kafka

import (
	"context"
	"log"

	"github.com/IBM/sarama"
)

type MessageHandler interface {
	HandleMessage(msg *sarama.ConsumerMessage)
}

type Consumer struct {
	Brokers []string
	GroupID string
	Topic   string
	Handler MessageHandler
}

func (c *Consumer) Start() error {
	log.Printf("Starting Kafka consumer for group %s on topic %s", c.GroupID, c.Topic)
	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0

	consumerGroup, err := sarama.NewConsumerGroup(c.Brokers, c.GroupID, config)
	if err != nil {
		return err
	}

	go func() {
		for {
			err := consumerGroup.Consume(context.Background(), []string{c.Topic}, &consumerGroupHandler{handler: c.Handler})
			if err != nil {
				log.Printf("[%s] Consume error: %v", c.GroupID, err)
			}
		}
	}()

	return nil
}

type consumerGroupHandler struct {
	handler MessageHandler
}

func (h *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }
func (h *consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		h.handler.HandleMessage(msg)
		sess.MarkMessage(msg, "")
	}
	return nil
}
