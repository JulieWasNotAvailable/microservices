package cart

import (
	"github.com/JulieWasNotAvailable/microservices/cart/api/presenters"
	"github.com/JulieWasNotAvailable/microservices/cart/internal/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	CreateCart(userId uuid.UUID) (*entities.Cart, error)
	AddProductToCart(cartId uuid.UUID, licenseId uint) error
	ReadCartByUser(userId uuid.UUID) (presenters.Cart, error)
	DeleteFromCart(userId uuid.UUID, licenseId uint) error
}

type repository struct {
	DB *gorm.DB
}

func NewRepo(db *gorm.DB) Repository {
	return &repository{
		DB: db,
	}
}

func (r *repository) CreateCart(userId uuid.UUID) (*entities.Cart, error) {
	uuid, err := uuid.NewV7()
	if err != nil{
		return nil, err	
	}
	cart := entities.Cart{
		ID : uuid,
		UserID: userId,
	}
	result := r.DB.Create(&cart)
	if result.Error != nil {
		return nil, result.Error
	}

	return &cart, nil
}

func (r *repository) AddProductToCart(cartId uuid.UUID, licenseId uint) error {
	cartItem := entities.CartLicenses{
		CartID:    cartId,
		LicenseID: licenseId,
	}
	
	result := r.DB.Create(&cartItem)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *repository) ReadCartByUser(userId uuid.UUID) (presenters.Cart, error) {
	var cart presenters.Cart

	result := r.DB.Where("user_id = ?", userId).Preload("Licenses").
		First(&cart)

	if result.Error != nil {
		return presenters.Cart{}, result.Error
	}

	return cart, nil
}

func (r *repository) DeleteFromCart(cartId uuid.UUID, licenseId uint) error {
	result := r.DB.Where("cart_id = ? AND license_id = ?", cartId, licenseId).
		Delete(&entities.CartLicenses{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
