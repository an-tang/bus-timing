package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"bus-timing/internal/aggregate"
	"bus-timing/internal/entity"
	"bus-timing/pkg/common"
	"bus-timing/pkg/location"
	"bus-timing/pkg/uwave"

	"github.com/stretchr/testify/assert"
)

type mockUWaveClient struct {
	getBusLines              func(ctx context.Context) (uwave.GetBusLineResponse, error)
	getRunningBusByBusLineID func(ctx context.Context, busLineID string) (uwave.GetRunningBusResponse, error)
}

func (m mockUWaveClient) GetBusLines(ctx context.Context) (uwave.GetBusLineResponse, error) {
	return m.getBusLines(ctx)
}

func (m mockUWaveClient) GetRunningBusByBusLineID(ctx context.Context, busLineID string) (uwave.GetRunningBusResponse, error) {
	return m.getRunningBusByBusLineID(ctx, busLineID)
}

func TestRunningBusService_EstimatedArrivalTime(t *testing.T) {
	t.Parallel()

	busLineData, _ := os.ReadFile("./../../../test_data/bus_line_less_data.json")
	busLine := uwave.GetBusLineResponse{}
	err := json.Unmarshal(busLineData, &busLine)
	assert.NoError(t, err)

	busLinePositionData, _ := os.ReadFile(fmt.Sprintf("./../../../test_data/bus_line_position_%s.json", "44480"))
	busLinePosition := uwave.GetRunningBusResponse{}
	err = json.Unmarshal(busLinePositionData, &busLinePosition)
	assert.NoError(t, err)

	t.Run("happy case", func(tt *testing.T) {
		busStopID := "377906"
		uwaveClient := mockUWaveClient{
			getBusLines: func(ctx context.Context) (uwave.GetBusLineResponse, error) {
				return busLine, nil
			},

			getRunningBusByBusLineID: func(ctx context.Context, busLineID string) (uwave.GetRunningBusResponse, error) {
				return busLinePosition, nil
			},
		}

		expected := []aggregate.IncomingBus{
			{
				Bus: entity.Bus{
					Bearing:      159.4,
					VehiclePlate: "PD1064Z",
				},
				BusPosition: entity.RunningBusPosition{
					Lat:        1.338066,
					Lng:        103.695944,
					CrowdLevel: common.LowCrowd,
				},
				BusLine: entity.BusLine{
					ID: "44481",
				},
				Distance:    168,
				ArrivalTime: 2,
			},
			{
				Bus: entity.Bus{
					Bearing:      159.4,
					VehiclePlate: "PD1064Z",
				},
				BusPosition: entity.RunningBusPosition{
					Lat:        1.338066,
					Lng:        103.695944,
					CrowdLevel: common.LowCrowd,
				},
				BusLine: entity.BusLine{
					ID: "44480",
				},
				Distance:    169,
				ArrivalTime: 2,
			},
		}

		svc := &RunningBusService{
			UWaveClient: uwaveClient,
		}

		resp, err := svc.EstimatedArrivalTime(context.Background(), busStopID)
		assert.NoError(tt, err)
		for i, val := range expected {
			assert.Equal(tt, val.ArrivalTime, resp[i].ArrivalTime)
			assert.Equal(tt, val.BusLine.ID, resp[i].BusLine.ID)
			assert.Equal(tt, val.Distance, resp[i].Distance)
			assert.Equal(tt, val.Bus.VehiclePlate, resp[i].Bus.VehiclePlate)
		}
	})

	t.Run("bad case: get data from uwave failed", func(tt *testing.T) {
		busStopID := "377906"
		uwaveClient := mockUWaveClient{
			getBusLines: func(ctx context.Context) (uwave.GetBusLineResponse, error) {
				return busLine, http.ErrServerClosed
			},
		}
		expectedError := http.ErrServerClosed

		svc := &RunningBusService{
			UWaveClient: uwaveClient,
		}

		resp, err := svc.EstimatedArrivalTime(context.Background(), busStopID)
		assert.Error(tt, err)
		assert.Equal(tt, expectedError, err)
		assert.Nil(tt, resp)
	})

	t.Run("cannot find bus stop", func(tt *testing.T) {
		busStopID := "-1"
		uwaveClient := mockUWaveClient{
			getBusLines: func(ctx context.Context) (uwave.GetBusLineResponse, error) {
				return busLine, nil
			},

			getRunningBusByBusLineID: func(ctx context.Context, busLineID string) (uwave.GetRunningBusResponse, error) {
				return busLinePosition, nil
			},
		}
		expectedError := fmt.Errorf("cannot find bus stop with ID: %s", busStopID)

		svc := &RunningBusService{
			UWaveClient: uwaveClient,
		}

		resp, err := svc.EstimatedArrivalTime(context.Background(), busStopID)
		assert.Error(tt, err)
		assert.Equal(tt, expectedError, err)
		assert.Nil(tt, resp)
	})

	t.Run("no running bus", func(tt *testing.T) {
		busStopID := "377906"
		uwaveClient := mockUWaveClient{
			getBusLines: func(ctx context.Context) (uwave.GetBusLineResponse, error) {
				return busLine, nil
			},

			getRunningBusByBusLineID: func(ctx context.Context, busLineID string) (uwave.GetRunningBusResponse, error) {
				return uwave.GetRunningBusResponse{}, nil
			},
		}
		expected := []aggregate.IncomingBus{}

		svc := &RunningBusService{
			UWaveClient: uwaveClient,
		}

		resp, err := svc.EstimatedArrivalTime(context.Background(), busStopID)
		assert.NoError(tt, err)
		assert.Equal(tt, expected, resp)
	})
}

func Test_findBusLineByBusStopID(t *testing.T) {
	t.Parallel()

	t.Run("happy case", func(tt *testing.T) {
		busLinesBusStops := mockBusLine()
		busStopID := "377906"

		expected := []entity.BusLine{
			{
				ID: "44481",
			},
			{
				ID: "44480",
			},
		}
		resp := findBusLineByBusStopID(busLinesBusStops, busStopID)
		for i, val := range expected {
			assert.Equal(tt, val.ID, resp[i].ID)
		}
	})
}

func Test_findNearestBusToBusStop(t *testing.T) {
	t.Parallel()

	t.Run("happy case", func(tt *testing.T) {
		busLinesBusStops := mockBusLine()
		busPositions := mockBusPosition()
		busStop := busLinesBusStops[0].BusStops[0]

		expected := aggregate.BusPosition{
			Bus: entity.Bus{
				VehiclePlate: "PD1064Z",
				Bearing:      159.4,
			},
			RunningBusPosition: entity.RunningBusPosition{
				Lat:        1.338066,
				Lng:        103.695944,
				CrowdLevel: "low",
			},
		}
		resp := findNearestBusToBusStop(busLinesBusStops[0].BusLine.BusLinePaths[0], busPositions, busStop)
		assert.Equal(tt, expected.Bus.VehiclePlate, resp.Bus.VehiclePlate)
		assert.Equal(tt, expected.RunningBusPosition.Lat, resp.RunningBusPosition.Lat)
		assert.Equal(tt, expected.RunningBusPosition.Lng, resp.RunningBusPosition.Lng)
	})

	t.Run("no running bus", func(tt *testing.T) {
		busLinesBusStops := mockBusLine()
		busPositions := []aggregate.BusPosition{}
		busStop := busLinesBusStops[0].BusStops[0]

		resp := findNearestBusToBusStop(busLinesBusStops[0].BusLine.BusLinePaths[0], busPositions, busStop)
		assert.Nil(tt, resp)
	})
}

func Test_findCurrentPath(t *testing.T) {
	t.Parallel()

	t.Run("happy case: position current path in journey", func(tt *testing.T) {
		busLinesBusStops := mockBusLine()
		position := location.Location{
			Lat: 1.33771,
			Lng: 103.69753,
		}

		expected := []int{1, 2}
		resp := findCurrentPath(busLinesBusStops[0].BusLine, position)
		assert.Equal(tt, expected[0], resp[0])
		assert.Equal(tt, expected[1], resp[1])
	})

	t.Run("happy case: position current path in the first path", func(tt *testing.T) {
		busLinesBusStops := mockBusLine()
		position := location.Location{
			Lat: 1.33771,
			Lng: 103.69736,
		}

		expected := []int{0, 1}
		resp := findCurrentPath(busLinesBusStops[0].BusLine, position)
		assert.Equal(tt, expected[0], resp[0])
		assert.Equal(tt, expected[1], resp[1])
	})

	t.Run("happy case: position current path in the last path", func(tt *testing.T) {
		busLinesBusStops := mockBusLine()
		position := location.Location{
			Lat: 1.33771,
			Lng: 103.69727,
		}

		expected := []int{2, 3}
		resp := findCurrentPath(busLinesBusStops[0].BusLine, position)
		assert.Equal(tt, expected[0], resp[0])
		assert.Equal(tt, expected[1], resp[1])
	})
}

func mockBusLine() []aggregate.BusLineBusStop {
	busLineData, _ := os.ReadFile("./../../../test_data/bus_line_less_data.json")
	busLine := uwave.GetBusLineResponse{}
	err := json.Unmarshal(busLineData, &busLine)
	if err != nil {
		log.Fatalln(err)
	}
	return toBusLinesBusStopAggregate(busLine)
}

func mockBusPosition() []aggregate.BusPosition {
	busLinePositionData, _ := os.ReadFile(fmt.Sprintf("./../../../test_data/bus_line_position_%s.json", "44480"))
	busLinePosition := uwave.GetRunningBusResponse{}
	err := json.Unmarshal(busLinePositionData, &busLinePosition)
	if err != nil {
		log.Fatalln(err)
	}

	return toRunningBusPositionEntity(busLinePosition)
}
