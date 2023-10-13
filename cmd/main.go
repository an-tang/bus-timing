package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	config "bus-timing/configuration"
	"bus-timing/internal/core/port"
	"bus-timing/internal/core/service"
	"bus-timing/pkg/middlewares/cors"
	"bus-timing/pkg/uwave"

	"github.com/gin-gonic/gin"
)

func RunServer() {
	router := SetupHTTP()
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", config.Config.Server.Host, config.Config.Server.Port),
		WriteTimeout: time.Second * time.Duration(config.Config.Server.WriteTimeout),
		ReadTimeout:  time.Second * time.Duration(config.Config.Server.ReadTimeout),
		IdleTimeout:  time.Second * time.Duration(config.Config.Server.IdleTimeout),
		Handler:      router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln("Server listen: ", err)
		}
	}()

	quit := make(chan os.Signal, 2)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutdown server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalln("server shutdown:", err)
	}

	if _, ok := <-ctx.Done(); ok {
		log.Println("timeout of 1 second.")
	}

	log.Println("server exiting")
}

func SetupHTTP() *gin.Engine {
	router := gin.Default()

	uWaveClient := uwave.UWaveClient{
		Endpoint: config.Config.UWaveConfig.Endpoint,
	}
	busLineService := service.BusLiveService{
		UWaveClient: &uWaveClient,
	}
	busPositionService := service.BusPositionService{
		UWaveClient: &uWaveClient,
	}
	runningBusService := service.RunningBusService{
		UWaveClient: &uWaveClient,
	}
	busLinePort := port.BusLinePort{
		BusLineService: &busLineService,
	}
	busPositionPort := port.BusPositionPort{
		BusPositionService: &busPositionService,
	}
	runningBusPort := port.RunningBusPort{
		BusTimingService: &runningBusService,
	}

	router.Use(gin.Recovery())
	router.Use(cors.CorsMiddleware())

	// health check
	router.GET("/health", port.HealthCheck)

	routerGroup := router.Group("api")
	// routerGroup.Use(jwt.Authorized())

	routerGroup.GET("/busPosition/:busLineID", busPositionPort.GetBusPosition)
	routerGroup.GET("/busLines", busLinePort.GetBusLines)
	routerGroup.GET("/busStop/:busStopID", runningBusPort.EstimatedArrival)

	return router
}
