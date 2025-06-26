package consumer

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/JulieWasNotAvailable/microservices/unpublished/internal/beatmetadata"
	"github.com/JulieWasNotAvailable/microservices/unpublished/internal/entities"
	"github.com/JulieWasNotAvailable/microservices/unpublished/internal/unpbeat"
	"github.com/google/uuid"
)

func StartConsumerFileUpdate(topic string, service unpbeat.Service, mservice beatmetadata.MetadataService, appQuit chan bool) {
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
				appQuit <- true
			case msg := <-consumer.Messages():
				msgCount++
				fmt.Printf("Received message Count %d: | Topic(%s) | Message(%s) \n", msgCount, string(msg.Topic), string(msg.Value))

				message := KafkaMessageURLUpdate{}
				err := json.Unmarshal(msg.Value, &message)
				if err != nil {
					log.Panic(err)
				}
				key, err := uuid.Parse(string(msg.Key))
				if err != nil {
					log.Panic(err)
				}

				updateDataFiles := entities.AvailableFiles{}
				updateDataBeat := entities.UnpublishedBeat{
					ID: key,
				}
				switch message.FileType {
				case "mp3":
					updateDataFiles.MP3Url = message.URL
					_, err := mservice.UpdateAvailableFilesByBeatId(key, &updateDataFiles)
					if err != nil {
						log.Println("couldn't update files in ConsumerFileUpdate", err)
					}
				case "wav":
					updateDataFiles.WAVUrl = message.URL
					_, err := mservice.UpdateAvailableFilesByBeatId(key, &updateDataFiles)
					if err != nil {
						log.Println("couldn't update files in ConsumerFileUpdate", err)
					}
				case "zip":
					updateDataFiles.ZIPUrl = message.URL
					_, err := mservice.UpdateAvailableFilesByBeatId(key, &updateDataFiles)
					if err != nil {
						log.Println("couldn't update files in ConsumerFileUpdate", err)
					}
				case "cover":
					updateDataBeat.Picture = message.URL
					beatInitial, _ := service.GetUnpublishedBeatByID(key)
					_, err = service.UpdateUnpublishedBeat(&updateDataBeat, beatInitial.BeatmakerID)
					if err != nil {
						log.Println("couldn't update files in ConsumerFileUpdate", err)
					}
				default:
					log.Println("error in ConsumerFileUpdate", err)
				}

			case <-sigchan:
				fmt.Println("Interrupt is detected")

				//It sends an empty struct to doneCh, signaling that the goroutine should terminate.
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
