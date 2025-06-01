package beat

import (
	"github.com/JulieWasNotAvailable/microservices/beatsUpload/pkg/dbconnection"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Presigner encapsulates the Amazon Simple Storage Service (Amazon S3) presign actions
// used in the examples.
// It contains PresignClient, a client that is used to presign requests to Amazon S3.
// Presigned requests contain temporary credentials and can be made from any HTTP client.
type Presigner struct {
	PresignClient *s3.PresignClient
}

func S3ConnectPresign (storage dbconnection.Storage) Presigner {
    newPresignClient := s3.NewPresignClient(storage.S3Client)
	presigner := Presigner{
		PresignClient: newPresignClient,
	}
    return presigner
}