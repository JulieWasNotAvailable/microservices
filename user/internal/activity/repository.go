package activity

import (
	"errors"

	"github.com/JulieWasNotAvailable/microservices/user/internal/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	CreateSub(userId uuid.UUID, beatmakerId uuid.UUID) (entities.User_Follows_Beatmakers, error)
	ReadSubsByUserId(userId uuid.UUID) ([]entities.User_Follows_Beatmakers, error)
	ReadSubsCountByBeatmakerId(beatmakerId uuid.UUID) (int, error)
	DeleteSub(userId uuid.UUID, beatmakerId uuid.UUID) (entities.User_Follows_Beatmakers, error)
}

type repository struct {
	DB *gorm.DB
}

func NewRepo(db *gorm.DB) Repository {
	return &repository{
		DB: db,
	}
}

func (r *repository) CreateSub(userId uuid.UUID, beatmakerId uuid.UUID) (entities.User_Follows_Beatmakers, error) {
	var user entities.User
	err := r.DB.Where("role_id = ?", 2).Where("id = ?", beatmakerId).First(&user).Error
	if err != nil && errors.As(err, &gorm.ErrRecordNotFound) {
		return entities.User_Follows_Beatmakers{}, errors.New("the beatmaker you tried to follow does not exist, or does not have beatmaker role")
	}

	sub := entities.User_Follows_Beatmakers{
		UserID: userId,
		BeatmakerID: beatmakerId,
	}
	err = r.DB.Create(&sub).Error
	if err != nil{
		return entities.User_Follows_Beatmakers{}, err
	}

	return sub, nil
}

// DeleteSub implements Repository.
func (r *repository) DeleteSub(userId uuid.UUID, beatmakerId uuid.UUID) (entities.User_Follows_Beatmakers, error) {
	sub := entities.User_Follows_Beatmakers{
		UserID: userId,
		BeatmakerID: beatmakerId,
	}
	result := r.DB.Delete(sub)
	if result.Error != nil {
		return entities.User_Follows_Beatmakers{}, result.Error
	}

	if result.RowsAffected == 0 {
		return entities.User_Follows_Beatmakers{}, errors.New("sub was not found")
	}

	return sub, nil
}

// ReadSubsByUserId implements Repository.
func (r *repository) ReadSubsByUserId(userId uuid.UUID) ([]entities.User_Follows_Beatmakers, error) {
	subs := []entities.User_Follows_Beatmakers{}
	err := r.DB.Preload("Beatmaker").Where("user_id = ?", userId).Find(&subs).Error
	if err != nil{
		return []entities.User_Follows_Beatmakers{}, err
	}
	return subs, nil
}

// ReadSubsCountByBeatmakerId implements Repository.
func (r *repository) ReadSubsCountByBeatmakerId(beatmakerId uuid.UUID) (int, error) {
	var count int64
	err := r.DB.Model(&entities.User_Follows_Beatmakers{}).Where("beatmaker_id = ?", beatmakerId).Count(&count).Error
	if err != nil{
		return 0, err
	}

    return int(count), nil
}


