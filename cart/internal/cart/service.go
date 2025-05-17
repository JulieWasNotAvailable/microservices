package cart

import (
	"errors"
	"log"

	"github.com/JulieWasNotAvailable/microservices/cart/api/presenters"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	CreateProductToCart(userId uuid.UUID, licenseId uint) error
	ReadCartByUser(userId uuid.UUID) (presenters.Cart, error)
	DeleteFromCart(userId uuid.UUID, licenseId uint) error
}

type service struct {
	repository Repository
}

// CreateProductToCart implements Service.
func (s *service) CreateProductToCart(userId uuid.UUID, licenseId uint) error {
	cart, err := s.repository.ReadCartByUser(userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newcart, err := s.repository.CreateCart(userId)
			if err != nil{
				return err
			}
			
			log.Println("new cart: ", newcart.ID)
			return s.repository.AddProductToCart(newcart.ID, licenseId)
		}
		return err
	}

	return s.repository.AddProductToCart(cart.ID, licenseId)
}

// DeleteFromCart implements Service.
func (s *service) DeleteFromCart(userId uuid.UUID, licenseId uint) error {
	cart, err := s.repository.ReadCartByUser(userId)
	if err != nil {
		return err
	}
	return s.repository.DeleteFromCart(cart.ID, licenseId)
}

// ReadCartByUser implements Service.
func (s *service) ReadCartByUser(userId uuid.UUID) (presenters.Cart, error) {
	return s.repository.ReadCartByUser(userId)
}

func NewService(r Repository) Service {
	return &service{repository: r}
}
