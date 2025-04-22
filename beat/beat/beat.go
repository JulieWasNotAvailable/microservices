package beat

import (
	"fmt"
	"net/http"
	"time"

	// "github.com/JulieWasNotAvailable/goBeatsBackend/beat"
	"github.com/JulieWasNotAvailable/goBeatsBackend/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Beat struct{
	ID uint
	Author *string
	Title *string
	License *string
	Mood *string
	Date time.Time
	Genre *string
	Url *string
	FreeForNonProfit *uint
}

type BeatByDateRequest struct {
	Date1 *string //pointer?
	Date2 *string 
}

type BeatByFeatureRequest struct {
	Feature *string //pointer?
	Value *string 
}

type UpdateBeatURLRequest struct {
	BeatId *uint
	Url *string
}

type Repository struct{
	DB *gorm.DB
}

func(r *Repository) CreateBeat (context *fiber.Ctx) error {
	beat := Beat{}

	beat.Date = time.Now()

	err := context.BodyParser(&beat) //binds a request body to a struct

	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	result := r.DB.Create(&beat)	
	
	if result.Error != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message":"could not create beat"})
			return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message":"beat has been added",
			"beat ID":beat.ID})

	return nil
}

func(r *Repository) UpdateBeatURL (req *UpdateBeatURLRequest) error {

	result := r.DB.Model(Beat{}).Where("id = ?", req.BeatId).Updates(Beat{Url: req.Url})

	fmt.Println(result, " ", result.Error)

	return nil
}

func (r *Repository) GetBeats(context *fiber.Ctx) error {
	beatModels := &[]model.Beat{}

	err:= r.DB.Find(beatModels).Error
	if err!=nil{
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message":"could not get books"})
			return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message":"beats fetched successfully",
		"data": beatModels})

	return nil
}

func (r *Repository) DeleteBeats (context *fiber.Ctx) error {
	//with fiber you can easily access params, request and response data
	//but you should be able to do stuff without fiber
	beatModel := model.Beat{}
	id := context.Params("id")

	if id == ""{
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
	}
	

	err := r.DB.Delete(beatModel, id).Error

	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
		&fiber.Map{"message": "couldn't delete"})
		return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "book deleted successfully"})

	return nil
}

func (r *Repository) GetBeatByID(context *fiber.Ctx) error {
	id := context.Params("id")
	beatModel :=  &model.Beat{}
	if id == ""{
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message":"id cannot be empty",
		})
	}
	fmt.Println("the id is: ", id)

	err:= r.DB.Where("id = ?", id).First(beatModel).Error
	if err !=nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get the beat"})
			return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "beat is fetched successfully",
		"data": beatModel,
	})

	return nil
}

// db.Where("created_at BETWEEN ? AND ?", lastWeek, today).Find(&users)
// SELECT * FROM users WHERE created_at BETWEEN '2000-01-01 00:00:00' AND '2000-01-08 00:00:00';

func (r Repository) GetBeatsByDate (context *fiber.Ctx) error {
	beatModels := &[]model.Beat{}
	request := &BeatByDateRequest{}

	err := context.BodyParser(&request) //binds a request body to a struct

	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message":"could not parse the json"})
			return err
		}
	
	err = r.DB.Where("Date BETWEEN ? AND ?", request.Date1, request.Date2).Find(beatModels).Error

	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message":"could not get the beats"})
			return err
		}
	
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "selected beats fetched successfully",
		"data": beatModels,
	})
	
	return nil
}

// func (r Repository) GetBeatsByFeature (context *fiber.Ctx){
// 	beat := &Beat{}
// 	beatModels := &[]model.Beat{}
// 	request := &BeatByFeatureRequest{}

// 	err := context.BodyParser(&request) //binds a request body to a struct

// 	if err != nil {
// 		context.Status(http.StatusBadRequest).JSON(
// 			&fiber.Map{"message":"could not parse the json"})
// 			return err
// 		}
	
// 	// feature := request.Feature
// 	//regular expressions?
// 	value := request.Value
// 	err = r.DB.Where(Beat{feature: value}).Find(beatModels).Error

// 	if err != nil {
// 		context.Status(http.StatusBadRequest).JSON(
// 			&fiber.Map{"message":"could not get the beats"})
// 			return err
// 		}
	
// 	context.Status(http.StatusOK).JSON(&fiber.Map{
// 		"message": "selected beats fetched successfully",
// 		"data": beatModels,
// 	})
	
// 	return nil
// }

func (r Repository) SetupRoutes (app *fiber.App){
	api := app.Group("/api/beats")
	api.Post("/create_beat", r.CreateBeat)
	api.Delete("/delete_beat/:id", r.DeleteBeats)
	api.Get("/get_beat/:id", r.GetBeatByID)
	api.Get("/beats", r.GetBeats)	
	api.Get("/get_beats_by_date", r.GetBeatsByDate)
	// api.Put("/updateBeatUrl/:id", r.UpdateBeatURL)
}