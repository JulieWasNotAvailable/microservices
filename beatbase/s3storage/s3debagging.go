package s3storage

import (
	// "context"
	// "errors"
	// "log"

	// "github.com/aws/aws-sdk-go-v2/aws"
	// "github.com/aws/aws-sdk-go-v2/config"
	// "github.com/aws/aws-sdk-go-v2/service/s3"
	// "github.com/aws/aws-sdk-go-v2/service/s3/types"
	// "github.com/aws/smithy-go"
	// // "github.com/gofiber/storage/s3/v2"
)

// func S3Connect() Storage {
//     //по примеру из яндекса
//     // Подгружаем конфигурацию из user/1/.aws/config*
//     cfg, err := config.LoadDefaultConfig(context.TODO())
//     if err != nil {
//         log.Fatal(err)
//     }

//     // Создаем клиента для доступа к хранилищу S3
//     client := s3.NewFromConfig(cfg)

//     storage := Storage{
//         S3Client: client,
//     }
//     return storage
// }

// func (storage Storage) ListObjectsByBucket (ctx context.Context, BucketName string) ([]types.Object, error) {
// 	var err error
// 	var output *s3.ListObjectsV2Output
// 	input := &s3.ListObjectsV2Input{
// 		Bucket: aws.String(BucketName),
// 	}
// 	var objects []types.Object
// 	objectPaginator := s3.NewListObjectsV2Paginator(storage.S3Client, input)
//     //пагинатор нужен для того, чтобы, когда объектов много, возвращать их в чанках через несколько запросов, а не выводить все разом
// 	for objectPaginator.HasMorePages() {
// 		output, err = objectPaginator.NextPage(ctx)
// 		if err != nil {
// 			var noBucket *types.NoSuchBucket
// 			if errors.As(err, &noBucket) {
// 				log.Printf("Bucket %s does not exist.\n", BucketName)
// 				err = noBucket
// 			}
// 			break
// 		} else {
// 			objects = append(objects, output.Contents...)
//             log.Println(objects)
// 		}
// 	}
//     log.Println("here are the objects: ", objects[0])
//     for _, object := range objects{
//         log.Println(aws.ToString(object.Key))
//     }
// 	return objects, err
// }

// func (storage Storage) GetMP3sByBucket(ctx context.Context, BucketName string) (error) {

//     output, err := storage.S3Client.ListObjectsV2(
//         ctx,
//         &s3.ListObjectsV2Input{
//             Bucket: aws.String(BucketName),
//         })

//     if err != nil{
//         log.Println(err)
//     }
//     log.Println("found the objects")
    
//     // m := make(map[string]int)
//     log.Println("first page results")
//     for _, object := range output.Contents {
//         // m[aws.ToString(object.Key)] = int(*object.Size)
//         log.Printf("key=%s size=%d", aws.ToString(object.Key), object.Size)
//     }

//     // ctx.Status(http.StatusOK).JSON(
//     //     &fiber.Map{"files in  bucket": m})
//     return nil
// }

// func (s Storage) UploadMP3(ctx context.Context, object Object) error {

//     file, err := os.Open(object.FileName)

//     if err != nil {
// 		log.Printf("Couldn't open file %v to upload. Here's why: %v\n", object.FileName, err)
// 	} else {
//         defer file.Close()
//         _, err = s.S3Client.PutObject(ctx, &s3.PutObjectInput{
//             Bucket: aws.String(object.BucketName),
//             Key: aws.String(object.ObjectKey),
//             Body: file,
//         })

//         //smithy is an IDL - interface description language
//         //AWS SDK для Go использует smithy для стандартизации ошибок. Здесь в коде нужен чтобы унифицировать ошибки и получить больше информации о них
//         if err != nil {
//             var apiErr smithy.APIError
//             //отлавливаем ошибку entity is too large
//             if errors.As(err, &apiErr) && apiErr.ErrorCode() == "EntityTooLarge" {
//                 log.Printf("Error while uploading object to %s. The object is too large.\n"+
// 					"To upload objects larger than 5GB, use the S3 console (160GB max)\n"+
// 					"or the multipart upload API (5TB max).", object.BucketName)
//             } else {
// 				log.Printf("Couldn't upload file %v to %v:%v. Here's why: %v\n",
// 					object.FileName, object.BucketName, object.ObjectKey, err)
// 			}
//         } else {
//             //ждёт ObjectExistsWaiter используется для ожидания, пока загруженный объект не станет доступен.
//             //Wait периодически проверяет, загрузился объект или нет
//             // headobject возвращает метаданные загруженного объекта, если он существует, и ошибку, если не сущесвует
//             // wait вызывает его периодически, интервалы заданы в SDK
            
//             err = s3.NewObjectExistsWaiter(s.S3Client).Wait(
//             ctx, &s3.HeadObjectInput{
//                 Bucket: aws.String(object.BucketName),
//                 Key: aws.String(object.ObjectKey)}, time.Minute)

//             if err != nil {
//                     log.Printf("Failed attempt to wait for object %s to exist.\n", object.ObjectKey)
//                 } else {
//                     log.Println("file uploaded successfully")
//                 }
//         }
//     }
//     return nil    
// }