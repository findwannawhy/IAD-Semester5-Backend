package api

import (
	"log"

	"github.com/findwannawhy/IAD-Semester5/internal/app/handler"
	"github.com/findwannawhy/IAD-Semester5/internal/app/repository"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func StartServer() {
	log.Println("Starting server")

	repo, err := repository.NewRepository()

	if err != nil {
		logrus.Error("ошибка инициализации репозитория")
	}

	handler := handler.NewHandler(repo)

	r := gin.Default()
	// загружаем все html-шаблоны
	r.LoadHTMLGlob("templates/*")
	// префикс для всей статики
	r.Static("/static", "./resources")
	// префикс для изображений
	// r.Static("/img", "./resources/img")

	r.GET("/materials", handler.GetMaterials)
	r.GET("/material/:id", handler.GetMaterial)
	r.GET("/order/:id", handler.GetOrder)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	log.Println("Server down")
}