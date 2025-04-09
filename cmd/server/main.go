package main

import (
	"apis/configs"
	"apis/db"
	"apis/internal/entity"
	"apis/internal/infra/database"
	"apis/internal/infra/webserver/handlers"
	"fmt"
	"net/http"

	_ "apis/docs"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email email@email.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
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
	// recovers from panics, logs the panic (and a backtrace), and returns a HTTP 500 (Internal Server Error) status if possible
	r.Use(middleware.Recoverer)
	r.Use(middleware.WithValue("jwt", cfg.TokenAuth))
	r.Use(middleware.WithValue("jwtExpiresIn", cfg.JwtExpiresIn))

	// Rotas públicas
	r.Post("/users", userHandler.Create)
	r.Post("/users/login", userHandler.GenerateJWT)
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

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
