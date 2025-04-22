package producer

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
)


func createProducer (brokersUrl []string) (sarama.AsyncProducer, error) {
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

//pushes Update Message to queue
func pushMessageToQueue (topic string, key []byte, message []byte) error {

	brokerUrl := []string{"localhost:9092"}

	producer, err := createProducer(brokerUrl)
	if err != nil {
		return err
	}

	defer producer.Close()
	
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key: sarama.StringEncoder(key),
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

func CreateMessage(url string, key string, topic string) error {

	urlInBytes, err := json.Marshal(url)
	if err != nil {
		log.Fatal(urlInBytes)
		}

	keyInBytes, err := json.Marshal(key)
	if err != nil {
		log.Fatal(keyInBytes)
		}

	err = pushMessageToQueue(topic, keyInBytes, urlInBytes)
	if err != nil {
		log.Println("created message successfully")
		}

	return nil
	}