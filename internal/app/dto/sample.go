package dto

import (
	"time"
)

type AddSampleToExperiment struct {
	SampleID                 uint       `json:"sample_id"`
	ExperimentId             uint       `json:"experiment_id"`
	ExperimentCreatedAt      time.Time  `json:"created_at"`
	CreatorLogin             string     `json:"creator_login"`
}

type UploadImage struct {
	SampleID     uint       `json:"id"`
	SampleTitle  string     `json:"title"`
	ImageURL       string     `json:"image_url"`
}