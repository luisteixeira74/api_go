package main

import (
	"apis/configs"
	"apis/db"
	"apis/internal/entity"
	"apis/internal/infra/database"
	"apis/internal/infra/webserver/handlers"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	cfg, err := configs.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("erro ao carregar as configs: %v", err))
	}

	sqlDB, err := db.Connect()
	if err != nil {
		panic("Erro ao conectar ao banco")
	}
	// Close está embutido no *sql.DB
	defer sqlDB.Close()

	// Conecta o GORM usando a conexão existente
	gormDB, err := gorm.Open(sqlite.Dialector{Conn: sqlDB}, &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("erro ao conectar o GORM: %v", err))
	}

	// Migração
	if err := gormDB.AutoMigrate(&entity.Product{}, &entity.User{}); err != nil {
		panic(fmt.Sprintf("erro ao migrar: %v", err))
	}

	// Handlers
	productHandler, userHandler := setupHandlers(gormDB, cfg)

	// Rotas
	r := setupRouter(cfg, productHandler, userHandler)

	fmt.Println("Servidor iniciado em :8080")
	http.ListenAndServe(":8080", r)
}

// inicializa os handlers com o banco de dados
func setupHandlers(db *gorm.DB, cfg *configs.Conf) (*handlers.ProductHandler, *handlers.UserHandler) {
	productDB := database.NewProduct(db)
	userDB := database.NewUser(db)
	// orderRepo := database.NewOrder(gormDB)  // camada de acesso ao banco

	productHandler := handlers.NewProductHandler(productDB)
	userHandler := handlers.NewUserHandler(userDB, cfg.TokenAuth, cfg.JwtExpiresIn)
	// orderHandler := handlers.NewOrderHandler(orderRepo)  // camada web

	return productHandler, userHandler
}

// configura as rotas do servidor
func setupRouter(cfg *configs.Conf, productHandler *handlers.ProductHandler, userHandler *handlers.UserHandler) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Rotas públicas
	r.Post("/users", userHandler.Create)
	r.Post("/users/login", userHandler.GenerateJWT)

	// Rotas protegidas
	registerProtectedRoutes(r, cfg, productHandler)

	return r
}

// registra as rotas protegidas
func registerProtectedRoutes(r chi.Router, cfg *configs.Conf, productHandler *handlers.ProductHandler) {
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(cfg.TokenAuth))
		r.Use(jwtauth.Authenticator(cfg.TokenAuth))

		r.Route("/products", func(r chi.Router) {
			r.Post("/", productHandler.Create)
			r.Get("/", productHandler.GetAll)
			r.Get("/{id}", productHandler.GetByID)
			r.Put("/{id}", productHandler.Update)
			r.Delete("/{id}", productHandler.Delete)
		})

	})
}
