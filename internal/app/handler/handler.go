package handler

import (
	"errors"

	"github.com/findwannawhy/IAD-Semester5/internal/app/repository"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

func (h *Handler) RegisterHandler(router *gin.Engine) {
	router.GET("/materials", h.GetMaterials)
	router.GET("/material/:id", h.GetMaterial)
	router.POST("/material/create-material", h.CreateMaterial)
	router.DELETE("/material/:id/soft-delete", h.DeleteMaterial)
	router.PUT("/material/:id/change-material", h.ChangeMaterial)
	router.POST("/material/:id/add-to-experiment", h.AddMaterialToExperiment)
	router.POST("/material/:id/create-image", h.UploadImage)

	router.GET("/experiment/experiment-cart", h.GetExperimentCart)	
	router.GET("/experiments", h.GetExperiments)
	router.GET("/experiment/:id", h.GetExperiment)
	router.PUT("/experiment/:id/change-experiment", h.ChangeExperiment)
	router.PUT("/experiment/:id/form", h.FormExperiment)
	router.PUT("/experiment/:id/finish", h.ModerateExperiment)
	router.POST("/experiment/:id/soft-delete", h.DeleteExperiment)

	router.DELETE("/experiment_item/:material_id/:experiment_id", h.DeleteItemFromExperiment)
	router.PUT("/experiment_item/:material_id/:experiment_id", h.UpdateExperimentItem)

	router.POST("/user/sign-up", h.CreateUser)
	router.GET("/user/profile", h.GetProfile)
	router.PUT("/user/profile", h.ChangeProfile)
	router.POST("/user/sign-in", h.SignIn)
	router.POST("/user/sign-out", h.SignOut)
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