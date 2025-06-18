package license

import (
	"errors"

	"github.com/JulieWasNotAvailable/microservices/cart/api/presenters"
	"github.com/JulieWasNotAvailable/microservices/cart/internal/entities"
	"github.com/google/uuid"
)

type Service interface {
	GetAllLicenseTemplateByBeatmakerId(beatmakerId uuid.UUID) (*[]presenters.LicenseTemplate, error)
	ReadLicenseTemplateById(id uint) (presenters.LicenseTemplate, error)
	UpdateLicenseTemplate(data presenters.LicenseTemplate) error
	InsertNewLicenseTemplate(userId uuid.UUID, data entities.LicenseTemplate) (entities.LicenseTemplate, error)

	GetLicenseByBeatId(beatId uuid.UUID) (*[]presenters.License, error)
	ReadLicenseById(id uint) (*presenters.License, error)
	InsertNewLicense(license entities.License) (entities.License, error)
	InsertNewLicenseList(beatId uuid.UUID, userId uuid.UUID, data []entities.License) ([]entities.License, error)

	//admin
	ReadAllLicenseTemplate() (*[]presenters.LicenseTemplate, error)
	ReadAllLicense() (*[]presenters.License, error)
}

type service struct {
	repository Repository
}

func (s *service) ReadAllLicense() (*[]presenters.License, error) {
	return s.repository.ReadAllLicense()
}

func (s *service) ReadAllLicenseTemplate() (*[]presenters.LicenseTemplate, error) {
	return s.repository.ReadAllLicenseTemplate()
}


func NewService(r Repository) Service {
	return &service{repository: r}
}

func (s *service) ReadLicenseTemplateById(id uint) (presenters.LicenseTemplate, error) {
	return s.repository.ReadLicenseTemplateById(id)
}

func (s *service) GetAllLicenseTemplateByBeatmakerId(beatmakerId uuid.UUID) (*[]presenters.LicenseTemplate, error) {
	return s.repository.ReadAllLicenseTemplateByBeamakerId(beatmakerId)
}

func (s *service) GetLicenseByBeatId(beatId uuid.UUID) (*[]presenters.License, error) {
	return s.repository.ReadLicenseByBeatId(beatId)
}

func (s *service) ReadLicenseById(id uint) (*presenters.License, error) {
	return s.repository.ReadLicenseById(id)
}

func (s *service) InsertNewLicense(license entities.License) (entities.License, error) {
	return s.repository.CreateNewLicense(license)
}

func (s *service) InsertNewLicenseList(beatId uuid.UUID, userId uuid.UUID, licenses []entities.License) ([]entities.License, error) {
	templateIds := []uint{}
	for _, license := range(licenses){
		if contains(templateIds, license.LicenseTemplateID) {
			return []entities.License{}, errors.New("license template was already used in this beat")
		}
		templateIds = append(templateIds, license.LicenseTemplateID)

		template, err := s.repository.ReadLicenseTemplateById(license.LicenseTemplateID)
		if template.UserID != userId {
			return []entities.License{}, errors.New("user does not own the template")
		}
		if err != nil{
			return []entities.License{}, err
		}
	}

	licenses, err := s.repository.CreateNewLicenseList(beatId, userId, licenses)
	if err != nil {
		return []entities.License{}, err
	}
	return licenses, nil
}

func (s *service) InsertNewLicenseTemplate(userId uuid.UUID, data entities.LicenseTemplate) (entities.LicenseTemplate, error) {
	return s.repository.CreateNewLicenseTemplate(userId, data)
}

func (s *service) UpdateLicenseTemplate(data presenters.LicenseTemplate) error {
	return s.repository.UpdateLicenseTemplate(data)
}

func contains(slice []uint, item uint) bool {
    for _, v := range slice {
        if v == item {
            return true
        }
    }
    return false
}
