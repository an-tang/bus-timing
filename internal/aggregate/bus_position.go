package aggregate

import "bus-timing/internal/entity"

type BusPosition struct {
	Bus                entity.Bus
	RunningBus         entity.RunningBus
	RunningBusPosition entity.RunningBusPosition
}
