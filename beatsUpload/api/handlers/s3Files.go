package handler

import (
	"errors"
	"strings"

	"net/http"
	"time"

	"github.com/JulieWasNotAvailable/microservices/beatsUpload/pkg/dbconnection"
	"github.com/JulieWasNotAvailable/microservices/beatsUpload/internal/beat"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
	"github.com/gofiber/fiber/v2"
)

//attributes need to start from capital so that they would be parsed
type request struct {
	UuidFileName string `json:"uuidFileName" example:"019623bd-3d0b-7dc2-8a1f-f782adeb42b4"`
	File string `json:"file" example:"019623bd-3d0b-7dc2-8a1f-f782adeb42b4"`
}

var lifetimeSecs int64 = 300

func typeOf(objectKey string) (string, error) {
	contentSplit := strings.Split(objectKey, ".")
	if len(contentSplit) == 1{
		return "", errors.New("file type should be in the following format: file.jpg")
	}
	contentType := contentSplit[1]
	picTypes := [3]string{"jpeg", "jpg", "png"}
	for _, ptype := range(picTypes){
		if contentType == ptype{
			contentType = "image"
		}
	}
	return contentType, nil
}

// @Summary List all S3 buckets
// @Description Get a list of all available S3 buckets
// @Tags Storage
// @Produce json
// @Success 200 {object} map[string]interface{} "Returns list of buckets"
// @Failure 401 {string} string "Access denied"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /buckets [get]
func GetBuckets (ctx *fiber.Ctx) error {
	storage := dbconnection.S3Connect()

	var err error

	buckets, err := storage.S3Client.ListBuckets(ctx.Context(), &s3.ListBucketsInput{})
	
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "AccessDenied" {
			err = apiErr
			return ctx.Status(http.StatusUnauthorized).JSON("AccessDenied")
		} else {
			return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": "Couldn't list buckets for your account",
				"error": err,
			})
		}
	}
	return ctx.JSON(fiber.Map{
		"message" : "success",
		"buckets" : buckets,
		})
}

// @Summary List objects in a bucket  
// @Description Get a list of objects from the specified S3 bucket  
// @Tags Storage  
// @Produce json  
// @Param bucket path string true "Bucket name"  
// @Success 200 {object} map[string]interface{} "Returns list of objects"  
// @Failure 400 {object} map[string]interface{} "Bad request"  
// @Router /getObjectsFromBucket/{bucket} [get]  
func GetObjectsFromBucket(ctx *fiber.Ctx) error {
	bucket := ctx.Params("bucket")
	storage := dbconnection.S3Connect()
	objects, err := storage.S3Client.ListObjects(ctx.Context(), &s3.ListObjectsInput{
		Bucket : aws.String(bucket),
	})
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"message":"error occured",
				"error": err.Error()})
	}
	
	return ctx.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message":"success",
			"objects": objects.Contents})
}

// @Summary List objects in a bucket  
// @Description Get a list of objects from the specified S3 bucket  
// @Tags Storage  
// @Produce json  
// @Param bucket path string true "Bucket name"  
// @Success 200 {object} map[string]interface{} "Returns list of objects"  
// @Failure 400 {object} map[string]interface{} "Bad request"  
// @Router /ObjectsFromBucket/{bucket} [get] 
func GetHeadObject(ctx *fiber.Ctx) error {
	bucket := ctx.Params("bucket")
	key := ctx.Params("key")
	storage := dbconnection.S3Connect()
	head, err := storage.S3Client.HeadObject(ctx.Context(), &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key)})

	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"message":"error occured",
				"error": err.Error()})
	}

	return ctx.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message":"success",
			"content type": head.ContentType})
}

// @Summary Generate presigned GET URL
// @Description Create a presigned URL to download an object
// @Tags Presigned
// @Accept json
// @Produce json
// @Param bucket path string true "Bucket name"
// @Param request body request true "Object details"
// @Success 200 {object} map[string]interface{} "Presigned URL generated"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Router /presigned/getPresignedGetRequest/{bucket} [post]
func GetObject (ctx *fiber.Ctx) error {
	storage := dbconnection.S3Connect()
	presigner := beat.S3ConnectPresign(storage)

	bucket := ctx.Params("bucket")

	if bucket == "" {
		return ctx.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message":"specify bucket"})
	}

	req := request{}
	if err := ctx.BodyParser(&req); err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(
				&fiber.Map{
					"message":"unparsable entity",
					"error": err.Error()})
		}
	
	_, err := storage.S3Client.HeadObject(ctx.Context(), &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key: aws.String(req.UuidFileName),
	})

	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"bucket": bucket,
				"file": req.UuidFileName,
				"error": err.Error(),
				"message":"Couldn't get the file"})
	}
	
	presignedRequest, err := presigner.PresignClient.PresignGetObject(ctx.Context(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(req.UuidFileName),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(lifetimeSecs * int64(time.Second))
	})

	if err != nil {
		ctx.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"message":"Couldn't generate a presigned get request",
				"error": err,
				"bucket": bucket,
				"key": req.UuidFileName})
			return err
		}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message":"generated presigned get request successfully",
		"data": presignedRequest})
}

// @Summary Generate presigned PUT URL
// @Description Create a presigned URL to upload an object
// @Tags Presigned
// @Accept json
// @Produce json
// @Param bucket path string true "Bucket name"
// @Param request body request true "Object details"
// @Success 200 {object} map[string]interface{} "Presigned URL generated"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /presigned/getPresignedPostRequest/{bucket} [post]
func PutObject(ctx *fiber.Ctx) error {
	storage := dbconnection.S3Connect()
	presigner := beat.S3ConnectPresign(storage)

	bucket := ctx.Params("bucket")
	if bucket == "" {
		return ctx.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message":"specify bucket"})
	}

	req := request{}
	err := ctx.BodyParser(&req)
	if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(
				&fiber.Map{
					"message":"could not parse the json",
					"error": err.Error()})
		}
	
	contentType, err := typeOf(req.File)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"message":"wrong format",
				"error": err.Error()})
		}

	if bucket[0:3] != contentType[0:3]{
		return ctx.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"message":"wrong bucket",
				"bucket" : bucket,
				"recognised content Type": contentType})
	}

	objectKeyStr := req.UuidFileName

	presignedRequest, err := presigner.PresignClient.PresignPutObject(ctx.Context(), &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objectKeyStr),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = 15 * time.Minute
		// Добавляем условия подписи
	})

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{
				"message":"Couldn't presign url. Here's why: ",
				"error": err.Error(),
				"accepted_data": req})
	}

	ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message":"generated presigned put request successfully",
		"generatedObjectKey" : objectKeyStr,
		"data": presignedRequest})

	return nil
}

// @Summary Generate presigned DELETE URL
// @Description Create a presigned URL to delete an object
// @Tags Presigned
// @Accept json
// @Produce json
// @Param bucket path string true "Bucket name"
// @Param request body request true "Object details"
// @Success 200 {object} map[string]interface{} "Presigned URL generated"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Router /presigned/getPresignedDeleteRequest/{bucket} [post]
func DeleteObject(ctx *fiber.Ctx) error {
	storage := dbconnection.S3Connect()
	presigner := beat.S3ConnectPresign(storage)

	bucket := ctx.Params("bucket")
	if bucket == "" {
		return ctx.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message":"specify bucket"})
	}

	req := request{}
	err := ctx.BodyParser(&req)
	if err != nil {
		ctx.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message":"could not parse the json"})
			return err
		}

	//depending on the content-type should change the bucket name
	_, err = storage.S3Client.HeadBucket(ctx.Context(), &s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})

	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"message":"Couldn't get the bucket. Here's why: ",
				"error": err,
				"accepted_data": req})
	}

	presignedRequest, err := presigner.PresignClient.PresignDeleteObject(ctx.Context(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(req.UuidFileName),
	})
	
	if err != nil {
		ctx.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"message":"Couldn't generate a presigned delete request. Here's why: ",
				"error": err,
				"accepted_data": req})
			return err
		}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message":"generated presigned delete request successfully",
		"data": presignedRequest})
}
