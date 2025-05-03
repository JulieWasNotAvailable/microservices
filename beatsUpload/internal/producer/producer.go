package producer

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

type FileType string

const (
	mp3   FileType = "mp3"
	wav   FileType = "wav"
	zip   FileType = "zip"
	cover FileType = "cover"
	pfp   FileType = "pfp" //profile picture
)

type KafkaMessage struct {
	FileType FileType
	URL      string
}

func createProducer(brokersUrl []string) (sarama.AsyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewAsyncProducer(brokersUrl, config)
	if err != nil {
		return nil, err
	}

	return producer, nil
}

// pushes Update Message to queue
func pushMessageToQueue(topic string, key []byte, message []byte) error {

	brokerUrl := []string{"broker:29092"}

	producer, err := createProducer(brokerUrl)
	if err != nil {
		return err
	}

	defer producer.Close()

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(message),
	}

	//<- is used to push the message to channel producer.Input()
	producer.Input() <- msg

	select {
	//If Return.Successes is true, you MUST read from this channel or the Producer will deadlock.
	case success := <-producer.Successes():
		// An offset is a unique identifier assigned to each message in a Kafka partition. Used to track the position of a consumer within a partition.
		fmt.Println("Message produced:", success.Offset)
		fmt.Printf("Message is stored in the topic (%s)/partition(%d)/offset(%d)\n", success.Topic, success.Partition, success.Offset)
	case err := <-producer.Errors():
		fmt.Println("Failed to produce message:", err)
		return err
	}

	return nil
}

func CreateMessage(message KafkaMessage, key string, topic string) error {

	messageInBytes, err := json.Marshal(message)
	if err != nil {
		log.Fatal(messageInBytes)
	}

	keyInBytes, err := json.Marshal(key)
	if err != nil {
		log.Fatal(keyInBytes)
	}

	err = pushMessageToQueue(topic, keyInBytes, messageInBytes)
	if err != nil {
		log.Println("created message successfully")
	}

	return nil
}
