package handler

import (
	"errors"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/JulieWasNotAvailable/microservices/beatsUpload/internal/producer"
	"github.com/JulieWasNotAvailable/microservices/beatsUpload/pkg/dbconnection"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
)

type Request struct {
	BucketName string
	ObjectKey string
	EntityType string
	Entity string
}

func checkContentType(contentType string, bucketName string, objectKey string, ctx *fiber.Ctx) error {
	client := dbconnection.S3Connect()

	re := regexp.MustCompile(`audio/mpeg|audio/wav|application/zip|image/jpeg|image/jpg|image/png`)
	if !re.MatchString(contentType) {
		return errors.New("content type is not available")
	}

	switch contentType {
	case "audio/mpeg":
		if bucketName != "mp3beats"{
			_, err := client.S3Client.DeleteObject(ctx.Context(), &s3.DeleteObjectInput{
				Bucket : aws.String(bucketName),
				Key : aws.String(objectKey),
			})
			if err != nil{
				return ctx.Status(http.StatusInternalServerError).JSON(
					&fiber.Map{
						"message": "couldn't delete the object with the wrong content-type. please, try again",
						"error":err.Error()})
			}
			return errors.New("wrong bucket. audio/mpeg should be in mp3beats. i'm deleting that")
		}
		
	case "audio/wav":
		if bucketName != "wavbeats"{
			_, err := client.S3Client.DeleteObject(ctx.Context(), &s3.DeleteObjectInput{
				Bucket : aws.String(bucketName),
				Key : aws.String(objectKey),
			})
			if err != nil{
				return ctx.Status(http.StatusInternalServerError).JSON(
					&fiber.Map{
						"message": "couldn't delete the object with the wrong content-type. please, try again",
						"error":err.Error()})
			}
			return errors.New("wrong bucket. audio/wav should be in wavbeats. i'm deleting that")
		}
	case "application/zip":
		if bucketName != "zipbeats"{
			_, err := client.S3Client.DeleteObject(ctx.Context(), &s3.DeleteObjectInput{
				Bucket : aws.String(bucketName),
				Key : aws.String(objectKey),
			})
			if err != nil{
				return ctx.Status(http.StatusInternalServerError).JSON(
					&fiber.Map{
						"message": "couldn't delete the object with the wrong content-type. please, try again",
						"error":err.Error()})
			}
			return errors.New("wrong bucket. application/zip should be in zipbeats. i'm deleting that")
		}
	case "image/jpeg", "image/jpg", "image/png":
		if bucketName != "imagesall"{
			_, err := client.S3Client.DeleteObject(ctx.Context(), &s3.DeleteObjectInput{
				Bucket : aws.String(bucketName),
				Key : aws.String(objectKey),
			})
			if err != nil{
				return ctx.Status(http.StatusInternalServerError).JSON(
					&fiber.Map{
						"message": "couldn't delete the object with the wrong content-type. please, try again",
						"error":err.Error()})
			}
			return errors.New("wrong bucket. images should be in imagesall. i'm deleting that")
		}
	default:
		return nil
	}

	return nil
}

// @Summary Validates the file, pushes to User or Beat Service.
// @Description Verify if a file exists in S3 and publish to Kafka
// @Tags Storage
// @Accept json
// @Produce json
// @Param request body Request true "File details"
// @Success 200 {object} map[string]interface{} "Successfully processed"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 422 {object} map[string]interface{} "Unprocessable entity"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /checkFileUpdateUrl [post]
func CheckFileAvailability(ctx *fiber.Ctx) error {
    client := dbconnection.S3Connect()

	req := Request{}

	err := ctx.BodyParser(&req)
	if err != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
	}

	err = s3.NewObjectExistsWaiter(client.S3Client).Wait(
		ctx.Context(), &s3.HeadObjectInput{
			Bucket: aws.String(req.BucketName),
			Key:    aws.String(req.ObjectKey)}, time.Second*10)
	if err != nil {
		log.Printf("Failed attempt to wait for object %s to exist.\n", req.ObjectKey)
		return ctx.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{
				"message": "couldn't wait for the object from s3. please, try again",
				"error":err.Error()})
	}

	objectMetaData, err := client.S3Client.HeadObject(ctx.Context(), &s3.HeadObjectInput{
		Bucket: aws.String(req.BucketName),
		Key:    aws.String(req.ObjectKey)})
	
	if err != nil {
		log.Printf("Failed attempt to wait for object %s to exist.\n", req.ObjectKey)
		return ctx.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{
				"message": "couldn't get the objects metadata. please, try again",
				"error":err})
			}

	var topic string
	if req.EntityType == "User" {
		topic = "profilepic_url_updates"
	} else if req.EntityType == "Beat"{
			topic = "beat_url_updates"
		} else {
		return ctx.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"message": "wrong request type, entity can only be User or Beat",
			})
	}
	
	allowedContentypes := [5]string{"audio/mpeg", "audio/wav", "application/zip", "image/jpeg", "image/png"}
	err = checkContentType(*objectMetaData.ContentType, req.BucketName, req.ObjectKey, ctx)		
	if err != nil {return ctx.Status(http.StatusInternalServerError).JSON(
		&fiber.Map{
			"message": "object had an incorrect type, so it was deleted",
			"err": err.Error(),
			"allowed content types" : allowedContentypes})	}
	

	address := req.BucketName + "/" + req.ObjectKey	
	err = producer.CreateMessage(address, req.Entity, topic)
	if err != nil{
		return ctx.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{
				"message": "couldn't send the message to kafka",
				"error":err.Error()})
	}

	return ctx.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message": "successfully pushed to kafka",
		})
}