package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/findwannawhy/IAD-Semester5/internal/app/ds"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetSample(ctx *gin.Context) {
	idStr := ctx.Param("id") 
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		logrus.Error(err)
		return
	}

	sample, err := h.Repository.GetSample(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		logrus.Error(err)
		return
	}

	ctx.HTML(http.StatusOK, "sample.html", gin.H{
		"sample": sample,
	})
}

func (h *Handler) GetSamples(ctx *gin.Context) {
	var samples []ds.AcidSolubleSample
	var err error
	creatorID := h.Repository.GetUser()

	searchSample := ctx.Query("search_sample")

	searchSample = strings.TrimSpace(searchSample)

	if searchSample == "" {
		samples, err = h.Repository.GetSamples()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			logrus.Error(err)
			return
		}
	} else {
		samples, err = h.Repository.SearchSamples(searchSample)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			logrus.Error(err)
			return
		}
	}
	currentExperiment, _ := h.Repository.CheckCurrentExperimentDraft(creatorID)

	ctx.HTML(http.StatusOK, "samples.html", gin.H{
		"samples":       samples,
		"count":           h.Repository.GetExperimentSamplesCount(),
		"search_sample": searchSample,
		"experimentId":    int(currentExperiment.ID),
	})
}

func (h *Handler) AddSampleToExperiment(ctx *gin.Context) {
	experiment, err := h.Repository.GetExperimentDraft(h.Repository.GetUser())
	experimentId := experiment.ID
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	sampleId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.AddSampleToExperiment(experimentId, uint(sampleId))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Redirect(http.StatusFound, "/acid-soluble-samples")
}