package activity

import (
	"errors"
	"fmt"

	"github.com/JulieWasNotAvailable/microservices/beat/internal/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	CreateLike(userid uuid.UUID, beatid uuid.UUID) (*entities.Like, error)
	DeleteLike(userid uuid.UUID, beatid uuid.UUID) (*entities.Like, error)
	ReadLikesByUserId(userid uuid.UUID) (*[]entities.Like, error)
	ReadLikesCountByBeatId(beatid uuid.UUID) (int, error) //likes number of this specific beat
	ReadLikesCountByUserId(userid uuid.UUID) (int, error) //likes number for a specific user (how many did he like)
	ReadLikesCountOfBeats(beatids []uuid.UUID) (int, error) //likes number of the group of beats

	CreateListened(userId uuid.UUID, beatId uuid.UUID) (entities.Listen, error)
	//view my listened
	BeatExists(beatId uuid.UUID) (bool, error) 
}

type repository struct {
	DB *gorm.DB
}

func NewRepo(db *gorm.DB) Repository {
	return &repository{
		DB: db,
	}
}

// CreateLike implements LikesRepository.
func (r *repository) CreateLike(userid uuid.UUID, beatid uuid.UUID) (*entities.Like, error) {
	beat := entities.Beat{
		ID: beatid,
	}
	err := r.DB.First(&beat).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("beat with this id does not exist")
	}
	
	like := entities.Like{
		UserID: userid,
		BeatID: beatid,
	}

	err = r.DB.Create(&like).Error
	if err != nil{
		return nil, err
	}

	return &like, nil
}

// DeleteLike implements LikesRepository.
func (r *repository) DeleteLike(userid uuid.UUID, beatid uuid.UUID) (*entities.Like, error) {
	result := r.DB.Where("user_id = ? AND beat_id = ?", userid, beatid).Delete(&entities.Like{})
    if result.Error != nil {
        return nil, result.Error
    }
    if result.RowsAffected == 0 {
        return nil, fmt.Errorf("no like found to delete")
    }
	like := entities.Like{
		UserID: userid,
		BeatID: beatid,
	}
    return &like, nil
}

// ReadLikesByUserId implements LikesRepository.
func (r *repository) ReadLikesByUserId(userid uuid.UUID) (*[]entities.Like, error) {
	var likes []entities.Like
    err := r.DB.Where("user_id = ?", userid).Preload("Beat").Find(&likes).Error
    if err != nil {
        return nil, err
    }
    return &likes, nil
}

// ReadLikesCountByBeatId implements LikesRepository.
func (r *repository) ReadLikesCountByBeatId(beatid uuid.UUID) (int, error) {
	var count int64
	beat := entities.Beat{
		ID: beatid,
	}
	err := r.DB.First(&beat).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, errors.New("beat with this id does not exist")
	}

    err = r.DB.Model(&entities.Like{}).Where("beat_id = ?", beatid).Count(&count).Error
    if err != nil {
        return 0, err
    }
    return int(count), nil
}

// ReadLikesCountByUserId implements LikesRepository.
func (r *repository) ReadLikesCountByUserId(userid uuid.UUID) (int, error) {
	var count int64
    err := r.DB.Model(&entities.Like{}).Where("user_id = ?", userid).Count(&count).Error
    if err != nil {
        return 0, err
    }
    return int(count), nil
}

func (r *repository) ReadLikesCountOfBeats(beatids []uuid.UUID) (int, error) {
	var count int64
    
    if len(beatids) == 0 {
        return 0, nil
    }

	//i don't think it's gonna work
    err := r.DB.Model(&entities.Like{}).
        Where("beat_id IN (?)", beatids).
        Count(&count).Error
        
    if err != nil {
        return 0, err
    }
    
    return int(count), nil
}

func (r *repository) CreateListened(userId uuid.UUID, beatId uuid.UUID) (entities.Listen, error) {
	record := entities.Listen{
		UserID: userId,
		BeatID: beatId,
	}

	err := r.DB.Create(&record).Error
	if err != nil{
		return entities.Listen{}, err
	}
	err = r.DB.First(&record).Error
	if err != nil{
		return entities.Listen{}, err
	}

	beat := &entities.Beat{
		ID: beatId,
	}
	_ = r.DB.First(&beat)
	beat.Plays = beat.Plays + 1
	err = r.DB.Updates(&beat).Error
	if err != nil {
		return entities.Listen{}, errors.New("listen record added, but faced an error while trying to increase the number of plays")
	}

	return record, nil
}



func (r *repository) BeatExists(beatId uuid.UUID) (bool, error) {
	err := r.DB.First(&entities.Beat{ID: beatId}).Error
	if errors.Is(err, gorm.ErrRecordNotFound){
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}