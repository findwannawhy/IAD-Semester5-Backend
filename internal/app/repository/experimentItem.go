package repository

import (
	"github.com/findwannawhy/IAD-Semester5/internal/app/ds"
	"github.com/sirupsen/logrus"
)

// GetCartCount для получения количества услуг в заявке (чатов в сообщении в моем случае)
func (r *Repository) GetExperimentItemsCount() int64 {
	var experimentID int
	var count int64
	creatorID := 1

	err := r.db.Model(&ds.Experiment{}).Where("creator_id = ? AND status = ?", creatorID, ds.StatusDraft).Select("id").First(&experimentID).Error
	if err != nil {
		 return 0
	}

	err = r.db.Model(&ds.ExperimentItem{}).Where("experiment_id = ?", experimentID).Count(&count).Error
	if err != nil {
		 logrus.Println("Error counting records in lists_experiment_items:", err)
	}

	return count
}