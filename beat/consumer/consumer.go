package consumer

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	// "reflect"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/JulieWasNotAvailable/microservices/beat/pkg/beat"
	"github.com/JulieWasNotAvailable/microservices/beat/pkg/entities"
)

type KafkaMessageBeatForPublishing struct {
	Beat entities.UnpublishedBeat
	MFCC entities.MFCC
}

func StartConsumerPublisher(topic string, service beat.Service) {
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

				//if err, service, that sent, should be notified?
				messageValue := KafkaMessageBeatForPublishing{}
				err = json.Unmarshal(msg.Value, &messageValue)
				if err != nil {
					log.Println(err)
				}
				// log.Println(" ")
				// log.Println("in consumer")
				// log.Println(messageValue)
				// log.Println(reflect.TypeOf(messageValue.MFCC.Col1))

				// log.Println(messageValue.Beat, messageValue.MFCC)
				_, err := service.CreateBeat(messageValue.Beat, messageValue.MFCC)
				if err != nil {
					log.Println(err)
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
