package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/findwannawhy/IAD-Semester5/internal/app/ds"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var errNoDraft = errors.New("no draft found")

func (r *Repository) GetExperiments(from, to time.Time, status string) ([]ds.Experiment, error) {
	var experiments []ds.Experiment
	sub := r.db.Where("status != 'deleted' and status != 'draft'")
	if !from.IsZero() {
		sub = sub.Where("created_at > ?", from)
	}
	if !to.IsZero() {
		sub = sub.Where("created_at < ?", to.Add(time.Hour*24))
	}

	if status != "" {
		sub = sub.Where("status = ?", status)
	}

	err := sub.Order("id").Find(&experiments).Error
	if err != nil {
		return nil, err
	}
	return experiments, nil
}

func (r *Repository) GetExperimentItems(experimentId uint) ([]ds.ExperimentItem, error) {
	var experimentItems []ds.ExperimentItem
	err := r.db.Where("experiment_id = ?", experimentId).Find(&experimentItems).Error
	if err != nil {
		return nil, err
	}
	return experimentItems, nil
}

func (r *Repository) GetExperimentItem(MaterialID uint, ExperimentID uint) (ds.ExperimentItem, error) {
	var experimentItem ds.ExperimentItem
	err := r.db.Where("material_id = ? and experiment_id = ?", MaterialID, ExperimentID).First(&experimentItem).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.ExperimentItem{}, fmt.Errorf("%w: experiment item not found", ErrNotFound)
		}
		return ds.ExperimentItem{}, err
	}
	return experimentItem, nil
}

func (r *Repository) GetExperimentMaterials(id uint) ([]ds.Material, ds.Experiment, error) {
	experiment, err := r.GetSingleExperiment(id)
	if err != nil {
		return []ds.Material{}, ds.Experiment{}, err
	}

	var materials []ds.Material
	sub := r.db.Table("experiment_items").Where("experiment_id = ?", experiment.ID)
	err = r.db.Order("id DESC").Where("id IN (?)", sub.Select("material_id")).Find(&materials).Error

	if err != nil {
		return []ds.Material{}, ds.Experiment{}, err
	}

	return materials, experiment, nil
}

func (r *Repository) CheckCurrentExperimentDraft(creatorID uint) (ds.Experiment, error) {	
	var experiment ds.Experiment
	res := r.db.Where("creator_id = ? AND status = ?", creatorID, "draft").Limit(1).Find(&experiment)
	if res.Error != nil {
		return ds.Experiment{}, res.Error
	} else if res.RowsAffected == 0 {
		return ds.Experiment{}, ErrNoDraft
	}
	return experiment, nil
}

func (r *Repository) GetExperimentDraft(creatorID uint) (ds.Experiment, bool, error) {
	experiment, err := r.CheckCurrentExperimentDraft(creatorID)
	if errors.Is(err, ErrNoDraft) {
		experiment = ds.Experiment{
			Status:     "draft",
			CreatorID:  creatorID,
			CreatedAt: time.Now(),
		}
		result := r.db.Create(&experiment)
		if result.Error != nil {
			return ds.Experiment{}, false, result.Error
		}
		return experiment, true, nil
	} else if err != nil {
		return ds.Experiment{}, false, err
	}
	return experiment, true, nil
}

func (r *Repository) GetExperimentCount(creatorID uint) int64 {
	if creatorID == 0 {
			return 0
	}
		
	var count int64
	experiment, err := r.CheckCurrentExperimentDraft(creatorID)
	if err != nil {
		return 0
	}
	err = r.db.Model(&ds.ExperimentItem{}).Where("experiment_id = ?", experiment.ID).Count(&count).Error
	if err != nil {
		logrus.Println("Error counting records in experiment_items:", err)
	}

	return count
}

func (r *Repository) GetSingleExperiment(id uint) (ds.Experiment, error) {
	var experiment ds.Experiment
	err := r.db.Where("id = ?", id).First(&experiment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.Experiment{}, fmt.Errorf("%w: эксперимент с id %d", ErrNotFound, id)
		}
		return ds.Experiment{}, err
	} else if experiment.Status == "deleted"  {
		return ds.Experiment{}, fmt.Errorf("%w: эксперимент удален", ErrNotAllowed)
	}
	return experiment, nil
}

func (r *Repository) FormExperiment(experimentId uint, status string) (ds.Experiment, error) {
	experiment, err := r.GetSingleExperiment(experimentId)
	if err != nil {
		return ds.Experiment{}, err
	}

	if experiment.Status != "draft" {
		return ds.Experiment{}, fmt.Errorf("эта заявка не может быть %s", status)
	}
	
	if status != "deleted"{
		if experiment.MolarVolume <= 0 {
		 	return ds.Experiment{}, errors.New("некорректный молярный объем")
		}
		ExperimentItems, _ := r.GetExperimentItems(experiment.ID)
		for _, experimentItem := range ExperimentItems {
			if experimentItem.MaterialMass <= 0{
				return ds.Experiment{}, errors.New("некорректная масса материала" )			
			}
			if experimentItem.GasVolume <= 0{
				return ds.Experiment{}, errors.New("некорректный объем газа" )			
			}
		}
	}	

	now := time.Now()
	err = r.db.Model(&experiment).Updates(ds.Experiment{
		Status:   ds.ExperimentStatus(status),
		FormedAt: &now,
	}).Error
	if err != nil {
		return ds.Experiment{}, err
	}

	return experiment, nil
}

func (r *Repository) UpdateExperiment(id uint, experiment ds.Experiment) (ds.Experiment, error) {
	dbExperiment := ds.Experiment{}

	if experiment.MolarVolume <= 0 {
		return ds.Experiment{}, errors.New("некорректный молярный объем")
  }

	err := r.db.Where("id = ? and status != 'deleted'", id).First(&dbExperiment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.Experiment{}, fmt.Errorf("%w: эксперимент с id %d", ErrNotFound, id)
		}
		return ds.Experiment{}, err
	}
	err = r.db.Model(&dbExperiment).Updates(experiment).Error
	if err != nil {
		return ds.Experiment{}, err
	}
	return dbExperiment, nil
}

func CalculateMassFraction(
	molarVolume float64,
	relativeMolecularMass float64,
	stoichiometricCoefficient float64,
	materialMass float64,
	gasVolume float64,
) (float64, error) {

	if molarVolume == 0 {
		return 0, errors.New("неправильный молярный объем")
	}

	if materialMass == 0 {
		return 0, errors.New("неправильная масса материала")
	}
	
	if gasVolume == 0 {
		return 0, errors.New("неправильный объем газа")
	}

	// Найдём количество вещества газа в молях
	nGas := gasVolume / molarVolume
	// Найдём количество вещества чистого вещества в молях
	nPureSubstance := nGas / stoichiometricCoefficient
	// Найдём массу чистого вещества в граммах
	mPureSubstance := nPureSubstance * relativeMolecularMass
	// Найдём массу примесей в граммах
	mImpurities := materialMass - mPureSubstance
	// Найдём массовую долю примесей в процентах
	massFractionPercentage := mImpurities / materialMass * 100
	return massFractionPercentage, nil
}

func (r *Repository) ModerateExperiment(id uint, status ds.ExperimentStatus) (ds.Experiment, error) {
	if status != ds.StatusFinished && status != ds.StatusRejected {
		return ds.Experiment{}, errors.New("неверный статус")
	}

	user, err := r.GetUserByID(r.GetUserID())
	if err != nil {
		return ds.Experiment{}, err
	}

	if !user.IsModerator {
		return ds.Experiment{}, fmt.Errorf("%w: вы не модератор", ErrNotAllowed)
	}

	experiment, err := r.GetSingleExperiment(id)
	if err != nil {
		return ds.Experiment{}, err
	} else if experiment.Status != "formed" {
		return ds.Experiment{}, fmt.Errorf("это исследование не может быть %s", status)
	}

	if status == ds.StatusFinished {
		experimentItems, err := r.GetExperimentItems(experiment.ID)
		if err != nil {
			return ds.Experiment{}, err
		}
		for _, experimentItem := range experimentItems {
			material, err := r.GetMaterial(experimentItem.MaterialID)
			if err != nil {
				return ds.Experiment{}, err
			}
			massFractionPercentage, err := CalculateMassFraction(experiment.MolarVolume, material.RelativeMolecularMass, material.StoichiometricCoefficient, experimentItem.MaterialMass, experimentItem.GasVolume)
			if err != nil {
				return ds.Experiment{}, err
			}
			err = r.db.Model(&experimentItem).Updates(ds.ExperimentItem{
				MassFractionPercentage: &massFractionPercentage,
			}).Error
			if err != nil {
				return ds.Experiment{}, err
			}
		}
	}

	now := time.Now()
	err = r.db.Model(&experiment).Updates(ds.Experiment{
		Status: status,
		FinishedAt: &now,
		ModeratorID: &user.ID,
	}).Error
	if err != nil {
		return ds.Experiment{}, err
	}

	return experiment, nil
}