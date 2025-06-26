package metadata

import (
	"errors"
	"log"
	"time"

	"github.com/JulieWasNotAvailable/microservices/beat/api/presenters"
	"github.com/JulieWasNotAvailable/microservices/beat/internal/entities"
	"gorm.io/gorm"
)

type Repository interface {
	// GetAll(keynote *entities.Keynote) (*entities.Keynote, error) //getall аггрерирование на сторное api gateway
	ReadAllGenres() (*[]entities.Genre, error)
	ReadPopularGenres() (*presenters.TrendingGenres, error)

	ReadAllMoods() (*[]entities.Mood, error)
	ReadAllKeys() (*[]entities.Keynote, error)
	ReadGenresWithCounts() (entities.GenresWithCount, error)
	// ReadAllInstruments() (*[]entities.Instrument, error)

	ReadRandomTags() (*[]entities.Tag, error)
	ReadTagByName(name string) (*entities.Tag, error)
	ReadTagsByNameLike(name string) (*[]entities.Tag, error)
	ReadPopularTags() (*presenters.TrendingTags, error)

	//admin
	ReadAllTags() (*[]entities.Tag, error)
	ReadAllTimestamps() (*[]entities.Timestamp, error)
	ReadAllMFCC() (*[]entities.MFCC, error)
	ReadAllAvailableFiles() (*[]entities.AvailableFiles, error)
}

type repository struct {
	DB *gorm.DB
}

func NewRepo(db *gorm.DB) Repository {
	return &repository{
		DB: db,
	}
}

func (r *repository) ReadAllGenres() (*[]entities.Genre, error) {
	var genres []entities.Genre
	err := r.DB.Find(&genres).Error
	if err != nil {
		return nil, err
	}
	return &genres, nil
}

func (r *repository) ReadAllMoods() (*[]entities.Mood, error) {
	var moods []entities.Mood
	err := r.DB.Find(&moods).Error
	if err != nil {
		return nil, err
	}
	return &moods, nil
}

func (r *repository) ReadAllKeys() (*[]entities.Keynote, error) {
	var keys []entities.Keynote
	err := r.DB.Find(&keys).Error
	if err != nil {
		return nil, err
	}
	return &keys, nil
}

func (r *repository) ReadGenresWithCounts() (entities.GenresWithCount, error) {
	genresWithCount := entities.GenresWithCount{}

	err := r.DB.Table("genres").
	Select("genres.id, genres.name, count(*) as count, genres.picture_url").
	Joins("JOIN beat_genres ON genres.id = beat_genres.genre_id").
	Group("genres.id").
	Order("count DESC").
	Find(&genresWithCount).Error
	if err != nil {
		return entities.GenresWithCount{}, err
	}

	return genresWithCount, nil
}

// func (r *repository) ReadAllInstruments() (*[]entities.Instrument, error) {
// 	var instruments []entities.Instrument
// 	err := r.DB.Find(&instruments).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &instruments, nil
// }

func (r *repository) ReadAllTags() (*[]entities.Tag, error) {
	var tags []entities.Tag
	err := r.DB.Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return &tags, nil
}

func (r *repository) ReadAllTimestamps() (*[]entities.Timestamp, error) {
	var timestamps []entities.Timestamp
	err := r.DB.Find(&timestamps).Error
	if err != nil {
		return nil, err
	}
	return &timestamps, nil
}

func (r *repository) ReadAllMFCC() (*[]entities.MFCC, error) {
	var mfccs []entities.MFCC
	err := r.DB.Find(&mfccs).Error
	if err != nil {
		return nil, err
	}
	return &mfccs, nil
}

func (r *repository) ReadAllAvailableFiles() (*[]entities.AvailableFiles, error) {
	var files []entities.AvailableFiles
	err := r.DB.Find(&files).Error
	if err != nil {
		return nil, err
	}
	return &files, nil
}

func (r *repository) ReadRandomTags() (*[]entities.Tag, error) {
	var tags []entities.Tag
	err := r.DB.Order("RANDOM()").Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return &tags, nil
}

func (r *repository) ReadTagByName(name string) (*entities.Tag, error) {
	var tag entities.Tag
	err := r.DB.Where("name = ?", name).First(&tag).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = ErrMetadataNotFound
		}
		return nil, err
	}
	return &tag, nil
}

func (r *repository) ReadTagsByNameLike(name string) (*[]entities.Tag, error) {
	var tags []entities.Tag
	err := r.DB.Where("name LIKE ?", name+"%").Limit(10).Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return &tags, nil
}

func (r *repository) ReadPopularTags() (*presenters.TrendingTags, error) {
	trendingTags := presenters.TrendingTags{}

	startTime := time.Now().AddDate(0, -1, 0).Unix()
	endTime := time.Now().Unix()
	log.Println(startTime)
	log.Println(endTime)

	err := r.DB.Table("tags").
		Select("id, name, count(*) as count").
		Joins("JOIN beat_tags ON beat_tags.tag_id = tags.id").
		Group("id").Where("created_at BETWEEN ? AND ?", startTime, endTime).Order("count DESC").
		Find(&trendingTags).Error

	if err != nil {
		return nil, err
	}

	return &trendingTags, nil
}

// how many times was used for the last month
func (r *repository) ReadPopularGenres() (*presenters.TrendingGenres, error) {
	trendingGenres := presenters.TrendingGenres{}

	startTime := time.Now().AddDate(0, -1, 0).Unix()
	endTime := time.Now().Unix()

	err := r.DB.Table("genres").Select("id, name, count(*) as count").
		Joins("JOIN beat_genres ON beat_genres.genre_id = genres.id").
		Group("id").Where("created_at BETWEEN ? AND ?", startTime, endTime).Order("count DESC").
		Find(&trendingGenres).Error

	if err != nil {
		return nil, err
	}

	return &trendingGenres, nil
}
