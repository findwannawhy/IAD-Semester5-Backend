package ds

import (
	"time"
)

type ExperimentStatus string

const (
	StatusDraft    ExperimentStatus = "draft"     // черновик (создан черновик)
	StatusDeleted  ExperimentStatus = "deleted"   // удалён (создатель удалил)
	StatusFormed   ExperimentStatus = "formed"    // сформирован (создатель завершил оформление)
	StatusFinished ExperimentStatus = "finished"  // завершён (модератор подтвердил)
	StatusRejected ExperimentStatus = "rejected"  // отклонён (модератор отклонил)
)

type Experiment struct {
  ID           uint             `gorm:"primaryKey"                                       json:"id"`
	MolarVolume  float64          `gorm:"not null;default:22.4"                            json:"molar_volume"`
	Status       ExperimentStatus `gorm:"type:varchar(16);not null;default:'draft';index"  json:"status"`
	CreatedAt    time.Time        `gorm:"<-:create;not null;index"                         json:"created_at"`            
	FormedAt    *time.Time        `gorm:""                                                 json:"formed_at"`   // «дата формирования» (действие создателя)
	FinishedAt  *time.Time        `gorm:""                                                 json:"finished_at"` // «дата завершения» (действие модератора)

	CreatorID    uint             `gorm:"not null;index"                                   json:"creator_id"`
	ModeratorID *uint             `gorm:"index"                                            json:"moderator_id"`

	Creator     *User             `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"   json:"creator"`
	Moderator   *User             `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"   json:"moderator"`

	// правило: у каждого пользователя не более одного эксперимента в статусе черновик
	_ struct{} `gorm:"uniqueIndex:uid_one_draft_per_user,where:status = 'draft';"`
}