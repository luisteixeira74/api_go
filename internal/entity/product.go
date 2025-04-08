package entity

import (
	"apis/pkg/entity"
	"time"

	"github.com/pkg/errors"
)

var (
	ErrInvalidId       = errors.New("Invalid ID")
	ErrIDISRequired    = errors.New("ID is required")
	ErrNameIsRequired  = errors.New("Name is required")
	ErrPriceIsRequired = errors.New("Price is required")
	ErrInvalidPrice    = errors.New("Invalid price")
)

type Product struct {
	ID        entity.ID `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	CreatedAt string    `json:"created_at"`
}

func NewProduct(name string, price float64) (*Product, error) {
	return &Product{
		ID:        entity.NewID(),
		Name:      name,
		Price:     price,
		CreatedAt: time.Now().GoString(),
	}, nil
}

func (p *Product) Validate() error {
	if p.ID.String() == "" {
		return ErrIDISRequired
	}
	if _, err := entity.ParseID(p.ID.String()); err != nil {
		// ID is not valid
		return ErrInvalidId
	}
	if p.Name == "" {
		return ErrNameIsRequired
	}
	if p.Price == 0 {
		return ErrPriceIsRequired
	}
	if p.Price <= 0 {
		return ErrInvalidPrice
	}
	return nil
}
