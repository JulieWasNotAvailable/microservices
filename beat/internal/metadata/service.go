package metadata

import (
	"github.com/JulieWasNotAvailable/microservices/beat/api/presenters"
	"github.com/JulieWasNotAvailable/microservices/beat/internal/entities"
)

type Service interface {
	ReadAllGenres() (*[]entities.Genre, error)
	ReadPopularGenres() (*presenters.TrendingGenres, error)

	ReadAllMoods() (*[]entities.Mood, error)
	ReadAllKeys() (*[]entities.Keynote, error)
	ReadAllInstruments() (*[]entities.Instrument, error)

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

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{repository: r}
}

// ReadAllAvailableFiles implements Service.
func (s *service) ReadAllAvailableFiles() (*[]entities.AvailableFiles, error) {
	return s.repository.ReadAllAvailableFiles()
}

// ReadAllGenres implements Service.
func (s *service) ReadAllGenres() (*[]entities.Genre, error) {
	return s.repository.ReadAllGenres()
}

// ReadAllInstruments implements Service.
func (s *service) ReadAllInstruments() (*[]entities.Instrument, error) {
	return s.repository.ReadAllInstruments()
}

// ReadAllKeys implements Service.
func (s *service) ReadAllKeys() (*[]entities.Keynote, error) {
	return s.repository.ReadAllKeys()
}

// ReadAllMFCC implements Service.
func (s *service) ReadAllMFCC() (*[]entities.MFCC, error) {
	return s.repository.ReadAllMFCC()
}

// ReadAllMoods implements Service.
func (s *service) ReadAllMoods() (*[]entities.Mood, error) {
	return s.repository.ReadAllMoods()
}

// ReadAllTags implements Service.
func (s *service) ReadAllTags() (*[]entities.Tag, error) {
	return s.repository.ReadAllTags()
}

// ReadAllTimestamps implements Service.
func (s *service) ReadAllTimestamps() (*[]entities.Timestamp, error) {
	return s.repository.ReadAllTimestamps()
}

// ReadPopularGenres implements Service.
func (s *service) ReadPopularGenres() (*presenters.TrendingGenres, error) {
	return s.repository.ReadPopularGenres()
}

// ReadPopularTags implements Service.
func (s *service) ReadPopularTags() (*presenters.TrendingTags, error) {
	return s.repository.ReadPopularTags()
}

// ReadRandomTags implements Service.
func (s *service) ReadRandomTags() (*[]entities.Tag, error) {
	return s.repository.ReadRandomTags()
}

// ReadTagsByName implements Service.
func (s *service) ReadTagByName(name string) (*entities.Tag, error) {
	return s.repository.ReadTagByName(name)
}

func (s *service) ReadTagsByNameLike(name string) (*[]entities.Tag, error) {
	return s.repository.ReadTagsByNameLike(name)
}

