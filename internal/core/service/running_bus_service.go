package service

import (
	"bus-timing/internal/aggregate"
	"bus-timing/internal/entity"
	"bus-timing/pkg/common"
	"bus-timing/pkg/location"
	"bus-timing/pkg/uwave"
	"context"
	"fmt"
	"time"
)

type RunningBusService struct {
	UWaveClient interface {
		GetBusLines(ctx context.Context) (uwave.GetBusLineResponse, error)
		GetRunningBusByBusLineID(ctx context.Context, busLineID string) (uwave.GetRunningBusResponse, error)
	}
}

func (service *RunningBusService) EstimatedArrivalTime(ctx context.Context, busStopID string) ([]aggregate.IncomingBus, error) {
	resp, err := service.UWaveClient.GetBusLines(ctx)
	if err != nil {
		return nil, err
	}

	// get bus line data, return if no data
	busLinesBusStops := toBusLinesBusStopAggregate(resp)
	if busLinesBusStops == nil {
		return nil, nil
	}

	// find bus line pass bus stop, return if no bus line existed
	busLines := findBusLineByBusStopID(busLinesBusStops, busStopID)
	if busLines == nil {
		return nil, nil
	}

	busStopInfo := getBusStopInfo(busLinesBusStops, busStopID)
	if busStopInfo == nil {
		return nil, fmt.Errorf("cannot find bus stop with ID: %s", busStopID)
	}
	incomingBus := []aggregate.IncomingBus{}
	for _, busLine := range busLines {
		resp, err := service.UWaveClient.GetRunningBusByBusLineID(ctx, busLine.ID)
		if err != nil {
			return nil, err
		}

		runningBusPositions := toRunningBusPositionEntity(resp)
		if len(runningBusPositions) == 0 {
			continue
		}

		busStopLocation := location.Location{
			Lat: busStopInfo.Lat,
			Lng: busStopInfo.Lng,
		}
		busStopPath := findCurrentPath(busLine, busStopLocation)
		if len(busStopPath) == 0 {
			continue
		}

		// find nearest bus to bus stop
		nearestBus := findNearestBusToBusStop(busLine.BusLinePaths[0], runningBusPositions, *busStopInfo)
		if nearestBus == nil {
			continue
		}
		nearestBusLocation := location.Location{
			Lat: nearestBus.RunningBusPosition.Lat,
			Lng: nearestBus.RunningBusPosition.Lng,
		}
		busCurrentPath := findCurrentPath(busLine, nearestBusLocation)
		if len(busCurrentPath) == 0 {
			continue
		}

		distance := calculateDistanceFromBusToBusStop(busLine, busCurrentPath, nearestBusLocation, busStopPath, busStopLocation)
		incomingBus = append(incomingBus, aggregate.IncomingBus{
			Bus:         nearestBus.Bus,
			BusLine:     busLine,
			BusPosition: nearestBus.RunningBusPosition,
			Distance:    distance,
			ArrivalTime: time.Duration(distance / common.MapCrowdLevelAndSpeed[nearestBus.RunningBusPosition.CrowdLevel]),
		})
	}

	return incomingBus, nil
}

func getBusStopInfo(busLines []aggregate.BusLineBusStop, busStopID string) *entity.BusStop {
	for _, v := range busLines {
		for _, busStop := range v.BusStops {
			if busStop.ID == busStopID {
				return &busStop
			}
		}
	}
	return nil
}

func findBusLineByBusStopID(busLinesBusStops []aggregate.BusLineBusStop, busStopID string) []entity.BusLine {
	if len(busLinesBusStops) == 0 {
		return nil
	}

	busLines := make([]entity.BusLine, 0)
	for _, val := range busLinesBusStops {
		for _, busStop := range val.BusStops {
			if busStop.ID == busStopID {
				busLines = append(busLines, val.BusLine)
				break
			}
		}
	}
	return busLines
}

func findNearestBusToBusStop(firstBusLinePath entity.BusLinePath, runningBusPositions []aggregate.BusPosition, busStop entity.BusStop) *aggregate.BusPosition {
	if len(runningBusPositions) == 0 {
		return nil
	}

	firstBusLineLocation := location.Location{
		Lat: firstBusLinePath.Lat,
		Lng: firstBusLinePath.Lng,
	}
	busStopLocation := location.Location{
		Lat: busStop.Lat,
		Lng: busStop.Lng,
	}

	// distance from bus to first path position
	minDistance := location.CalculateStraightLine(firstBusLineLocation,
		location.Location{Lat: runningBusPositions[0].RunningBusPosition.Lat, Lng: runningBusPositions[0].RunningBusPosition.Lng})
	minBusIdx := 0
	for i, busPosition := range runningBusPositions {
		busLocation := location.Location{
			Lat: busPosition.RunningBusPosition.Lat,
			Lng: busPosition.RunningBusPosition.Lng,
		}
		// check bus position is between start of bus line and bus stop
		if !location.IsPointBetween(firstBusLineLocation, busStopLocation, busLocation) {
			continue
		}
		distance := location.CalculateStraightLine(firstBusLineLocation, busLocation)
		if distance < minDistance {
			minDistance = distance
			minBusIdx = i
		}
	}

	return &runningBusPositions[minBusIdx]
}

func findCurrentPath(busLine entity.BusLine, position location.Location) []int {
	pathPositionNearest := findNearestPathPositionIndex(busLine, position)
	switch {
	case pathPositionNearest == 0:
		return []int{0, 1}
	case pathPositionNearest == len(busLine.BusLinePaths)-1:
		return []int{len(busLine.BusLinePaths) - 2, len(busLine.BusLinePaths) - 1}
	case pathPositionNearest > 0 && pathPositionNearest < len(busLine.BusLinePaths)-1:
		prePosition := busLine.BusLinePaths[pathPositionNearest-1]
		nearestPosition := busLine.BusLinePaths[pathPositionNearest]
		if location.IsPointBetween(location.Location{Lat: prePosition.Lat, Lng: prePosition.Lng}, location.Location{Lat: nearestPosition.Lat, Lng: nearestPosition.Lng}, position) {
			return []int{pathPositionNearest - 1, pathPositionNearest}
		}
		nextPosition := busLine.BusLinePaths[pathPositionNearest+1]
		if location.IsPointBetween(location.Location{Lat: nearestPosition.Lat, Lng: nearestPosition.Lng}, location.Location{Lat: nextPosition.Lat, Lng: nextPosition.Lng}, position) {
			return []int{pathPositionNearest, pathPositionNearest + 1}
		}
	}

	return nil
}

func findNearestPathPositionIndex(busLine entity.BusLine, position location.Location) int {
	if len(busLine.BusLinePaths) == 0 {
		return -1
	}
	minDistance := location.CalculateStraightLine(location.Location{Lat: busLine.BusLinePaths[0].Lat, Lng: busLine.BusLinePaths[0].Lng}, position)
	minIdx := 0
	for idx, v := range busLine.BusLinePaths {
		distance := location.CalculateStraightLine(location.Location{Lat: v.Lat, Lng: v.Lng}, position)
		if distance < minDistance {
			minDistance = distance
			minIdx = idx
		}
	}

	return minIdx
}

func calculateDistanceFromBusToBusStop(busLine entity.BusLine, busCurrentPath []int, busLocation location.Location, busStopPath []int, busStopLocation location.Location) float64 {
	nearestPositionWithBus := busLine.BusLinePaths[busCurrentPath[1]]
	distanceFromBusToNearestPathPosition := location.CalculateStraightLine(location.Location{Lat: nearestPositionWithBus.Lat, Lng: nearestPositionWithBus.Lng}, busLocation)

	nearestPrePositionWithBusStop := busLine.BusLinePaths[busStopPath[0]]
	distanceFromPathPositionToBusStop := location.CalculateStraightLine(location.Location{Lat: nearestPrePositionWithBusStop.Lat, Lng: nearestPrePositionWithBusStop.Lng}, busStopLocation)

	// calculate distance between each path
	distanceInPath := 0.0
	for i := busCurrentPath[1]; i <= busStopPath[0]-1; i++ {
		l1 := busLine.BusLinePaths[i]
		l2 := busLine.BusLinePaths[i+1]
		distanceInPath += location.CalculateStraightLine(location.Location{Lat: l1.Lat, Lng: l1.Lng}, location.Location{Lat: l2.Lat, Lng: l2.Lng})
	}

	// distance from bus to bus stop = distance from bus to nearest path position + distance between each path + distance of nearest path position with bus stop
	return distanceFromBusToNearestPathPosition + distanceFromPathPositionToBusStop + distanceInPath
}
