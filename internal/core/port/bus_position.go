package port

import (
	"bus-timing/internal/aggregate"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetRunningBusRequest struct {
}

type GetBusPositionResponse struct {
	Payload []RunningBusPayload `json:"payload"`
	Status  int                 `json:"status"`
}

type RunningBusPayload struct {
	Bearing      float64 `json:"bearing"`
	CrowdLevel   string  `json:"crowdLevel"`
	Lat          float64 `json:"lat"`
	Lng          float64 `json:"lng"`
	VehiclePlate string  `json:"vehiclePlate"`
}

type BusPositionPort struct {
	BusPositionService interface {
		GetBusPosition(ctx context.Context, busLineID string) ([]aggregate.BusPosition, error)
	}
}

func (port *BusPositionPort) GetBusPosition(ctx *gin.Context) {
	busLineID := ctx.Param("busLineID")
	if busLineID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid bus line: %s", busLineID)})
		return
	}
	busLines, err := port.BusPositionService.GetBusPosition(ctx, busLineID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := transformBusPositionsResponse(busLines)

	ctx.JSON(http.StatusOK, resp)
}

func transformBusPositionsResponse(runningBuses []aggregate.BusPosition) GetBusPositionResponse {
	payload := make([]RunningBusPayload, 0, len(runningBuses))
	for _, val := range runningBuses {
		payload = append(payload, RunningBusPayload{
			Bearing:      val.Bus.Bearing,
			CrowdLevel:   string(val.RunningBusPosition.CrowdLevel),
			Lat:          val.RunningBusPosition.Lat,
			Lng:          val.RunningBusPosition.Lng,
			VehiclePlate: val.Bus.VehiclePlate,
		})
	}
	return GetBusPositionResponse{
		Payload: payload,
		Status:  statusSuccess,
	}
}
