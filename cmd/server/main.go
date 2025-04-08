package main

import (
	"apis/db"
	"apis/internal/entity"
	"apis/internal/infra/database"
	"apis/internal/infra/webserver/handlers"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := db.Connect()
	if err != nil {
		panic("Erro ao conectar ao banco")
	}
	// Close está embutido no *sql.DB
	defer db.Close()

	// Conecta o GORM usando a conexão existente
	gormDB, err := gorm.Open(sqlite.Dialector{Conn: db}, &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("erro ao conectar o GORM: %v", err))
	}

	err = gormDB.AutoMigrate(&entity.Product{}, entity.User{})
	if err != nil {
		panic(fmt.Sprintf("erro ao migrar: %v", err))
	}

	productDB := database.NewProduct(gormDB)
	productHandler := handlers.NewProductHandler(productDB)

	r := chi.NewRouter()
	// Injetando Logs
	r.Use(middleware.Logger)
	// Products
	r.Post("/products", productHandler.Create)
	r.Get("/products/{id}", productHandler.GetByID)
	r.Get("/products", productHandler.GetAll)
	r.Put("/products/{id}", productHandler.Update)
	r.Delete("/products/{id}", productHandler.Delete)

	http.ListenAndServe(":8080", r)
}
