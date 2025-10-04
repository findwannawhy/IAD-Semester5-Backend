package repository

import (
	"github.com/findwannawhy/IAD-Semester5/internal/app/ds"
	"github.com/sirupsen/logrus"
)

// GetCartCount для получения количества услуг в заявке (чатов в сообщении в моем случае)
func (r *Repository) GetExperimentSamplesCount() int64 {
	var count int64
	creatorID := 1

	var draftExperiment ds.ImpurityFractionExperiment
	err := r.db.Where("creator_id = ? AND status = ?", creatorID, "draft").Select("id").First(&draftExperiment).Error
	if err != nil {
		return 0
	}

	err = r.db.Model(&ds.ExperimentSample{}).Where("experiment_id = ?", draftExperiment.ID).Count(&count).Error
	if err != nil {
		logrus.Println("Error counting records in lists_experiment_samples:", err)
	}

	return count
}