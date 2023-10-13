package aggregate

import "bus-timing/internal/entity"

type BusLineBusStop struct {
	BusLine  entity.BusLine
	BusStops []entity.BusStop
}
