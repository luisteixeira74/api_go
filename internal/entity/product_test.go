package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProduct(t *testing.T) {
	product, err := NewProduct("Test Product", 10.0)
	assert.Nil(t, err)
	assert.NotNil(t, product)
	assert.NotEmpty(t, product.ID)
	assert.Equal(t, "Test Product", product.Name)
	assert.Equal(t, 10.0, product.Price)
}

func TestProductNameIsRequired(t *testing.T) {
	product, err := NewProduct("", 10.0)
	assert.Nil(t, product)
	assert.NotNil(t, err)
	assert.Equal(t, ErrNameIsRequired, err)
}

func TestProductPriceIsInvalid(t *testing.T) {
	product, err := NewProduct("Test Product", -10.0)
	assert.Nil(t, product) // valida que o produto não foi criado
	assert.NotNil(t, err)  // valida que houve erro
	assert.Equal(t, ErrInvalidPrice, err)
}

func TestProductValidate(t *testing.T) {
	product, err := NewProduct("Test Product", 10.0)
	assert.Nil(t, err)
	assert.NotNil(t, product)
	assert.NoError(t, product.Validate())
}
