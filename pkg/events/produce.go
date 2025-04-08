package events

import "github.com/IBM/sarama"

func (k *KafkaClient) Produce(topic, key string, message []byte) (int32, int64, error) {

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(message),
	}

	return k.producer.SendMessage(msg)
}
