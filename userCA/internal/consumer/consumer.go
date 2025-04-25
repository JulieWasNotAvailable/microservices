package consumer

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/JulieWasNotAvailable/microservices/user/api/presenters"
	"github.com/JulieWasNotAvailable/microservices/user/pkg/user"
	"github.com/google/uuid"
)

type KafkaMessageURLUpdate struct {
	FileType string
	URL      string
}

func StartConsumer(topic string, uservice user.Service) {
	worker, err := connectConsumer([]string{os.Getenv("KAFKA_BROKER_URL")})

	if err != nil {
		panic(err)
	}

	consumer, err := worker.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}

	fmt.Println(("consumer started"))
	sigchan := make(chan os.Signal, 1)

	// Ctrl + C is SIGINT
	// SIGTERM is gracefull shutdown
	// A graceful shutdown refers to the process where a system, service, or application is brought down in a managed and orderly way.
	// Notify listens to these values and sends them to signchan channel
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	msgCount := 0

	// doneChannel is made to notify that "hey i'm done consuming from streams"
	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				fmt.Println(err)
			case msg := <-consumer.Messages():
				msgCount++
				fmt.Printf("Received message Count %d: | Topic(%s) | Message(%s) \n", msgCount, string(msg.Topic), string(msg.Value))

				//if err, service, that sent, should be notified?
				key, err := uuid.Parse(string(msg.Key))
				if err != nil {
					log.Println(err)
				}
				message := KafkaMessageURLUpdate{}
				err = json.Unmarshal(msg.Value, &message)
				if err != nil {
					log.Println(err)
				}
				userUpdate := presenters.User{}
				if message.FileType == "pfp" {
					userUpdate.ID = key
					userUpdate.ProfilePictureUrl = message.URL
					_, err = uservice.UpdateUser(&userUpdate)
					if err != nil {
						log.Println(err)
					}
				} else {
					log.Println("wrong request")
				}

			case <-sigchan:
				fmt.Println("Interrupt is detected")
				//It sends an empty struct to doneCh, signaling that the goroutine should terminate.
				doneCh <- struct{}{}
			}
		}
	}()

	//we're waiting for a response from this channel
	<-doneCh
	fmt.Println("Processed", msgCount, "messages")

	if err := worker.Close(); err != nil {
		panic(err)
	}
}

func connectConsumer(brokersUrl []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// Create new consumer
	conn, err := sarama.NewConsumer(brokersUrl, config)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
