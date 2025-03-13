package updateurl

import (
	"log"
	"net/http"
	"time"

	"github.com/JulieWasNotAvailable/microservices/beatsUpload/dbconnection"
	"github.com/JulieWasNotAvailable/microservices/beatsUpload/producer"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
)

//если неправильно прописать реквест, будет паника (всё отрубится)
type Request struct {
	BucketName *string
	ObjectKey *string
	BeatId *int
}

func CheckFileAvailability(ctx *fiber.Ctx) error {
    client := dbconnection.S3Connect()

	req := Request{}

	err := ctx.BodyParser(&req)

	if err != nil {
		ctx.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	err = s3.NewObjectExistsWaiter(client.S3Client).Wait(
		ctx.Context(), &s3.HeadObjectInput{
			Bucket: aws.String(*req.BucketName),
			Key:    aws.String(*req.ObjectKey)}, time.Minute)

	if err != nil {
		log.Printf("Failed attempt to wait for object %s to exist.\n", *req.ObjectKey)
	} else {
		log.Println("sending message to kafka")
	}

	address := "https://storage.yandexcloud.net/" + *req.BucketName + "/" + *req.ObjectKey

	producer.CreateMessage(&address, req.BeatId)

	return nil
}
