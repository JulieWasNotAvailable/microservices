package dbconnection

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Storage struct {
	S3Client *s3.Client
}


type Bucket struct{
    Name *string
}

func S3Connect() Storage {
    //по примеру из яндекса
    // Подгружаем конфигурацию из user/1/.aws/config*
    cfg, err := config.LoadDefaultConfig(context.TODO())
    if err != nil {
        log.Fatal(err)
    }

    // Создаем клиента для доступа к хранилищу S3
    client := s3.NewFromConfig(cfg)

    storage := Storage{
        S3Client: client,
    }
    return storage
}