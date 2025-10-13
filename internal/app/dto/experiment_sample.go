package dto

type DeleteExperimentSampleResp struct {
	Message string `json:"message"`
	RemovedSampleID int `json:"removed_sample_id"`
	ExperimentID int `json:"experiment_id"`
}

type UpdateExperimentSampleReq struct {
	SampleMass       *float64 `json:"sample_mass"`
	EvolvedGasVolume *float64 `json:"evolved_gas_volume"`
}