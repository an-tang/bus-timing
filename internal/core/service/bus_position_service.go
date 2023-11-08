package service

import (
	"context"

	"bus-timing/internal/aggregate"
	"bus-timing/internal/entity"
	"bus-timing/pkg/common"
	"bus-timing/pkg/uwave"
)

type BusPositionService struct {
	UWaveClient interface {
		GetRunningBusByBusLineID(ctx context.Context, busLineID string) (uwave.GetRunningBusResponse, error)
	}
}

func (service *BusPositionService) GetBusPosition(ctx context.Context, busLineID string) ([]aggregate.BusPosition, error) {
	resp, err := service.UWaveClient.GetRunningBusByBusLineID(ctx, busLineID)
	if err != nil {
		return nil, err
	}

	runningBusPositions := toRunningBusPositionEntity(resp)
	return runningBusPositions, nil
}

func toRunningBusPositionEntity(object uwave.GetRunningBusResponse) []aggregate.BusPosition {
	if len(object.Payload) == 0 {
		return nil
	}

	runningBuses := make([]aggregate.BusPosition, 0, len(object.Payload))
	for _, val := range object.Payload {
		bus := entity.Bus{
			VehiclePlate: val.VehiclePlate,
			Bearing:      val.Bearing,
		}
		runningBus := entity.RunningBus{
			// Status: common.RunningBusStatus(val.CrowdLevel),
		}
		runningBusPosition := entity.RunningBusPosition{
			Lat:        val.Lat,
			Lng:        val.Lng,
			CrowdLevel: common.CrowdLevel(val.CrowdLevel),
		}

		runningBuses = append(runningBuses, aggregate.BusPosition{
			Bus:                bus,
			RunningBus:         runningBus,
			RunningBusPosition: runningBusPosition,
		})
	}

	return runningBuses
}
