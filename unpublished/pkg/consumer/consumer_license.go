package consumer

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/JulieWasNotAvailable/microservices/unpublished/internal/entities"
	"github.com/JulieWasNotAvailable/microservices/unpublished/internal/unpbeat"
	"github.com/JulieWasNotAvailable/microservices/unpublished/pkg/producer"
	// "github.com/JulieWasNotAvailable/microservices/user/api/presenters"
	// "github.com/JulieWasNotAvailable/microservices/user/pkg/user"
)

func StartConsumerLicense(topic string, service unpbeat.Service) {
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

				messageValue := KafkaMessageLicenseCreated{}
				err := json.Unmarshal(msg.Value, &messageValue)
				if err != nil {
					messageValue.Error = "couldn't unmarshal"
					log.Println("cannot parse the message")
				}
				if messageValue.Error != "" {
					toUpdateErr := entities.UnpublishedBeat{
						Status: entities.StatusDraft,
						Err: messageValue.Error,
					}
					beatInitial, _ := service.GetUnpublishedBeatByID(messageValue.BeatId)
					_, err = service.UpdateUnpublishedBeat(&toUpdateErr, beatInitial.BeatmakerID)
					if err != nil {
						log.Println("couldn't update the beat error after receving license error message")
					}
				} else {
					beat, err := service.GetUnpublishedBeatByID(messageValue.BeatId)
					if err != nil {
						log.Println("couldn't update the beat error after receving license error message")
					}
					toMfccMsg := producer.KafkaMessageToMFCC{
						ID : messageValue.BeatId.String(),
						Filename: beat.AvailableFiles.MP3Url,
					}
					toMfccMsgBytes, _ := json.Marshal(toMfccMsg)
					producer.CreateMessage(toMfccMsgBytes, "track_for_mfcc")
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
