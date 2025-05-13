package beatmetadata

import (
	"github.com/JulieWasNotAvailable/microservices/unpublished/internal/entities"
	"github.com/google/uuid"
)

type MetadataService interface {
	CreateAvailableFiles(availableFiles *entities.AvailableFiles) (entities.AvailableFiles, error)
	ReadAllAvailableFiles() (*[]entities.AvailableFiles, error)
	ReadAvailableFilesByBeatId(uuid uuid.UUID) (*entities.AvailableFiles, error)
	UpdateAvailableFiles(availableFiles *entities.AvailableFiles) (entities.AvailableFiles, error)
	UpdateAvailableFilesByBeatId(beatId uuid.UUID, updateData *entities.AvailableFiles) (*entities.AvailableFiles, error)
	DeleteFileById(id uuid.UUID, fileType string) error

	CreateInstrument(instrument *entities.Instrument) (*entities.Instrument, error)
	GetAllInstruments() (*[]entities.Instrument, error)

	CreateGenre(genre *entities.Genre) (*entities.Genre, error)
	GetAllGenres() (*[]entities.Genre, error)

	CreateTimestamp(timestamp *entities.Timestamp) (*entities.Timestamp, error)
	GetAllTimestamps() (*[]entities.Timestamp, error)
	DeleteTimestampById(id uint) error

	CreateTag(tag *entities.Tag) (*entities.Tag, error)
	GetAllTags() (*[]entities.Tag, error)
	GetTagById(id uint) (*entities.Tag, error)
	GetTagsByName(name string) (*[]entities.Tag, error)
	DeleteTagById(id uint) error

	CreateMood(mood *entities.Mood) (*entities.Mood, error)
	GetAllMoods() (*[]entities.Mood, error)

	CreateKeynote(keynote *entities.Keynote) (*entities.Keynote, error)
	GetAllKeynotes() (*[]entities.Keynote, error)
}

type metadataService struct {
	repo MetadataRepository
}

func NewMetadataService(repo MetadataRepository) MetadataService {
	return &metadataService{repo: repo}
}

func (s *metadataService) CreateAvailableFiles(availableFiles *entities.AvailableFiles) (entities.AvailableFiles, error) {
	return s.repo.CreateAvailableFiles(availableFiles)
}

// ReadAllAvailableFiles implements MetadataService.
func (s *metadataService) ReadAllAvailableFiles() (*[]entities.AvailableFiles, error) {
	return s.repo.ReadAllAvailableFiles()
}

func (s *metadataService) ReadAvailableFilesByBeatId(uuid uuid.UUID) (*entities.AvailableFiles, error) {
	return s.repo.ReadAvailableFilesByBeatId(uuid)
}

// UpdateAvailableFiles implements MetadataService.
func (s *metadataService) UpdateAvailableFiles(availableFiles *entities.AvailableFiles) (entities.AvailableFiles, error) {
	return s.repo.UpdateAvailableFiles(availableFiles)
}

// UpdateAvailableFilesByBeatId implements MetadataService.
func (s *metadataService) UpdateAvailableFilesByBeatId(beatId uuid.UUID, updateData *entities.AvailableFiles) (*entities.AvailableFiles, error) {
	return s.repo.UpdateAvailableFilesByBeatId(beatId, updateData)
}

// DeleteFileById implements MetadataService.
func (s *metadataService) DeleteFileById(id uuid.UUID, fileType string) error {
	return s.repo.DeleteFileById(id, fileType)
}

func (s *metadataService) CreateInstrument(instrument *entities.Instrument) (*entities.Instrument, error) {
	return s.repo.CreateInstrument(instrument)
}

func (s *metadataService) GetAllInstruments() (*[]entities.Instrument, error) {
	return s.repo.ReadAllInstruments()
}

func (s *metadataService) CreateGenre(genre *entities.Genre) (*entities.Genre, error) {
	return s.repo.CreateGenre(genre)
}

func (s *metadataService) GetAllGenres() (*[]entities.Genre, error) {
	return s.repo.ReadAllGenres()
}

func (s *metadataService) CreateTimestamp(timestamp *entities.Timestamp) (*entities.Timestamp, error) {
	return s.repo.CreateTimestamp(timestamp)
}

func (s *metadataService) GetAllTimestamps() (*[]entities.Timestamp, error) {
	return s.repo.ReadAllTimestamps()
}

func (s *metadataService) DeleteTimestampById(id uint) error {
	return s.repo.DeleteTimestampById(id)
}

func (s *metadataService) CreateTag(tag *entities.Tag) (*entities.Tag, error) {
	return s.repo.CreateTag(tag)
}

func (s *metadataService) GetAllTags() (*[]entities.Tag, error) {
	return s.repo.ReadAllTags()
}

func (s *metadataService) GetTagById(id uint) (*entities.Tag, error) {
	return s.repo.ReadTagById(id)
}

func (s *metadataService) GetTagsByName(name string) (*[]entities.Tag, error){
	return s.repo.ReadTagsByName(name)
}

func (s *metadataService) DeleteTagById(id uint) error {
	return s.repo.DeleteTagById(id)
}

func (s *metadataService) CreateMood(mood *entities.Mood) (*entities.Mood, error) {
	return s.repo.CreateMood(mood)
}

func (s *metadataService) GetAllMoods() (*[]entities.Mood, error) {
	return s.repo.ReadAllMoods()
}

func (s *metadataService) CreateKeynote(keynote *entities.Keynote) (*entities.Keynote, error) {
	return s.repo.CreateKeynote(keynote)
}

func (s *metadataService) GetAllKeynotes() (*[]entities.Keynote, error) {
	return s.repo.ReadAllKeynotes()
}
