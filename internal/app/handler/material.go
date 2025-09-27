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
	resp := make([]ds.Material, 0, len(materials))
	for _, material := range materials {
		resp = append(resp, material)
	}
	ctx.JSON(http.StatusOK, resp)
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

func (h *Handler) DeleteMaterial(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	id := uint(id64)

	err = h.Repository.DeleteMaterial(id)
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

func (h *Handler) ChangeMaterial(ctx *gin.Context) {
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

	material, err := h.Repository.ChangeMaterial(id, materialJSON)
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
	experiment, created, err := h.Repository.GetExperimentDraft(h.Repository.GetUserID())
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

func (h *Handler) UploadImage(ctx *gin.Context) {
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

	material, err := h.Repository.UploadImage(ctx, materialId, file)
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