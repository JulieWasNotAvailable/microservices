package handler

import (
	"errors"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/JulieWasNotAvailable/microservices/beatsUpload/pkg/producer"
	"github.com/JulieWasNotAvailable/microservices/beatsUpload/pkg/dbconnection"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
)

type UpdateRequest struct {
	Id string `json:"id" example:"019623bd-3d0b-7dc2-8a1f-f782adeb42b4"`
	ObjectKey string `json:"objectKey" example:"019623bd-3d0b-7dc2-8a1f-f782adeb42b4"`
}

func checkContentType(contentType string, bucketName string, objectKey string, ctx *fiber.Ctx) error {
	client := dbconnection.S3Connect()

	re := regexp.MustCompile(`audio/mpeg|audio/wav|application/zip|image/jpeg|image/jpg|image/png`)
	if !re.MatchString(contentType) {
		_, err := client.S3Client.DeleteObject(ctx.Context(), &s3.DeleteObjectInput{
			Bucket : aws.String(bucketName),
			Key : aws.String(objectKey),
		})
		if err != nil{
			return err
		}
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
				return err
			}
			return errors.New("wrong bucket. audio/mpeg should be in mp3beats. i'm deleting that")
		}
		
	case "audio/wav":
		log.Println(bucketName)
		if bucketName != "wavbeats"{
			_, err := client.S3Client.DeleteObject(ctx.Context(), &s3.DeleteObjectInput{
				Bucket : aws.String(bucketName),
				Key : aws.String(objectKey),
			})
			if err != nil{
				return err
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
				return err
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
				return err
			}
			return errors.New("wrong bucket. images should be in imagesall. i'm deleting that")
		}
	default:
		return nil
	}

	return nil
}

// @Summary Validates the file, pushes to User or Beat Service.
// @Description Verify if a file exists, file type in S3 and publish to Kafka
// @Tags Update
// @Accept json
// @Produce json
// @Param entity path string true "User or Beat"
// @Param filetype path string false "fileType (mp3, wav, zip, cover or pfp)"
// @Param UpdateRequest body UpdateRequest true "UpdateRequest"
// @Success 200 {object} map[string]interface{} "Successfully processed"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 422 {object} map[string]interface{} "Unprocessable entity"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /updateURL/{entity}/{filetype} [post]
func UpdateFile(ctx *fiber.Ctx) error {
	client := dbconnection.S3Connect()
	
	entity := ctx.Params("entity") //user or beat
	filetype := ctx.Params("filetype") //mp3 wav or zip

	var bucketName string
	switch {
	case entity == "user":
		bucketName = "imagesall"
	case entity == "beat":
		if filetype == "mp3"{
			bucketName = "mp3beats"
		}
		if filetype == "wav"{
			bucketName = "wavbeats"
		}
		if filetype == "zip"{
			bucketName = "zipbeats"
		}
		if filetype == "cover"{
			bucketName = "imagesall"
		}
	default:
		return ctx.Status(http.StatusNotFound).JSON(
			&fiber.Map{"message": "request should be in the followng form: /api/user or /api/beat/mp3"})
	}

	req := UpdateRequest{}
	
	err := ctx.BodyParser(&req)
	if err != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "unprocessable entity"})
	}

	err = s3.NewObjectExistsWaiter(client.S3Client).Wait(
		ctx.Context(), &s3.HeadObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(req.ObjectKey)}, time.Second*10)
	if err != nil {
		log.Printf("Failed attempt to wait for object %s to exist.\n", req.ObjectKey)
		return ctx.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{
				"message": "cannot find this file in this bucket/not available yet",
				"error":err.Error()})
	}

	objectMetaData, err := client.S3Client.HeadObject(ctx.Context(), &s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(req.ObjectKey)})
	
	if err != nil {
		// log.Printf("Failed attempt to wait for object %s to exist.\n", req.ObjectKey)
		return ctx.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{
				"message": "the object does not exist. please, try again",
				"error":err})
			}

	var topic string
	if entity == "user" {
		topic = "profilepic_url_updates"
	} else if entity == "beat"{
			topic = "beat_files_updates"
		} else {
		return ctx.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"message": "wrong request type, entity can only be User or Beat",
			})
	}
	
	allowedContentypes := [5]string{"audio/mpeg", "audio/wav", "application/zip", "image/jpeg", "image/png"}
	err = checkContentType(*objectMetaData.ContentType, bucketName, req.ObjectKey, ctx)		
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(
		&fiber.Map{
			"message": "object had an incorrect type, so it was deleted",
			"err": err.Error(),
			"allowed content types" : allowedContentypes})	}
	
	url := "storage.yandexcloud.net/"
	address := url + bucketName + "/" + req.ObjectKey	
	message := producer.KafkaMessage{
		FileType: producer.FileType(filetype),
		URL : address,
	}
	err = producer.CreateMessage(message, req.Id, topic)
	if err != nil{
		return ctx.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{
				"message": "couldn't send the message to kafka",
				"error":err.Error()})
	}

	return ctx.Status(http.StatusOK).JSON(
		&fiber.Map{
			"status" : true,
			"message": "file checked, update request successfully pushed to kafka",
		})
}