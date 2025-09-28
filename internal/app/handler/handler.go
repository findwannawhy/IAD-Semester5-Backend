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
	router.DELETE("/material/:id/soft-delete", h.SoftDeleteMaterial)
	router.PUT("/material/:id/update-material", h.UpdateMaterial)
	router.POST("/material/:id/add-to-cart", h.AddMaterialToExperiment)
	router.POST("/material/:id/update-image", h.UpdateImage)

	router.GET("/experiment/cart", h.GetExperimentCart)	
	router.GET("/experiments", h.GetExperiments)
	router.GET("/experiment/:id", h.GetExperiment)
	router.PUT("/experiment/:id/update-experiment", h.UpdateExperiment)
	router.PUT("/experiment/:id/form", h.FormExperiment)
	router.PUT("/experiment/:id/moderate", h.ModerateExperiment)
	router.DELETE("/experiment/:id/soft-delete", h.SoftDeleteExperiment)

	router.DELETE("/experiment_item/:material_id/:experiment_id", h.DeleteItemFromExperiment)
	router.PUT("/experiment_item/:material_id/:experiment_id", h.UpdateExperimentItem)

	router.POST("/user/sign-up", h.CreateUser)
	router.GET("/user/profile", h.GetProfile)
	router.PUT("/user/profile", h.UpdateProfile)
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