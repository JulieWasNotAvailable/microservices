package license

import (
	"github.com/JulieWasNotAvailable/microservices/cart/api/presenters"
	"github.com/JulieWasNotAvailable/microservices/cart/internal/entities"
	"github.com/google/uuid"
)

type Service interface {
	GetAllLicenseTemplateByBeatmakerId(beatmakerId uuid.UUID) (*[]presenters.LicenseTemplate, error)
	ReadLicenseTemplateById(id uint) (*presenters.LicenseTemplate, error)
	UpdateLicenseTemplate(data presenters.LicenseTemplate) error
	InsertNewLicenseTemplate(userId uuid.UUID, data entities.LicenseTemplate) (entities.LicenseTemplate, error)

	GetLicenseByBeatId(beatId uuid.UUID) (*[]presenters.License, error)
	ReadLicenseById(id uint) (*presenters.License, error)
	InsertNewLicense(license entities.License) (entities.License, error)

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

func (s *service) ReadLicenseTemplateById(id uint) (*presenters.LicenseTemplate, error) {
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

func (s *service) InsertNewLicenseTemplate(userId uuid.UUID, data entities.LicenseTemplate) (entities.LicenseTemplate, error) {
	return s.repository.CreateNewLicenseTemplate(userId, data)
}

func (s *service) UpdateLicenseTemplate(data presenters.LicenseTemplate) error {
	return s.repository.UpdateLicenseTemplate(data)
}
