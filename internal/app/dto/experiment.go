package dto

import (
	"time"

	"github.com/findwannawhy/IAD-Semester5/internal/app/ds"
)

type ExperimentSampleCard struct {
	SampleID                   uint
	Title                      string
	Formula                    string
	ImageURL                   string
	RelativeMolecularMass      float64
	StoichiometricCoefficient  float64
	SampleMass                 string
	EvolvedGasVolume           string
	MassFractionPercentage     string
}

type DraftExperimentResponse struct {
	ExperimentID uint `json:"experiment_id"`
	SampleCount int64   `json:"sample_count"`
}

type ExperimentsResponse struct {
	ID uint                `json:"id"`
	MolarVolume float64    `json:"molar_volume"`
	Status string          `json:"status"`
	CreatedAt time.Time    `json:"created_at"`
	FormedAt *time.Time    `json:"formed_at"`
	FinishedAt *time.Time  `json:"finished_at"`
	CreatorLogin   string  `json:"creator_login"`
	ModeratorLogin string  `json:"moderator_login"`
}

type ExperimentResponse struct {
	Experiment             ds.ImpurityFractionExperiment `json:"experiment"`
	ExperimentSampleCards  []ExperimentSampleCard `json:"experiment_sample_cards"`
	CreatorLogin           string `json:"creator_login"`
	ModeratorLogin         string `json:"moderator_login"`
}

type UpdateExperiment struct {
	Experiment ds.ImpurityFractionExperiment `json:"experiment"`
	CreatorLogin string      `json:"creator_login"`
}

type FormExperiment struct {
	Experiment   ds.ImpurityFractionExperiment `json:"experiment"`
	CreatorLogin string        `json:"creator_login"`
}

type SoftDeleteExperiment struct {
	ExperimentID   uint `json:"id"`
	Status string `json:"status"`
	FormedAt *time.Time `json:"deleted_at"`
	CreatorLogin   string `json:"creator_login"`
}

type ModerateExperiment struct {
	Experiment ds.ImpurityFractionExperiment `json:"experiment"`
	CreatorLogin string      `json:"creator_login"`
	ModeratorLogin string    `json:"moderator_login"`
}