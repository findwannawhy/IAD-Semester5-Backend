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
	// загружаем все html-шаблоны, чтобы потом их можно было использовать
	r.LoadHTMLGlob("templates/*")
	// делаем статические файлы доступными по URL‑префиксу /static
	// слева название папки, в которую выгрузится наша статика
	// справа путь к папке, в которой лежит статика
	r.Static("/static", "./resources")

	r.GET("/materials", handler.GetMaterials)
	r.GET("/material/:id", handler.GetMaterial)
	r.GET("/order", handler.GetOrder)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	log.Println("Server down")
}