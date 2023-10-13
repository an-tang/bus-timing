package port

import (
	"bus-timing/internal/aggregate"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type RunningBusPort struct {
	BusTimingService interface {
		EstimatedArrivalTime(ctx context.Context, busStopID string) ([]aggregate.IncomingBus, error)
	}
}

type GetIncomingBusRequest struct {
}

type IncomingBusResponse struct {
	Payload []BusLine `json:"payload"`
	Status  int       `json:"status"`
}

type BusLine struct {
	ID        string `json:"busLineID"`
	FullName  string `json:"fullName"`
	ShortName string `json:"shortName"`
	Origin    string `json:"origin"`
	Bus       Bus    `json:"bus"`
}

type Bus struct {
	Lat          float64       `json:"lat"`
	Lng          float64       `json:"lng"`
	VehiclePlate string        `json:"vehiclePlate"`
	TimeDuration time.Duration `json:"timeDuration"`
	Distance     float64       `json:"distance"`
}

func (port *RunningBusPort) EstimatedArrival(ctx *gin.Context) {
	busStopID := ctx.Param("busStopID")
	if busStopID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid bus stop: %s", busStopID)})
		return
	}

	incomingBuses, err := port.BusTimingService.EstimatedArrivalTime(ctx, busStopID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, transformIncomingBusToEstimatedArrival(incomingBuses))
}

func transformIncomingBusToEstimatedArrival(incomingBuses []aggregate.IncomingBus) IncomingBusResponse {
	payload := make([]BusLine, 0, len(incomingBuses))
	for _, val := range incomingBuses {
		payload = append(payload, BusLine{
			ID:        val.BusLine.ID,
			FullName:  val.BusLine.FullName,
			ShortName: val.BusLine.ShortName,
			Origin:    val.BusLine.Origin,
			Bus: Bus{
				Lat:          val.BusPosition.Lat,
				Lng:          val.BusPosition.Lng,
				VehiclePlate: val.Bus.VehiclePlate,
				TimeDuration: val.ArrivalTime,
				Distance:     val.Distance,
			},
		})
	}
	return IncomingBusResponse{
		Payload: payload,
		Status:  statusSuccess,
	}
}
