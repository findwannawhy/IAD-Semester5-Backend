package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/findwannawhy/IAD-Semester5/internal/app/ds"
	"github.com/findwannawhy/IAD-Semester5/internal/app/repository"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetMaterials(ctx *gin.Context) {
	var materials []ds.Material
	var err error

	searchQuery := ctx.Query("material_name")
	if searchQuery == "" {
		materials, err = h.Repository.GetMaterials()
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
	} else {
		materials, err = h.Repository.GetMaterialsByName(searchQuery)
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
	}
	if materials == nil {
		materials = make([]ds.Material, 0)
	}
	ctx.JSON(http.StatusOK, materials)
}

func (h *Handler) GetMaterial(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	id := uint(id64)

	material, err := h.Repository.GetMaterial(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, material)
}

func (h *Handler) CreateMaterial(ctx *gin.Context) {
	var materialJSON ds.Material
	if err := ctx.BindJSON(&materialJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	material, err := h.Repository.CreateMaterial(materialJSON)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Header("Location", fmt.Sprintf("/materials/%v", material.ID))
	ctx.JSON(http.StatusCreated, material)
}

func (h *Handler) SoftDeleteMaterial(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	id := uint(id64)

	err = h.Repository.SoftDeleteMaterial(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "deleted",
	})
}

func (h *Handler) UpdateMaterial(ctx *gin.Context) {
	var materialJSON ds.Material
	if err := ctx.BindJSON(&materialJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	idStr := ctx.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	id := uint(id64)

	material, err := h.Repository.UpdateMaterial(id, materialJSON)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, material)
}

func (h *Handler) AddMaterialToExperiment(ctx *gin.Context) {
	userID := h.Repository.GetUserID()
	if userID == 0 {
		h.errorHandler(ctx, http.StatusUnauthorized, errors.New("user not authenticated"))
		return
	}

	experiment, created, err := h.Repository.GetExperimentDraft(userID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	experimentId := experiment.ID

	materialId64, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	materialId := uint(materialId64)

	err = h.Repository.AddMaterialToExperiment(uint(experimentId), materialId)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else if errors.Is(err, repository.ErrAlreadyExists) {
			h.errorHandler(ctx, http.StatusConflict, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	status := http.StatusOK

	if created {
		ctx.Header("Location", fmt.Sprintf("/experiment/%v", experiment.ID))
		status = http.StatusCreated
	}

	creatorLogin, moderatorLogin, err := h.Repository.GetModeratorAndCreatorLogin(experiment)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(status, gin.H{
		"experiment":      experiment,
		"creator_login":    creatorLogin,
		"moderator_login":  moderatorLogin,
	})
}

func (h *Handler) UpdateImage(ctx *gin.Context) {
	materialId64, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	materialId := uint(materialId64)

	file, err := ctx.FormFile("image")
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	material, err := h.Repository.UpdateImage(ctx, materialId, file)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":   "uploaded",
		"material": material,
	})
}