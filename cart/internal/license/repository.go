package license

import (
	"github.com/JulieWasNotAvailable/microservices/cart/api/presenters"
	"github.com/JulieWasNotAvailable/microservices/cart/internal/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	CreateNewLicenseTemplate(userId uuid.UUID, data entities.LicenseTemplate) (entities.LicenseTemplate, error)
	ReadLicenseTemplateById(id uint) (*presenters.LicenseTemplate, error)
	ReadAllLicenseTemplateByBeamakerId(beatmakerId uuid.UUID) (*[]presenters.LicenseTemplate, error)
	UpdateLicenseTemplate(data presenters.LicenseTemplate) error

	ReadLicenseByBeatId(beatId uuid.UUID) (*[]presenters.License, error)
	ReadLicenseById(id uint) (*presenters.License, error)
	CreateNewLicense(license entities.License) (entities.License, error)

	//admin
	ReadAllLicenseTemplate() (*[]presenters.LicenseTemplate, error)
	ReadAllLicense() (*[]presenters.License, error)
}

type repository struct {
	DB *gorm.DB
}

func NewRepo(db *gorm.DB) Repository {
	return &repository{
		DB: db,
	}
}

func (r *repository) CreateNewLicense(license entities.License) (entities.License, error) {
	result := r.DB.Create(&license)
	if result.Error != nil {
		return entities.License{}, result.Error
	}

	return license, nil
}

func (r *repository) CreateNewLicenseTemplate(userId uuid.UUID, data entities.LicenseTemplate) (entities.LicenseTemplate, error) {
	data.UserID = userId

	result := r.DB.Create(&data)
	if result.Error != nil {
		return entities.LicenseTemplate{}, result.Error
	}

	return data, nil
}

func (r *repository) ReadLicenseTemplateById(id uint) (*presenters.LicenseTemplate, error) {
	var licenseTemplateModel presenters.LicenseTemplate
	result := r.DB.Where("id = ?", id).First(&licenseTemplateModel)
	if result.Error == nil {
		return &licenseTemplateModel, result.Error
	}

	return &licenseTemplateModel, nil
}

func (r *repository) ReadAllLicenseTemplateByBeamakerId(beatmakerId uuid.UUID) (*[]presenters.LicenseTemplate, error) {
	var templates []presenters.LicenseTemplate

	result := r.DB.Model(&entities.LicenseTemplate{}).
		Where("user_id = ?", beatmakerId).
		Find(&templates)

	if result.Error != nil {
		return nil, result.Error
	}

	return &templates, nil
}

func (r *repository) ReadLicenseByBeatId(beatId uuid.UUID) (*[]presenters.License, error) {
	var licenses []presenters.License

	result := r.DB.Model(&entities.License{}).
		Where("beat_id = ?", beatId).
		Find(&licenses)

	if result.Error != nil {
		return nil, result.Error
	}

	return &licenses, nil
}

func (r *repository) ReadLicenseById(id uint) (*presenters.License, error) {
	license := presenters.License{}
	result := r.DB.Where("id = ?", id).First(&license)
	if result.Error != nil {
		return nil, result.Error
	}
	return &license, nil
}

func (r *repository) UpdateLicenseTemplate(data presenters.LicenseTemplate) error {
	result := r.DB.Model(&entities.LicenseTemplate{}).
		Where("id = ?", data.ID).
		Updates(data)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}


func (r *repository) ReadAllLicenseTemplate() (*[]presenters.LicenseTemplate, error) {
	var licenseTemplates []presenters.LicenseTemplate
	err := r.DB.Find(&licenseTemplates).Error
	if err != nil {
		return nil, err
	}
	return &licenseTemplates, nil
}

func (r *repository) ReadAllLicense() (*[]presenters.License, error) {
	var licenses []presenters.License
	err := r.DB.Find(&licenses).Error
	if err != nil {
		return nil, err
	}
	return &licenses, nil
}