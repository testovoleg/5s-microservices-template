package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/testovoleg/5s-microservice-template/api_gateway_service/config"
	"github.com/testovoleg/5s-microservice-template/api_gateway_service/internal/server"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
)

// @contact.name Oleg Testov
// @contact.url https://github.com/testovoleg
// @contact.email o.testov@5systems.ru
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Bearer token. Example: 'Bearer xxxxxxxxxxxxxxxx'
func main() {
	flag.Parse()

	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	appLogger := logger.NewAppLogger(cfg.Logger)
	appLogger.InitLogger()
	appLogger.WithName("ApiGateway")

	s := server.NewServer(appLogger, cfg)
	appLogger.Fatal(s.Run())
}

func init() {
	//init env
	if err := godotenv.Load(); err != nil {
		os.Chdir("../../")
		if err := godotenv.Load(); err != nil {
			fmt.Printf("No .env file found")
		}
	}
}
