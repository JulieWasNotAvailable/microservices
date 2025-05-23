package consumer

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	// "github.com/JulieWasNotAvailable/microservices/user/api/presenters"
	// "github.com/JulieWasNotAvailable/microservices/user/pkg/user"
)

type KafkaMessageValue struct {
	ID       string    `json:"beat_id"`
	Features []float64 `json:"features"`
	Err      string    `json:"error"`
}

type KafkaMessageToMFCC struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
}

func StartConsumer(topic string, channel chan<- KafkaMessageValue) {
	brokerUrl := []string{"localhost:9092"}

	worker, err := connectConsumer(brokerUrl)
	if err != nil {
		panic(err)
	}

	consumer, err := worker.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}

	fmt.Println(("consumer started"))

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	msgCount := 0
	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				fmt.Println(err)
			case msg := <-consumer.Messages():
				msgCount++
				fmt.Printf("Received message Count %d: | Topic(%s) | Message(%s) \n", msgCount, string(msg.Topic), string(msg.Value))

				messageValue := KafkaMessageValue{}
				err := json.Unmarshal(msg.Value, &messageValue)
				if err != nil {
					messageValue.Err = "couldn't unmarshal"
					log.Println("cannot parse the message")
				}
				channel <- messageValue
			case <-sigchan:
				fmt.Println("Interrupt is detected")

				//It sends an empty struct to doneCh, signaling that the goroutine should terminate.
				doneCh <- struct{}{}
			}
		}
	}()

	<-doneCh

	if err := worker.Close(); err != nil {
		panic(err)
	}
	//we're waiting for a response from this channel
	fmt.Println("Processed", msgCount, "messages")
}

func connectConsumer(brokersUrl []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	conn, err := sarama.NewConsumer(brokersUrl, config)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
