package database

import (
	"apis/internal/entity"
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestNewProduct(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}

	db.Migrator().DropTable(&entity.Product{})
	db.AutoMigrate(&entity.Product{})

	product, err := entity.NewProduct("Test Product", 10.0)
	assert.NoError(t, err)
	productDB := NewProduct(db)
	err = productDB.Create(product)
	assert.NoError(t, err)
	assert.NotEmpty(t, product.ID)
}

func TestAllProducts(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}

	// Limpa completamente a tabela antes de começar o teste
	db.Migrator().DropTable(&entity.Product{})
	db.AutoMigrate(&entity.Product{})

	for i := 1; i < 24; i++ {
		product, err := entity.NewProduct(fmt.Sprintf("Product %d", i), rand.Float64()*100)
		assert.NoError(t, err)
		db.Create(product)
	}

	productDB := NewProduct(db)
	// Testando a paginação
	products, err := productDB.GetAll(1, 10, "asc")
	fmt.Println(products[0].Name)
	assert.NoError(t, err)
	assert.Len(t, products, 10)
	assert.Equal(t, "Product 1", products[0].Name)
	assert.Equal(t, "Product 10", products[9].Name)

	products, err = productDB.GetAll(2, 10, "asc")
	assert.NoError(t, err)
	assert.Len(t, products, 10)
	assert.Equal(t, "Product 11", products[0].Name)
	assert.Equal(t, "Product 20", products[9].Name)

	products, err = productDB.GetAll(3, 10, "asc")
	assert.NoError(t, err)
	assert.Len(t, products, 3)
	assert.Equal(t, "Product 21", products[0].Name)
	assert.Equal(t, "Product 23", products[2].Name)

}

func TestGetByID(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}

	db.Migrator().DropTable(&entity.Product{})
	db.AutoMigrate(&entity.Product{})

	product, err := entity.NewProduct("Test Product 1", 10.0)
	assert.NoError(t, err)
	db.Create(product)

	productDB := NewProduct(db)
	productFound, err := productDB.GetByID(product.ID.String())
	assert.NoError(t, err)
	assert.Equal(t, "Test Product 1", productFound.Name)
}

func TestUpdate(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}

	db.Migrator().DropTable(&entity.Product{})
	db.AutoMigrate(&entity.Product{})

	product, err := entity.NewProduct("Test Product 2", 10.0)
	assert.NoError(t, err)
	db.Create(product)

	productDB := NewProduct(db)
	productFound, err := productDB.GetByID(product.ID.String())
	assert.NoError(t, err)
	productFound.Name = "Updated Product 2"
	err = productDB.Update(productFound.ID.String(), productFound)
	assert.NoError(t, err)

	productUpdated, err := productDB.GetByID(product.ID.String())
	assert.NoError(t, err)
	assert.Equal(t, "Updated Product 2", productUpdated.Name)
}

func TestDelete(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}

	db.Migrator().DropTable(&entity.Product{})
	db.AutoMigrate(&entity.Product{})

	product, err := entity.NewProduct("Test Product 3", 10.0)
	assert.NoError(t, err)
	db.Create(product)
	productDB := NewProduct(db)

	err = productDB.Delete(product.ID.String())
	assert.NoError(t, err)
	productFound, err := productDB.GetByID(product.ID.String())
	assert.Error(t, err)
	assert.Nil(t, productFound)
}
