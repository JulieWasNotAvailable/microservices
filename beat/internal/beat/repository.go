package beat

import (
	// "github.com/JulieWasNotAvailable/microservices/beat/pkg/entities"
	"errors"
	"strings"

	"github.com/JulieWasNotAvailable/microservices/beat/api/presenters"
	"github.com/JulieWasNotAvailable/microservices/beat/internal/entities"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

type Repository interface {
	//publish
	CreateBeat(beat entities.Beat) (entities.Beat, error)

	// helpers
	ReadBeats() (*[]entities.Beat, error)
	ReadBeatById(id uuid.UUID) (*presenters.Beat, error)
	ReadBeatsByBeatmakerId(id uuid.UUID) (*[]presenters.Beat, error)

	//filters
	ReadFilteredBeats(filters presenters.Filters) (*[]presenters.Beat, error)
	ReadBeatsByMoodId(moodId uint) (*[]presenters.Beat, error)
	ReadBeatsByDate(from int64, to int64) (*[]presenters.Beat, error)
	FindBeatsWithAllMoods(moodIDs []uint) ([]presenters.Beat, error)

	//for service
	DeleteBeatById(id uuid.UUID) error 
}

type repository struct {
	DB *gorm.DB
}

func NewRepo(db *gorm.DB) Repository {
	return &repository{
		DB: db,
	}
}

func (r repository) CreateBeat(beat entities.Beat) (entities.Beat, error) {
	result := r.DB.Create(&beat)

	if result.Error != nil {
		return entities.Beat{}, result.Error
	}

	return beat, nil
}

func (r *repository) ReadBeats() (*[]entities.Beat, error) {
	var beatModels []entities.Beat

	result := r.DB.Model(beatModels).
	Preload("AvailableFiles").
	Preload("Tags").Preload("Genres").
	Preload("Moods").Preload("Timestamps").
	Preload("Instruments").Preload("MFCC").
	Find(&beatModels)
	if result.Error != nil {
		return nil, result.Error
	}

	return &beatModels, nil
}

func (r *repository) ReadBeatById(id uuid.UUID) (*presenters.Beat, error) {
	var beatModel presenters.Beat

	result := r.DB.Model(&beatModel).Preload("Tags").Preload("Genres").
		Preload("Moods").Preload("Timestamps").Preload("Instruments").
		Where("id = ?", id).First(&beatModel)
	if result.Error != nil {
		return nil, result.Error
	}

	return &beatModel, nil
}

func (r *repository) ReadBeatsByBeatmakerId(id uuid.UUID) (*[]presenters.Beat, error) {
	var beatModels []presenters.Beat

	result := r.DB.Model(beatModels).Preload("Tags").Preload("Genres").
		Preload("Moods").Preload("Timestamps").Preload("Instruments").
		Where("beatmaker_id = ?", id).Find(&beatModels)
	if result.Error != nil {
		return nil, result.Error
	}

	return &beatModels, nil
}

func (r *repository)ReadBeatsByMoodId(moodId uint) (*[]presenters.Beat, error){
	var beatModels []presenters.Beat

	err := r.DB.Table("beats").Joins("JOIN beat_moods ON beat_moods.beat_id = beats.id").
	Where("mood_id = ?", moodId).
	Scopes(WithBasicPreloads()).
	Find(&beatModels).Error
	if err != nil {
		return nil, err
	}

	return &beatModels, nil
}

func (r *repository)ReadBeatsByDate(from int64, to int64) (*[]presenters.Beat, error){
	var beatModels []presenters.Beat
	err := r.DB.
	Where("created_at BETWEEN ? AND ?", from, to).
	Find(&beatModels).Error
	if err != nil{
		return nil, err
	}
	return &beatModels, nil
}

func (r *repository) DeleteBeatById(id uuid.UUID) error {
	beat := entities.Beat{}
	result := r.DB.Where("id = ?", id).Delete(beat)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("unpublished beat not found")
	}

	return nil
}

func (r *repository) FindBeatsWithAllMoods(moodIDs []uint) ([]presenters.Beat, error) {
    var beats []presenters.Beat
	// err := r.DB.Scopes(WithAllMoods(moodIDs)).Find(&beats).Error
	// if err != nil{
	// 	return nil, err
	// }
	return beats, nil
}

func (r *repository) ReadFilteredBeats(filters presenters.Filters) (*[]presenters.Beat, error){
	var beats []presenters.Beat
	err := r.DB.Scopes(WithAllMoodsGenres(filters.Moods, filters.Genres, filters.Tags)).Find(&beats).Error
	if err != nil{
		return nil, err
	}
	
	result, err := FilterBeatsByNumericFields(beats, filters)
	if err != nil {
		return nil, err
	}
	
	return result, nil
}

func FilterBeatsByNumericFields(beats []presenters.Beat, filters presenters.Filters) (*[]presenters.Beat, error) {
    var filteredBeats []presenters.Beat

    for _, beat := range beats {
        match := true

		if filters.Keynote != 0 && beat.KeynoteID != filters.Keynote {
            match = false
        }

        if filters.MinPrice != 0 && beat.Price < filters.MinPrice {
            match = false
        }

        if filters.MaxPrice != 0 && beat.Price > filters.MaxPrice {
            match = false
        }

        if filters.MinBPM != 0 && beat.BPM < filters.MinBPM {
            match = false
        }

        if filters.MaxBPM != 0 && beat.BPM > filters.MaxBPM {
            match = false
        }

        if match {
            filteredBeats = append(filteredBeats, beat)
        }
    }

    if len(filteredBeats) == 0 {
        return nil, errors.New("no beats match the specified numeric filters")
    }

    return &filteredBeats, nil
}

//scopes

func WithAllMoodsGenres(moodIDs []uint, genreIDs []uint, tagIDs []uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(moodIDs) == 0 && len(genreIDs) == 0 && len(tagIDs) == 0{
			return db.Where("1 = 1")
		}

		allparams := append(moodIDs, genreIDs...)
		allparams = append(allparams, tagIDs...)
		args := make([]interface{}, len(allparams))
		for i, id := range allparams {
			args[i] = id
		}

		var builder strings.Builder
		
		builder.Write([]byte("WITH "))
		isFirst := true
		if len(moodIDs) != 0{
			builder = *moodBuilder(&builder, moodIDs)
			isFirst = false
		}
		if len(genreIDs) != 0{
			if !isFirst{
				builder.Write([]byte(")), "))
			}
			builder = *genreBuilder(&builder, genreIDs)
		}
		if len(tagIDs) != 0{
			if !isFirst{
				builder.Write([]byte(")), "))
			}
			builder = *tagBuilder(&builder, tagIDs)
		}
		builder.Write([]byte(")) "))
		
		isFirst = true
		builder.Write([]byte("SELECT * FROM beats b WHERE "))
		if len(moodIDs) != 0{
			builder.Write([]byte("b.id in (SELECT id FROM moods_CTE)"))
			isFirst = false
		}
		if len(genreIDs) != 0 {
			if isFirst{
				builder.Write([]byte("b.id in (SELECT id FROM genres_CTE)"))	
			} else {
				builder.Write([]byte(" AND b.id IN (SELECT id FROM genres_CTE)"))
			}
			isFirst = false
		}
		if len(tagIDs) != 0 {
			if isFirst{
				builder.Write([]byte("b.id IN (SELECT id FROM tags_CTE)"))
			} else {
				builder.Write([]byte(" AND b.id IN (SELECT id FROM tags_CTE)"))
			}
		}

		query := builder.String()
		return db.Raw(query, args...)
	}
}

func moodBuilder(builder *strings.Builder, moodIDs []uint) *strings.Builder {
	builder.Write([]byte("moods_CTE AS (SELECT b.* FROM beats b WHERE b.id IN ("))
	for n := range moodIDs {
		if n == 0{
			builder.Write([]byte("SELECT beat_id FROM beat_moods WHERE mood_id in (?) "))		
		} else {
			builder.Write([]byte("INTERSECT SELECT beat_id FROM beat_moods WHERE mood_id in (?) "))
		}
	}
	
	return builder	
}

func genreBuilder(builder *strings.Builder, genreIDs []uint) *strings.Builder {
	builder.Write([]byte("genres_CTE as (select beats.id from beats where beats.id in ("))
		for n := range genreIDs{
			if n == 0{
				builder.Write([]byte("SELECT beat_id FROM beat_genres WHERE genre_id in (?) "))		
			} else {
				builder.Write([]byte("INTERSECT SELECT beat_id FROM beat_genres WHERE genre_id in (?) "))
			}
		}
	return builder
}

func tagBuilder (builder *strings.Builder, tagIDs []uint) *strings.Builder {
	builder.Write([]byte("tags_CTE as (select beats.id from beats where beats.id in ("))
	for n := range tagIDs{
		if n == 0{
			builder.Write([]byte("SELECT beat_id FROM beat_tags WHERE tag_id in (?) "))		
		} else {
			builder.Write([]byte("INTERSECT SELECT beat_id FROM beat_tags WHERE tag_id in (?) "))
		}
	}
	return builder
}

func WithAllGenres(genreIDs []uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(genreIDs) == 0 {
			return db.Where("1 = 0") // Return empty result if no genres provided
		}

		args := make([]interface{}, len(genreIDs))
		for i, id := range genreIDs {
			args[i] = id
		}

		// Build INTERSECT query dynamically
		query := `
			WITH genres_CTE AS (SELECT b.id FROM beats b 
			WHERE b.id IN (
				SELECT beat_id FROM beat_genres WHERE genre_id = ? 
				` + strings.Repeat("INTERSECT SELECT beat_id FROM beat_genres WHERE genre_id = ? ", len(genreIDs)-1) + `
			))`

		return db.Raw(query, args...)
	}
}

func WithPriceMax(max *int) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        if max != nil {
            return db.Where("price <= ?", *max)
        }
        return db
    }
}

func WithPriceMin(min *int) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        if min != nil {
            return db.Where("price >= ?", *min)
        }
        return db
    }
}

func WithBasicPreloads() func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
		return db.
			Preload("Tags").Preload("Genres").
			Preload("Moods").Preload("Timestamps").
			Preload("Instruments")
	}
}
