package consumer

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"os"
// 	"os/signal"
// 	"syscall"

// 	"github.com/IBM/sarama"
// 	"github.com/JulieWasNotAvailable/microservices/beat/internal/producer"
// 	"github.com/JulieWasNotAvailable/microservices/beat/pkg/beat"
// )

// func StartConsumerPublisherErr (topic string, service beat.Service) {
// 	worker, err := connectConsumer([]string{os.Getenv("KAFKA_BROKER_URL")})

// 	if err != nil {
// 		panic(err)
// 	}

// 	consumer, err := worker.ConsumePartition(topic, 0, sarama.OffsetOldest)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println(("consumer started"))
// 	sigchan := make(chan os.Signal, 1)

// 	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

// 	msgCount := 0

// 	doneCh := make(chan struct{})
// 	go func() {
// 		for {
// 			select {
// 			case err := <-consumer.Errors():
// 				fmt.Println(err)
// 			case msg := <-consumer.Messages():
// 				msgCount++
// 				fmt.Printf("Received message Count %d: | Topic(%s) | Message(%s) \n", msgCount, string(msg.Topic), string(msg.Value))

// 				message := producer.KafkaMessageValue{}
// 				err = json.Unmarshal(msg.Value, &message)
// 				if err != nil {
// 					log.Println(err)
// 				}

// 				err = service.DeleteBeatById(message.ID)
// 				if err != nil {
// 					log.Println(err)
// 				}

// 			case <-sigchan:
// 				fmt.Println("Interrupt is detected")
// 				//It sends an empty struct to doneCh, signaling that the goroutine should terminate.
// 				doneCh <- struct{}{}
// 			}
// 		}
// 	}()

// 	//we're waiting for a response from this channel
// 	<-doneCh
// 	fmt.Println("Processed", msgCount, "messages")

// 	if err := worker.Close(); err != nil {
// 		panic(err)
// 	}
// }
