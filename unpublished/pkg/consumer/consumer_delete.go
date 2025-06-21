package consumer

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/JulieWasNotAvailable/microservices/unpublished/internal/unpbeat"
)

func StartConsumerDeleteApprove(topic string, service unpbeat.Service) {
	brokerUrl := []string{"broker:29092"}

	fmt.Printf("starting consumer with brokerurl %s on topic: %s \n", brokerUrl[0], topic)

	worker, err := connectConsumer(brokerUrl)
	if err != nil {
		panic(err)
	}

	consumer, err := worker.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}

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

				messageValue := KafkaMessageDeleteApprove{}
				err := json.Unmarshal(msg.Value, &messageValue)
				if err != nil {
					log.Println("cannot parse the delete_approve message ")
				} else if messageValue.Error != "" {
					log.Println("couldnt create the beat in publish service")
				} else {
					log.Println("deleting the beat with id", messageValue.BeatId)
					err := service.DeleteUnpublishedBeat(messageValue.BeatId)
					if err != nil {
						log.Println("cannot delete the beat in unpublished microservice.")
						log.Println("error", err)
					}
				}

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
