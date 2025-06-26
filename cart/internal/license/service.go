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
	InsertNewLicenseList(beatId uuid.UUID, userId uuid.UUID, templatesAndPrices []entities.TemplateAndPrice) ([]entities.License, error) 

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

func (s *service) InsertNewLicenseList(beatId uuid.UUID, userId uuid.UUID, templatesAndPrices []entities.TemplateAndPrice) ([]entities.License, error) {
	licenses := []entities.License{}
	templateIds := []uint{}
	for _, tandp := range templatesAndPrices {
		template, err := s.repository.ReadLicenseTemplateById(tandp.TemplateId)
		if err != nil {
			return []entities.License{}, err
		}
		if template.UserID != userId {
			return []entities.License{}, errors.New("темплейт, который вы попытались использовать, не принадлежит вам")
		}
		if contains(templateIds, template.ID) {
			return []entities.License{}, errors.New("данный шаблон лицензии уже используется в бите")
		}
		templateIds = append(templateIds, tandp.TemplateId)
		license := entities.License{
			UserID: userId,
			BeatID: beatId,
			Price: tandp.Price,
			LicenseTemplateID: tandp.TemplateId,
		}
		licenses = append(licenses, license)
	}

	createdLicenses, err := s.repository.CreateNewLicenseList(beatId, userId, licenses)
	if err != nil {
		return []entities.License{}, err
	}
	return createdLicenses, nil
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
