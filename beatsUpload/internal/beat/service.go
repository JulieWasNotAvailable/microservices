package beat

import (
	"context"
	"log"

	"github.com/JulieWasNotAvailable/microservices/beatsUpload/pkg/dbconnection"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type AvailableFiles struct {
	ID                uuid.UUID `json:"id"`
	MP3Url            string    `json:"mp3url"`
	WAVUrl            string    `json:"wavurl"`
	ZIPUrl            string    `json:"zipurl"`
	UnpublishedBeatID uuid.UUID `json:"unpublishedBeatId"`
}

type KafkaMessage struct {
	AvailableFiles AvailableFiles `json:"availableFiles"`
	Cover          string         `json:"cover"`
}

type Service interface{
	DeleteObject(filename string, bucket string)(string, error)
	DeleteAttachedToBeatFiles(message KafkaMessage) (error)
}

// s3 storage
type service struct {
	storage dbconnection.Storage
}

func NewService(s dbconnection.Storage) Service {
	return &service{storage: s}
}

func (s *service) DeleteObject (filename string, bucket string)(string, error) {
	storage := dbconnection.S3Connect()
	object, err := storage.S3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket : aws.String(bucket),
	})
	log.Println(object.ResultMetadata)
	if err != nil{
		return "", err
	}
	return filename, nil
}

func (s *service) DeleteAttachedToBeatFiles (message KafkaMessage) (error) {
	if message.AvailableFiles.MP3Url != "" {
		_, err := s.DeleteObject(message.AvailableFiles.MP3Url, "mp3beats")
		if err != nil {
			return (err)
		}
	}
	if message.AvailableFiles.WAVUrl != ""{
		_, err := s.DeleteObject(message.AvailableFiles.WAVUrl, "mp3beats")
		if err != nil {
			return (err)
		}
	}
	if message.AvailableFiles.ZIPUrl != ""{
		_, err := s.DeleteObject(message.AvailableFiles.ZIPUrl, "mp3beats")
		if err != nil {
			return (err)
		}
	}
	if message.Cover != ""{
		_, err := s.DeleteObject(message.Cover, "mp3beats")
		if err != nil {
			return (err)
		}
	}
	return nil
}


				