package consumer

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/JulieWasNotAvailable/microservices/cart/internal/entities"
	"github.com/JulieWasNotAvailable/microservices/cart/internal/license"
	"github.com/JulieWasNotAvailable/microservices/cart/pkg/producer"
	"github.com/google/uuid"
)

func StartConsumer(topic string, service license.Service) {
	brokerUrl := []string{"localhost:9092"}

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

				messageValue := entities.CreateLicense{}
				err := json.Unmarshal(msg.Value, &messageValue)
				if err != nil {
					errMessage := producer.KafkaMessageCreateLicense{
						BeatId: uuid.Nil,
						Error:  "could not unmarshal the message",
					}
					errMsgBytes, _ := json.Marshal(errMessage)
					producer.CreateMessage(errMsgBytes, "license_was_created")
				} else {
					_, err = service.InsertNewLicenseList(messageValue.BeatId, messageValue.UserId, messageValue.LicenseList)
					if err != nil {
						errMessage := producer.KafkaMessageCreateLicense{
							BeatId: messageValue.BeatId,
							Error:  err.Error(),
						}
						errMsgBytes, _ := json.Marshal(errMessage)
						producer.CreateMessage(errMsgBytes, "license_was_created")
					} else {
						successMessage := producer.KafkaMessageCreateLicense{
							BeatId: messageValue.BeatId,
							Error:  "",
						}
						errMsgBytes, _ := json.Marshal(successMessage)
						producer.CreateMessage(errMsgBytes, "license_was_created")
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

func connectConsumer(brokersUrl []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	conn, err := sarama.NewConsumer(brokersUrl, config)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
