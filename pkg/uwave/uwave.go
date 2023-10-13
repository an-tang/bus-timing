package uwave

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

type UWaveClient struct {
	Endpoint string
}

type GetBusLineRequest struct {
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

type GetRunningBusRequest struct {
}

type GetRunningBusResponse struct {
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

type RunningBusPort struct {
	RunningBusService interface {
		GetRunningBus(ctx context.Context, busLineID string)
	}
}

func (u *UWaveClient) GetBusLines(ctx context.Context) (GetBusLineResponse, error) {
	requestURL := fmt.Sprintf("%s/busLines", u.Endpoint)
	res, err := http.Get(requestURL)
	if err != nil {
		return GetBusLineResponse{}, errors.Wrap(err, "UWaveClient.GetBusLines")
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
	}

	resp := GetBusLineResponse{}
	if err := json.Unmarshal(resBody, &resp); err != nil {
		return GetBusLineResponse{}, errors.Wrap(err, "UWaveClient.GetBusLines")
	}

	// plan, _ := os.ReadFile("./test_data/bus_line.json")
	// resp := GetBusLineResponse{}
	// err := json.Unmarshal(plan, &resp)
	// if err != nil {
	// 	return GetBusLineResponse{}, err
	// }

	return resp, nil
}

func (u *UWaveClient) GetRunningBusByBusLineID(ctx context.Context, busLineID string) (GetRunningBusResponse, error) {
	requestURL := fmt.Sprintf("%s/busPositions/%s", u.Endpoint, busLineID)
	res, err := http.Get(requestURL)
	if err != nil {
		return GetRunningBusResponse{}, errors.Wrap(err, "UWaveClient.GetRunningBusByBusLineID")
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		return GetRunningBusResponse{}, errors.Wrap(err, "UWaveClient.GetRunningBusByBusLineID")
	}

	resp := GetRunningBusResponse{}
	if err := json.Unmarshal(resBody, &resp); err != nil {
		return GetRunningBusResponse{}, errors.Wrap(err, "UWaveClient.GetRunningBusByBusLineID")
	}

	// plan, _ := os.ReadFile(fmt.Sprintf("./test_data/bus_line_position_%s.json", busLineID))
	// resp := GetRunningBusResponse{}
	// err := json.Unmarshal(plan, &resp)
	// if err != nil {
	// 	return GetRunningBusResponse{}, err
	// }
	return resp, nil
}
