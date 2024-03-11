package cmd

import (
	"applicationDesignTest/config"
	"applicationDesignTest/internal/repository"
	"applicationDesignTest/internal/service"
	"applicationDesignTest/internal/transport/rest"
	"applicationDesignTest/internal/transport/rest/controllers"
	"applicationDesignTest/pkg/log"
	"context"
	"os"
	"os/signal"
	"syscall"
)

func InitApp(ctx context.Context) {
	cfg := config.NewConfig()
	logger := log.NewLogger()

	roomRepo := repository.NewRoom()
	orderRepo := repository.NewOrder()

	roomSrv := service.NewRoomService(roomRepo, logger)
	orderSrv := service.NewOrderService(orderRepo, roomSrv, logger)

	orderController := controllers.NewOrderController(orderSrv)
	server := rest.NewRestServer(cfg.GetAddress(), orderController, logger)

	go server.StartServer()
	defer server.StopServer(ctx)

	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
	case <-quit:
	}
}
