package handler

import (
	"errors"
	"net/http"

	"github.com/findwannawhy/IAD-Semester5/internal/app/repository"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

// RegisterHandler godoc
// @title IAD API
// @version 1.0
// @description API для управления экспериментами по определению массовой доли примесей в образце
// @contact.name API Support
// @contact.url http://localhost:8080
// @contact.email support@impurity.fraction.com
// @license.name MIT
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func (h *Handler) RegisterHandler(router *gin.Engine) {
	api := router.Group("/api/v1")

	unauthorized := api.Group("/")

	unauthorized.POST("/users/sign-up", h.SignUp)	
	unauthorized.GET("/api/soluble-samples", h.GetSamples) // список образцов
	unauthorized.GET("/api/soluble-samples/:id", h.GetSample) // получить образец по id
	unauthorized.POST("/users/sign-in", h.SignIn)

	authorized := api.Group("/")
	authorized.Use(h.ModeratorMiddleware(false))

	authorized.POST("/api/soluble-samples", h.CreateSample) // создать образец
	authorized.DELETE("/api/soluble-samples/:id", h.SoftDeleteSample) // удалить образец (soft)
	authorized.PUT("/api/soluble-samples/:id", h.UpdateSample) // обновить образец
	authorized.POST("/api/soluble-samples/:id/impurity-experiments/draft", h.AddSampleToExperimentDraft) // добавить образец в корзину
	authorized.POST("/api/soluble-samples/:id/image", h.UpdateImage) // обновить изображение образца

	authorized.GET("/api/impurity-experiments/draft", h.GetExperimentCart)	// текущий черновик пользователя
	authorized.GET("/api/impurity-experiments", h.GetExperiments) // список экспериментов
	authorized.GET("/api/impurity-experiments/:id", h.GetExperiment) // получить эксперимент по id
	authorized.PUT("/api/impurity-experiments/:id", h.UpdateExperiment) // обновить эксперимент
	authorized.PUT("/api/impurity-experiments/:id/formation", h.FormExperiment) // поменять статус эксперимента (user)
	authorized.DELETE("/api/impurity-experiments/:id", h.SoftDeleteExperiment) // удалить эксперимент (soft)

	authorized.DELETE("/api/impurity-experiments/:id/soluble-samples/:sample_id", h.DeleteSampleFromExperiment) // удалить образец из эксперимента
	authorized.PUT("/api/impurity-experiments/:id/soluble-samples/:sample_id", h.UpdateExperimentSample) // обновить позицию образца в эксперименте

	authorized.GET("/api/users/me", h.GetProfile) // профиль текущего пользователя
	authorized.PUT("/api/users/me", h.UpdateProfile) // обновить профиль
	authorized.POST("/api/users/logout", h.SignOut) // выход из системы

	moderator := api.Group("/")
	moderator.Use(h.ModeratorMiddleware(true))
	moderator.PUT("/api/impurity-experiments/:id/moderation", h.ModerateExperiment) // поменять статус эксперимента (moderator)

	swaggerURL := ginSwagger.URL("/swagger/doc.json")
	router.Any("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerURL))
	router.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
}

func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.Static("/static", "./resources")
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	
	var errorMessage string
	switch {
	case errors.Is(err, repository.ErrNotFound):
		errorMessage = "Не найден"
	case errors.Is(err, repository.ErrAlreadyExists):
		errorMessage = "Уже существует"
	case errors.Is(err, repository.ErrNotAllowed):
		errorMessage = "Доступ запрещен"
	case errors.Is(err, repository.ErrNoDraft):
		errorMessage = "Черновик не найден"
	default:
		errorMessage = err.Error()
	}
	
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": errorMessage,
	})
}