package common

type (
	RunningBusStatus string

	CrowdLevel string
)

const (
	HighCrowd   CrowdLevel = "high"
	MediumCrowd CrowdLevel = "medium"
	LowCrowd    CrowdLevel = "low"
)

var MapCrowdLevelAndSpeed = map[CrowdLevel]float64{
	HighCrowd:   40.0,
	MediumCrowd: 50.0,
	LowCrowd:    60.0,
}
