package repository

import (
	"errors"
	"fmt"

	"github.com/findwannawhy/IAD-Semester5/internal/app/ds"
)

var errNoDraft = errors.New("no draft found")

func (r *Repository) GetExperiment(id int) ([]ds.ItemCard, ds.Experiment, error) {

	creatorID := r.GetUser()

	var experiment ds.Experiment
	err := r.db.Where("id = ?", id).First(&experiment).Error
	if err != nil {
		return []ds.ItemCard{}, ds.Experiment{}, err
	} else if creatorID != int(experiment.CreatorID) {
		return []ds.ItemCard{}, ds.Experiment{}, errors.New("you are not allowed")
	} else if experiment.Status == ds.StatusDeleted {
		return []ds.ItemCard{}, ds.Experiment{}, errors.New("you can`t watch deleted calculations")
	}

	var experimentItems []ds.ExperimentItem
	var materials []ds.Material
	sub := r.db.Table("experiment_items").Where("experiment_id = ?", experiment.ID).Find(&experimentItems)
	err = r.db.Where("id IN (?)", sub.Select("material_id")).Find(&materials).Error
	if err != nil {
		return []ds.ItemCard{}, ds.Experiment{}, err
	}


	var cards []ds.ItemCard
	materialsMap := make(map[uint]ds.Material)

	for _, material := range materials {
		materialsMap[material.ID] = material
	}

	for _, item := range experimentItems {
		material, ok := materialsMap[item.MaterialID]
		if ok {
			materialMassStr := ""
			if item.MaterialMass != nil {
				materialMassStr = fmt.Sprintf("%.2f", *item.MaterialMass)
			}
			gasVolumeStr := ""
			if item.GasVolume != nil {
				gasVolumeStr = fmt.Sprintf("%.3f", *item.GasVolume)
			}
			massFractionPercentageStr := ""
			if item.MassFractionPercentage != nil {
				massFractionPercentageStr = fmt.Sprintf("%.1f", *item.MassFractionPercentage)
			}

			cards = append(cards, ds.ItemCard{
				ItemID:                    item.ID,
				Title:                     material.Title,
				Description:               material.Description,
				ImageURL:                  material.ImageURL,
				RelativeMolecularMass:     material.RelativeMolecularMass,
				StoichiometricCoefficient: material.StoichiometricCoefficient,
				MaterialMass:              materialMassStr,
				GasVolume:                 gasVolumeStr,
				MassFractionPercentage:    massFractionPercentageStr,
			})
		}
	}

	return cards, experiment, nil
}


func (r *Repository) CheckCurrentExperimentDraft(creatorID int) (ds.Experiment, error) {
	var experiment ds.Experiment

	res := r.db.Where("creator_id = ? AND status = ?", creatorID, "draft").Limit(1).Find(&experiment)
	if res.Error != nil {
		return ds.Experiment{}, res.Error
	} else if experiment.ID == 0 {
		return ds.Experiment{}, errNoDraft
	}
	return experiment, nil
}

func (r *Repository) GetExperimentDraft(creatorID int) (ds.Experiment, error) {
	experiment, err := r.CheckCurrentExperimentDraft(creatorID)
	if err == errNoDraft {
		experiment = ds.Experiment{
			CreatorID: uint(creatorID),
			Status: ds.StatusDraft,
		}
		res := r.db.Create(&experiment)
		if res.Error != nil {
			return ds.Experiment{}, res.Error
		}
		return experiment, nil
	} else if err != nil {
		return ds.Experiment{}, err
	}
	return experiment, nil
}

func (r *Repository) SoftDeleteExperiment(experimentId int) error{
	return r.db.Exec("UPDATE experiments SET status = 'deleted' WHERE id = ?", experimentId).Error
}