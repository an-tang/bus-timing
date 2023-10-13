package entity

import "bus-timing/pkg/common"

type RunningBusPosition struct {
	ID         string
	Lat        float64
	Lng        float64
	CrowdLevel common.CrowdLevel
}
