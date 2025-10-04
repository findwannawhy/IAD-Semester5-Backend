package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/findwannawhy/IAD-Semester5/internal/app/ds"
	"github.com/findwannawhy/IAD-Semester5/internal/app/dto"
	"github.com/findwannawhy/IAD-Semester5/internal/app/repository"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetSamples(ctx *gin.Context) {
	var samples []ds.AcidSolubleSample
	var err error

	searchQuery := ctx.Query("search_sample")
	if searchQuery == "" {
		samples, err = h.Repository.GetSamples()
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
	} else {
		samples, err = h.Repository.GetSamplesByName(searchQuery)
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
	}
	if samples == nil {
		samples = make([]ds.AcidSolubleSample, 0)
	}
	ctx.JSON(http.StatusOK, samples)
}

func (h *Handler) GetSample(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	id := uint(id64)

	sample, err := h.Repository.GetSample(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, sample)
}

func (h *Handler) CreateSample(ctx *gin.Context) {
	var sampleJSON ds.AcidSolubleSample
	if err := ctx.BindJSON(&sampleJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	sample, err := h.Repository.CreateSample(sampleJSON)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Header("Location", fmt.Sprintf("/samples/%v", sample.ID))
	ctx.JSON(http.StatusCreated, sample)
}

func (h *Handler) SoftDeleteSample(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	id := uint(id64)

	err = h.Repository.SoftDeleteSample(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"sample_id": id,
		"status": "deleted",
	})
}

func (h *Handler) UpdateSample(ctx *gin.Context) {
	var sampleJSON ds.AcidSolubleSample
	if err := ctx.BindJSON(&sampleJSON); err != nil {
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

	sample, err := h.Repository.UpdateSample(id, sampleJSON)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, sample)
}

func (h *Handler) AddSampleToExperimentDraft(ctx *gin.Context) {
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

	sampleId64, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	sampleId := uint(sampleId64)

	err = h.Repository.AddSampleToExperimentDraft(uint(experimentId), sampleId)
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

	creatorLogin, _,err := h.Repository.GetModeratorAndCreatorLogin(experiment)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(status, dto.AddSampleToExperiment{
		SampleID: sampleId,
		ExperimentId: experiment.ID,
		ExperimentCreatedAt: experiment.CreatedAt,
		CreatorLogin: creatorLogin,
	})
}

func (h *Handler) UpdateImage(ctx *gin.Context) {
	sampleId64, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	sampleId := uint(sampleId64)

	file, err := ctx.FormFile("image")
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	sample, err := h.Repository.UpdateImage(ctx, sampleId, file)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, dto.UploadImage{
		SampleID: sample.ID,
		SampleTitle: sample.Title,
		ImageURL: *sample.ImageURL,
	})
}

