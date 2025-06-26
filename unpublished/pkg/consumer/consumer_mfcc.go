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
)

func StartConsumerMFCC(topic string, service unpbeat.Service, appQuit chan bool) {
	brokerUrl := []string{"localhost:9092"}

	fmt.Printf("starting consumer with brokerurl %s on topic: %s \n", brokerUrl[0], topic)

	worker, err := connectConsumer(brokerUrl)
	if err != nil {
		panic(err)
	}

	consumer, err := worker.ConsumePartition(topic, 0, sarama.OffsetOldest)
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
				appQuit <- true
			case msg := <-consumer.Messages():
				msgCount++
				fmt.Printf("Received message Count %d: | Topic(%s) | Message(%s) \n", msgCount, string(msg.Topic), string(msg.Value))

				messageValue := KafkaMessageMFCC{}
				err := json.Unmarshal(msg.Value, &messageValue)
				if err != nil {
					log.Println("couldn't process the kafka message in consumer_Mfcc")
				}
				log.Println(msg.Offset)

				if messageValue.Error != "" {
					toUpdateErr := entities.UnpublishedBeat{
						Status: entities.StatusDraft,
						Err:    messageValue.Error,
					}
					beatInitial, _ := service.GetUnpublishedBeatByID(messageValue.ID)
					_, err = service.UpdateUnpublishedBeat(&toUpdateErr, beatInitial.BeatmakerID)
					if err != nil {
						log.Println("couldn't update the beat error after receving mfcc error message")
					}
				} else {
					beat, err := service.GetUnpublishedBeatByID(messageValue.ID)
					if err != nil {
						log.Println("couldn't get the beat while trying to sent the message to publish_beat_main")
					} else {
						toBeatMsg := producer.KafkaMessageBeatForPublishing{
							Beat: *beat,
							MFCC: messageValue.Features,
						}
						toBeatMsgBytes, _ := json.Marshal(toBeatMsg)
						producer.CreateMessage(toBeatMsgBytes, "publish_beat_main")
					}
				}

			case <-sigchan:
				fmt.Println("Interrupt is detected")

				doneCh <- struct{}{}
				appQuit <- true
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
