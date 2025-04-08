package configs

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

var cfg *conf

func NewConfig() *conf {
	return cfg
}

// type conf struct {
// 	DBDriver      string `mapstructure:"DB_DRIVER"`
// 	DBHost        string `mapstructure:"DB_HOST"`
// 	DBPort        string `mapstructure:"DB_PORT"`
// 	DBUser        string `mapstructure:"DB_USER"`
// 	DBPassword    string `mapstructure:"DB_PASSWORD"`
// 	DBName        string `mapstructure:"DB_NAME"`
// 	WebServerPort string `mapstructure:"WEB_SERVER_PORT"`
// 	JWTSecret     string `mapstructure:"JWT_SECRET"`
// 	JwtExpiresIn  int    `mapstructure:"JWT_EXPIRES_IN"`
// 	TokenAuth     *jwtauth.JWTAuth
// }

// Conf representa as configurações do banco
type conf struct {
	DBFile    string
	DBMode    string
	DBTimeout string
}

// LoadConfig carrega as configurações do .env e retorna uma instância de Conf
func LoadConfig() (*conf, error) {
	err := godotenv.Load("cmd/server/.env")
	if err != nil {
		return nil, fmt.Errorf("erro ao carregar .env: %v", err)
	}

	config := &conf{
		DBFile:    os.Getenv("DB_FILE"),
		DBMode:    os.Getenv("DB_MODE"),
		DBTimeout: os.Getenv("DB_TIMEOUT"),
	}

	return config, nil
}
