package application

import (
	"beli-mang/config"
	"beli-mang/database"
	"beli-mang/pkg/log"
	"beli-mang/server"
	"beli-mang/version"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Start(cfg *config.Config) {
	// init logger
	logger, err := log.New(zapcore.DebugLevel, version.ServiceID, version.Version)
	if err != nil {
		panic(err)
	}

	db, err := database.NewDatabase(cfg)
	if err != nil {
		logger.Error("error opening database", zap.Error(err))
		panic(err)
	}
	defer db.Close()

	s := server.NewServer(db, logger)
	s.RegisterRoute(cfg)

	logger.Fatal("failed run app", zap.Error(s.Start()))
}
