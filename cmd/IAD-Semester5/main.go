// RegisterHandler godoc
// @title IAD-Semester5 API
// @version 1.0
// @description API для управления экспериментами по определению массовой доли примесей в образце
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
package main

import (
	"fmt"

	_ "github.com/findwannawhy/IAD-Semester5/docs"
	"github.com/findwannawhy/IAD-Semester5/internal/app/config"
	"github.com/findwannawhy/IAD-Semester5/internal/app/dsn"
	"github.com/findwannawhy/IAD-Semester5/internal/app/handler"
	"github.com/findwannawhy/IAD-Semester5/internal/app/repository"
	"github.com/findwannawhy/IAD-Semester5/internal/pkg"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	router := gin.Default()
	
	// Настройка CORS для Tauri и браузера
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://localhost:3000", "http://localhost:3000", "https://tauri.localhost", "http://tauri.localhost"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			// Разрешаем все tauri:// и localhost запросы
			return origin == "tauri://localhost" || 
				   origin == "http://localhost:3000" || 
				   origin == "https://localhost:3000" ||
				   origin == "http://tauri.localhost" ||
				   origin == "https://tauri.localhost"
		},
	}))
	conf, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}

	postgresString := dsn.FromEnv()
	fmt.Println(postgresString)

	rep, errRep := repository.NewRepository(postgresString)
	if errRep != nil {
		logrus.Fatalf("error initializing repository: %v", errRep)
	}

	hand := handler.NewHandler(rep)

	application := pkg.NewApp(conf, router, hand)
	application.RunApp()
}