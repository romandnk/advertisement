package main

import (
	"context"
	"github.com/romandnk/advertisement/configs"
	"github.com/romandnk/advertisement/internal/logger"
	"github.com/romandnk/advertisement/internal/server/http"
	"github.com/romandnk/advertisement/internal/service"
	"github.com/romandnk/advertisement/internal/storage/postgres"
	"go.uber.org/zap"
	"net"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	config, err := configs.NewConfig()
	if err != nil {
		panic("error initialising config: " + err.Error())
	}

	log, err := logger.NewZapLogger(config.ZapLogger.Level, config.ZapLogger.Encoding, config.ZapLogger.OutputPath, config.ZapLogger.ErrorOutputPath)
	if err != nil {
		panic("error initialising zap logger: " + err.Error())
	}

	log.Log.Info("using zap logger...")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	db, err := postgres.NewPostgresDB(ctx, config.Postgres)
	if err != nil {
		log.Error("error connecting postgres db", zap.String("error", err.Error()))
		return
	}
	defer db.Close()

	log.Log.Info("using postgres db",
		zap.String("address", net.JoinHostPort(config.Postgres.Host, strconv.Itoa(config.Postgres.Port))))

	storage := postgres.NewPostgresStorage(db)

	services := service.NewService(storage, log, config.PathToImages)

	handler := http.NewHandler(services, log)

	server := http.NewServer(config.Server.Host, config.Server.Port,
		config.Server.ReadTimeout, config.Server.WriteTimeout, handler.InitRoutes())

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			log.Error("error stopping server", zap.String("error", err.Error()))
			return
		}

		log.Info("app stopped")
	}()

	log.Info("app is starting...")

	if err := server.Start(); err != nil {
		log.Error("error starting server", zap.String("error", err.Error()))
		return
	}
}
