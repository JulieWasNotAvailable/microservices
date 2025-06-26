package consumer

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/JulieWasNotAvailable/microservices/beat/internal/beat"
	"github.com/JulieWasNotAvailable/microservices/beat/internal/entities"
)

type KafkaMessageBeatForPublishing struct {
	Beat entities.UnpublishedBeat
	MFCC []float64
}

func StartConsumerPublisher(topic string, service beat.Service, appQuit chan bool) {
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
				fmt.Printf("Received message Count %d: | Topic(%s) | Message(written below) \n", msgCount, string(msg.Topic))

				message := KafkaMessageBeatForPublishing{}
				err = json.Unmarshal(msg.Value, &message)
				if err != nil {
					log.Println(err)
				} else {
					log.Println("Message (mfcc excluded): ", message.Beat)
					_, err := service.CreateBeat(message.Beat, message.MFCC)
					if err != nil {
						log.Println(err)
					}
				}

			case <-sigchan:
				fmt.Println("Interrupt is detected")
				doneCh <- struct{}{}
				appQuit <- true
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
