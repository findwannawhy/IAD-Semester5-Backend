package ds

type ExperimentItem struct {
	ID                       int        `gorm:"primaryKey"                                        json:"id"`

	ExperimentID             int        `gorm:"not null;index"                                    json:"experiment_id"`
	MaterialID               int        `gorm:"not null;index"                                    json:"material_id"`

	MaterialMass             float64    `gorm:"default:0.0"                                       json:"material_mass"`
	GasVolume                float64    `gorm:"default:0.0"                                       json:"gas_volume"`
	MassFractionPercentage  *float64    `gorm:""                                                  json:"mass_fraction_percentage"`

	Experiment              *Experiment `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"    json:"experiment"`
	Material                *Material   `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"    json:"material"`
}
