package s3files

import (
	"log"
	"net/http"
	"time"

	"github.com/JulieWasNotAvailable/microservices/beatsUpload/dbconnection"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
)

// Presigner encapsulates the Amazon Simple Storage Service (Amazon S3) presign actions
// used in the examples.
// It contains PresignClient, a client that is used to presign requests to Amazon S3.
// Presigned requests contain temporary credentials and can be made from any HTTP client.
type Presigner struct {
	PresignClient *s3.PresignClient
}

//attributes need to start from capital so that they would be parsed
type request struct {
	BucketName *string
	ObjectKey *string
}

var lifetimeSecs int64 = 300

func S3ConnectPresign (storage dbconnection.Storage) Presigner {
    newPresignClient := s3.NewPresignClient(storage.S3Client)
	presigner := Presigner{
		PresignClient: newPresignClient,
	}
    return presigner
}

// GetObject makes a presigned request that can be used to get an object from a bucket.
// The presigned request is valid for the specified number of seconds.
func (presigner Presigner) GetObject (ctx *fiber.Ctx) error {
	req := request{}
	
	err := ctx.BodyParser(&req)

	log.Println("tryin to get an object")

	if err != nil {
		ctx.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message":"could not parse the json"})
			return err
		}
	
	//write a func to check if this object exists
	presignedRequest, err := presigner.PresignClient.PresignGetObject(ctx.Context(), &s3.GetObjectInput{
		Bucket: aws.String(*req.BucketName),
		Key:    aws.String(*req.ObjectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(lifetimeSecs * int64(time.Second))
	})

	log.Println("here's the presigned request: ", presignedRequest)

	if err != nil {
		ctx.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"message":"Couldn't generate a presigned get request. Here's why: ",
				"error": err,
				"accepted_data": req})
			return err
		}

	ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message":"generated presigned get request successfully",
		"data": presignedRequest})

	return nil
}

// PutObject makes a presigned request that can be used to put an object in a bucket.
// The presigned request is valid for the specified number of seconds.
func (presigner Presigner) PutObject(ctx *fiber.Ctx) error {
	req := request{}
	
	err := ctx.BodyParser(&req)

	if err != nil {
		ctx.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message":"could not parse the json"})
			return err
		}

	presignedRequest, err := presigner.PresignClient.PresignPutObject(ctx.Context(), &s3.PutObjectInput{
		Bucket: aws.String(*req.BucketName),
		Key:    aws.String(*req.ObjectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(lifetimeSecs * int64(time.Second))
	})
	
	if err != nil {
		ctx.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"message":"Couldn't generate a presigned put request. Here's why: ",
				"error": err,
				"accepted_data": req})
			return err
		}

	ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message":"generated presigned put request successfully",
		"data": presignedRequest})

	return nil
}

// DeleteObject makes a presigned request that can be used to delete an object from a bucket.
func (presigner Presigner) DeleteObject(ctx *fiber.Ctx) error {
	req := request{}
	
	err := ctx.BodyParser(&req)

	if err != nil {
		ctx.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message":"could not parse the json"})
			return err
		}

	presignedRequest, err := presigner.PresignClient.PresignDeleteObject(ctx.Context(), &s3.DeleteObjectInput{
		Bucket: aws.String(*req.BucketName),
		Key:    aws.String(*req.ObjectKey),
	})
	
	if err != nil {
		ctx.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"message":"Couldn't generate a presigned delete request. Here's why: ",
				"error": err,
				"accepted_data": req})
			return err
		}

	ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message":"generated presigned delete request successfully",
		"data": presignedRequest})

	return nil
}

func (p Presigner) SetupRoutes (app *fiber.App) {
	api := app.Group("/api/presigned")
	api.Get("/getPresignedGetRequest", p.GetObject)
	api.Get("/getPresignedPostRequest", p.PutObject)
	api.Get("/getPresignedDeleteRequest", p.DeleteObject)
}