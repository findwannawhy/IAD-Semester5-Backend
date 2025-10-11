package ds

import (
	"time"

	"github.com/google/uuid"
)

type ImpurityFractionExperiment struct {
  ID           uint         `gorm:"primaryKey"                                       json:"id"`
	MolarVolume *float64      `gorm:"type:double precision"                            json:"molar_volume"`
	Status       string       `gorm:"type:varchar(16);not null;default:'draft';index;check:status IN ('draft','deleted','formed','finished','rejected')" json:"status"`
	CreatedAt    time.Time    `gorm:"<-:create;not null;index"                         json:"created_at"`            
	FormedAt    *time.Time    `gorm:""                                                 json:"formed_at"`
	FinishedAt  *time.Time    `gorm:""                                                 json:"finished_at"`

	CreatorID    uuid.UUID         `gorm:"not null;index;uniqueIndex:uid_one_draft_per_user,where:status = 'draft'" json:"creator_id"`
	ModeratorID *uuid.UUID         `gorm:"index"                                                                    json:"moderator_id"`

	Creator     *User         `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"   json:"creator"`
	Moderator   *User         `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"   json:"moderator"`

}