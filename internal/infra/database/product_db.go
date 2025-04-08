package database

import (
	"apis/internal/entity"

	"gorm.io/gorm"
)

type Product struct {
	DB *gorm.DB
}

// Implement the ORM
func NewProduct(db *gorm.DB) *Product {
	return &Product{DB: db}
}

func (p *Product) Create(product *entity.Product) error {
	return p.DB.Create(product).Error
}

func (p *Product) GetByID(id string) (*entity.Product, error) {
	var product entity.Product
	err := p.DB.Where("id = ?", id).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (p *Product) Update(id string, product *entity.Product) error {
	_, err := p.GetByID(product.ID.String())
	if err != nil {
		return err
	}
	return p.DB.Save(product).Error
}

func (p *Product) Delete(id string) error {
	product, err := p.GetByID(id)
	if err != nil {
		return err
	}
	return p.DB.Delete(product).Error
}

func (p *Product) GetAll(page, limit int, sort string) ([]entity.Product, error) {
	var products []entity.Product

	// Sanitiza o valor de sort
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	// Prepara a query base
	query := p.DB.Order("created_at " + sort)

	// Aplica paginação se necessário
	if page > 0 && limit > 0 {
		query = query.Offset((page - 1) * limit).Limit(limit)
	}

	// Executa a consulta
	if err := query.Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}
