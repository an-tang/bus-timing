package port

import (
	"context"
	"net/http"

	"bus-timing/internal/aggregate"

	"github.com/gin-gonic/gin"
)

var (
	statusSuccess = 100000
)

type BusLinePort struct {
	BusLineService interface {
		GetBusLines(ctx context.Context) ([]aggregate.BusLineBusStop, error)
	}
}

type GetBusLinesRequest struct {
}

type GetBusLineResponse struct {
	Payload []BusLinePayload `json:"payload"`
	Status  int              `json:"status"`
}

type BusLinePayload struct {
	BusStops  []BusStop   `json:"busStops"`
	FullName  string      `json:"fullName"`
	ID        string      `json:"id"`
	Origin    string      `json:"origin"`
	Path      [][]float64 `json:"path"`
	ShortName string      `json:"shortName"`
}

type BusStop struct {
	ID   string  `json:"id"`
	Lat  float64 `json:"lat"`
	Lng  float64 `json:"lng"`
	Name string  `json:"name"`
}

func (port *BusLinePort) GetBusLines(ctx *gin.Context) {
	busLines, err := port.BusLineService.GetBusLines(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := transformBusLinesResponse(busLines)

	ctx.JSON(http.StatusOK, resp)
}

func transformBusLinesResponse(busLineBusStops []aggregate.BusLineBusStop) GetBusLineResponse {
	busLinePayloads := make([]BusLinePayload, 0)
	for _, val := range busLineBusStops {
		busStops := make([]BusStop, 0, len(val.BusStops))
		for _, busStop := range val.BusStops {
			busStops = append(busStops, BusStop{
				ID:   busStop.ID,
				Name: busStop.Name,
				Lat:  busStop.Lat,
				Lng:  busStop.Lng,
			})
		}
		paths := make([][]float64, 0, len(val.BusLine.BusLinePaths))
		for _, path := range val.BusLine.BusLinePaths {
			paths = append(paths, []float64{path.Lat, path.Lng})
		}

		busLine := BusLinePayload{
			BusStops:  busStops,
			FullName:  val.BusLine.FullName,
			ShortName: val.BusLine.ShortName,
			Origin:    val.BusLine.Origin,
			ID:        val.BusLine.ID,
			Path:      paths,
		}

		busLinePayloads = append(busLinePayloads, busLine)
	}

	return GetBusLineResponse{
		Payload: busLinePayloads,
		Status:  statusSuccess,
	}
}
