package consumer

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"github.com/IBM/sarama"
	"github.com/JulieWasNotAvailable/goBeatsBackend/beat"
	"github.com/JulieWasNotAvailable/goBeatsBackend/dbconnection"
	"github.com/JulieWasNotAvailable/goBeatsBackend/model"
)

func StartConsumer(topic string){
	pgconfig := dbconnection.NewConfig()
	pgdb, err := dbconnection.NewConnection(pgconfig)
	if err != nil{
		log.Fatal("could not reach the DB")
	}
	err = model.MigrateBeats(pgdb)
	if err != nil {
		log.Fatal("could not migrate")
	}
	pgrepo := beat.Repository{
		DB: pgdb,
	}
	
	worker, err := connectConsumer([]string{"localhost:9092"})
	if err != nil {
		panic (err)
	}

	consumer, err := worker.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}

	fmt.Println(("consumer started"))
	//chan - chanel, used for goroutines to exchange messages
	/*
	make(chan os.Signal, 1)  creates a new channel of type os.Signal
	ctrl+c interruptions or terminal signals
	this channel can only hold 1 signal
	if the second signal is send, the goroutine that sends messages to this channel will be blocked
	*/

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
				
				key := uint(msg.Key[0])
				log.Println(key)
				valuebytes := msg.Value
				value := string(valuebytes)
				log.Println(value)
				
				req := beat.UpdateBeatURLRequest{
					BeatId: &key,
					Url: &value,
				}
				log.Println("sending the req to an UpdateBeatUrl func")
				pgrepo.UpdateBeatURL(&req)

			case <- sigchan:
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