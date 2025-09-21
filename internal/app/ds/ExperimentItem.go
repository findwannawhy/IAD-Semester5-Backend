package ds

type ExperimentItem struct {
	ID                       uint       `gorm:"primaryKey"                                        json:"id"`

	ExperimentID             uint       `gorm:"not null;index"                                    json:"experiment_id"`
	MaterialID               uint       `gorm:"not null;index"                                    json:"material_id"`

	MaterialMass             float64    `gorm:"type:double precision;default:0.0"                                       json:"material_mass"`
	GasVolume                float64    `gorm:"type:double precision;default:0.0"                                       json:"gas_volume"`
	MassFractionPercentage  *float64    `gorm:"type:double precision"                                                  json:"mass_fraction_percentage"`

	Experiment              *Experiment `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"    json:"experiment"`
	Material                *Material   `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"    json:"material"`
}
