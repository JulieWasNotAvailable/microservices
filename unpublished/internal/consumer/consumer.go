package consumer

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"log"
	"encoding/json"

	"github.com/IBM/sarama"
	// "github.com/JulieWasNotAvailable/microservices/user/api/presenters"
	// "github.com/JulieWasNotAvailable/microservices/user/pkg/user"
)

type KafkaMessage struct{
	Key string	
	Value string
	Err string
}

type MessageData struct{	
	Value json.RawMessage
	Err string `json:"error"`
}

func StartConsumer(topic string, channel chan<- KafkaMessage){
	worker, err := connectConsumer([]string{os.Getenv("KAFKA_BROKER_URL")})
	if err != nil {
		panic (err)
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

				messageModel := MessageData{}
				err := json.Unmarshal(msg.Value, &messageModel)
				if err != nil{
					kafkaMessage := KafkaMessage{
						Key: string(msg.Key),
						Value: "", //mfcc, error
						Err : err.Error(),
					}
					channel <- kafkaMessage
				} else {
					kafkaMessage := KafkaMessage{
						Key: string(msg.Key),
						Value: string(messageModel.Value), //mfcc, error
						Err : messageModel.Err,
					}
					
					log.Println(kafkaMessage.Value)
					channel <- kafkaMessage
				}

			case <- sigchan:
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