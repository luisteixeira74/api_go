package database

import "apis/internal/entity"

type UserInterface interface {
	Create(user *entity.User) error
	GetByEmail(email string) (*entity.User, error)
}

type ProductInterface interface {
	Create(product *entity.Product) error
	GetAll(page, limit int, sort string) ([]entity.Product, error)
	GetByID(id string) (*entity.Product, error)
	Update(id string, product *entity.Product) error
	Delete(id string) error
}
