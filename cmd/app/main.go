package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/OlegGibadulin/proxy-server/internal/models"
	reqHndlr "github.com/OlegGibadulin/proxy-server/internal/request/delivery"
	reqRp "github.com/OlegGibadulin/proxy-server/internal/request/repository"
	reqUcs "github.com/OlegGibadulin/proxy-server/internal/request/usecases"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
}

func getDBConnection() *gorm.DB {
	user, _ := os.LookupEnv("DB_USER")
	password, _ := os.LookupEnv("DB_PWD")
	host, _ := os.LookupEnv("DB_HOST")
	port, _ := os.LookupEnv("DB_PORT")
	dbname, _ := os.LookupEnv("DB_NAME")

	conn_data := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		user, password, host, port, dbname)

	dbConnection, err := gorm.Open("postgres", conn_data)
	if err != nil {
		log.Fatal(err)
	}
	return dbConnection
}

func main() {
	dbConnection := getDBConnection()
	dbConnection.AutoMigrate(&models.Request{})
	defer dbConnection.Close()

	requestRepo := reqRp.NewRequestPgRepository(dbConnection)
	requestUseCase := reqUcs.NewRequestUseCase(requestRepo)
	requestHandler := reqHndlr.NewRequestHandler(requestUseCase)

	e := echo.New()
	e.Use(middleware.Logger())

	requestHandler.Configure(e)

	proxyPort, _ := os.LookupEnv("PROXY_PORT")
	if proxyPort == "" {
		logrus.Info("Default proxy port 8080 was used")
		proxyPort = "8080"
	}
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", proxyPort),
		Handler: http.HandlerFunc(requestHandler.HandleRequest),
	}

	go func() {
		log.Fatal(server.ListenAndServe())
	}()

	serverPort, _ := os.LookupEnv("SERVER_PORT")
	if serverPort == "" {
		logrus.Info("Default server port 8000 was used")
		serverPort = "8000"
	}
	log.Fatal(e.Start(fmt.Sprintf(":%s", serverPort)))
}
