package service

import (
	"context"

	"bus-timing/internal/aggregate"
	"bus-timing/internal/entity"
	"bus-timing/pkg/uwave"
)

type BusLiveService struct {
	UWaveClient interface {
		GetBusLines(ctx context.Context) (uwave.GetBusLineResponse, error)
	}
}

func (service *BusLiveService) GetBusLines(ctx context.Context) ([]aggregate.BusLineBusStop, error) {
	resp, err := service.UWaveClient.GetBusLines(ctx)
	if err != nil {
		return nil, err
	}

	busLinesBusStops := toBusLinesBusStopAggregate(resp)
	return busLinesBusStops, nil
}

func toBusLinesBusStopAggregate(object uwave.GetBusLineResponse) []aggregate.BusLineBusStop {
	if len(object.Payload) == 0 {
		return nil
	}

	busLineBusStops := make([]aggregate.BusLineBusStop, 0, len(object.Payload))
	for _, val := range object.Payload {
		busStops := make([]entity.BusStop, 0, len(val.BusStops))
		for _, busStop := range val.BusStops {
			busStops = append(busStops, entity.BusStop{
				ID:   busStop.ID,
				Name: busStop.Name,
				Lat:  busStop.Lat,
				Lng:  busStop.Lng,
			})
		}

		busLinePaths := make([]entity.BusLinePath, 0, len(val.Path))
		for _, path := range val.Path {
			busLinePaths = append(busLinePaths, entity.BusLinePath{
				Lat: path[0],
				Lng: path[1],
			})
		}
		busLine := entity.BusLine{
			ID:           val.ID,
			FullName:     val.FullName,
			ShortName:    val.ShortName,
			Origin:       val.Origin,
			BusLinePaths: busLinePaths,
		}

		busLineBusStops = append(busLineBusStops, aggregate.BusLineBusStop{
			BusLine:  busLine,
			BusStops: busStops,
		})
	}

	return busLineBusStops
}
