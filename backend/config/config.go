package config

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

var (
	Port          = 0 //Porta onde a API roda
	DbName        = ""
	DbUser        = ""
	DbPassword    = ""
	DbPort        = ""
	SecretKey     = ""
	CertFile      = ""
	KeyFile       = ""
	Production    = false
	EmailPassword = ""
)

func Load() {
	var err error = nil

	if err := godotenv.Load(); err != nil {
		log.Fatal(err.Error()) //API PARA DE RODAR
	}

	if Port, err = strconv.Atoi(os.Getenv("PORT")); err != nil {
		log.Fatal(err.Error())
	}

	DbName = os.Getenv("DB_NAME")
	DbUser = os.Getenv("DB_USER")
	DbPassword = os.Getenv("DB_PASSWORD")
	DbPort = os.Getenv("DB_PORT")
	SecretKey = os.Getenv("SECRET_KEY")

	CertFile = os.Getenv("CERT_FILE")
	KeyFile = os.Getenv("KEY_FILE")
	if Production, err = strconv.ParseBool(os.Getenv("PRODUCTION")); err != nil {
		log.Fatal("Erro ao converter a chave PRODUCTION:", err.Error())
	}

	EmailPassword = os.Getenv("EMAIL_PASSWORD")
}

func GetContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}
