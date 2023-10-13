package aggregate

import (
	"time"

	"bus-timing/internal/entity"
)

type IncomingBus struct {
	Bus         entity.Bus
	BusLine     entity.BusLine
	BusPosition entity.RunningBusPosition
	Distance    float64
	ArrivalTime time.Duration
}
