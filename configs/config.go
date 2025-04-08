package configs

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/go-chi/jwtauth/v5"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

// Conf representa as configurações do banco
type Conf struct {
	DBFile       string
	DBMode       string
	DBTimeout    string
	TokenAuth    *jwtauth.JWTAuth
	JwtExpiresIn int
}

// LoadConfig carrega as configurações do .env e retorna uma instância de Conf
func LoadConfig() (*Conf, error) {
	err := godotenv.Load("cmd/server/.env")
	if err != nil {
		return nil, fmt.Errorf("erro ao carregar .env: %v", err)
	}

	expString := os.Getenv("JWT_EXPIRATION")
	if expString == "" {
		expString = "3600" // valor padrão de 1 hora
	}
	expInt, err := strconv.Atoi(expString) // converte para inteiro
	if err != nil {
		log.Fatalf("JWT_EXPIRATION inválido: %v", err)
	}

	config := &Conf{
		DBFile:       os.Getenv("DB_FILE"),
		DBMode:       os.Getenv("DB_MODE"),
		DBTimeout:    os.Getenv("DB_TIMEOUT"),
		TokenAuth:    jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET")), nil),
		JwtExpiresIn: expInt, // Tempo de expiração do token em segundos
	}

	return config, nil
}
