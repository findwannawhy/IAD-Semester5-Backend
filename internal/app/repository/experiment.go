package repository

import (
	"errors"
	"fmt"
	"strings"

	"github.com/findwannawhy/IAD-Semester5/internal/app/ds"
)

var errNoDraft = errors.New("no draft found")

type ExperimentSampleCard struct {
	SampleID                   uint
	Title                      string
	Formula                    string
	ImageURL                  *string
	RelativeMolecularMass      float64
	StoichiometricCoefficient  float64
	SampleMass                 string
	EvolvedGasVolume           string
	MassFractionPercentage     string
}

func (r *Repository) GetExperiment(id int) ([]ExperimentSampleCard, ds.ImpurityFractionExperiment, error) {

	creatorID := r.GetUser()

	var experiment ds.ImpurityFractionExperiment
	err := r.db.Where("id = ?", id).First(&experiment).Error
	if err != nil {
		return []ExperimentSampleCard{}, ds.ImpurityFractionExperiment{}, err
	} else if creatorID != int(experiment.CreatorID) {
		return []ExperimentSampleCard{}, ds.ImpurityFractionExperiment{}, errors.New("you are not allowed")
	} else if experiment.Status == "deleted" {
		return []ExperimentSampleCard{}, ds.ImpurityFractionExperiment{}, errors.New("you can`t watch deleted experiment")
	}

	var experimentSamples []ds.ExperimentSample
	var samples []ds.AcidSolubleSample
	sub := r.db.Table("experiments_samples").Where("experiment_id = ?", experiment.ID).Find(&experimentSamples)
	err = r.db.Where("id IN (?)", sub.Select("sample_id")).Find(&samples).Error
	if err != nil {
		return []ExperimentSampleCard{}, ds.ImpurityFractionExperiment{}, err
	}

	var cards []ExperimentSampleCard
	samplesMap := make(map[uint]ds.AcidSolubleSample)

	for _, sample := range samples {
		samplesMap[sample.ID] = sample
	}

	for _, experimentSample := range experimentSamples {
		sample, ok := samplesMap[experimentSample.SampleID]
		if ok {
			massFractionPercentageStr := ""
			if experimentSample.MassFractionPercentage != nil {
				massFractionPercentageStr = fmt.Sprintf("%.2f", *experimentSample.MassFractionPercentage)
			}
			sampleMassStr := ""
			if experimentSample.SampleMass != nil {
				sampleMassStr = fmt.Sprintf("%g", *experimentSample.SampleMass)
				if !strings.Contains(sampleMassStr, ".") {
					sampleMassStr += ".0"
				}
			}

			gasVolumeStr := ""
			if experimentSample.EvolvedGasVolume != nil {
				gasVolumeStr = fmt.Sprintf("%g", *experimentSample.EvolvedGasVolume)
				if !strings.Contains(gasVolumeStr, ".") {
					gasVolumeStr += ".0"
				}
			}

			cards = append(cards, ExperimentSampleCard{
				SampleID:                  sample.ID,
				Title:                     sample.Title,
				Formula:                   sample.Formula,
				ImageURL:                  sample.ImageURL,
				RelativeMolecularMass:     sample.RelativeMolecularMass,
				StoichiometricCoefficient: sample.StoichiometricCoefficient,
				SampleMass:                sampleMassStr,
				EvolvedGasVolume:          gasVolumeStr,
				MassFractionPercentage:    massFractionPercentageStr,
			})
		}
	}

	return cards, experiment, nil
}


func (r *Repository) CheckCurrentExperimentDraft(creatorID int) (ds.ImpurityFractionExperiment, error) {
	var experiment ds.ImpurityFractionExperiment

	res := r.db.Where("creator_id = ? AND status = ?", creatorID, "draft").Limit(1).Find(&experiment)
	if res.Error != nil {
		return ds.ImpurityFractionExperiment{}, res.Error
	} else if experiment.ID == 0 {
		return ds.ImpurityFractionExperiment{}, errNoDraft
	}
	return experiment, nil
}

func (r *Repository) GetExperimentDraft(creatorID int) (ds.ImpurityFractionExperiment, error) {
	experiment, err := r.CheckCurrentExperimentDraft(creatorID)
	if err == errNoDraft {
		experiment = ds.ImpurityFractionExperiment{
			CreatorID: uint(creatorID),
			Status: "draft",
		}
		res := r.db.Create(&experiment)
		if res.Error != nil {
			return ds.ImpurityFractionExperiment{}, res.Error
		}
		return experiment, nil
	} else if err != nil {
		return ds.ImpurityFractionExperiment{}, err
	}
	return experiment, nil
}

func (r *Repository) SoftDeleteExperiment(experimentId int) error{
	return r.db.Exec("UPDATE impurity_fraction_experiments SET status = 'deleted' WHERE id = ?", experimentId).Error
}