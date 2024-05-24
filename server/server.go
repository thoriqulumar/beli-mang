package server

import (
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

type Server struct {
	db        *sqlx.DB
	app       *echo.Echo
	validator *validator.Validate
	logger    *zap.Logger
}

func NewServer(db *sqlx.DB, logger *zap.Logger) *Server {
	app := echo.New()
	validate := validator.New()

	app.Use(middleware.Recover())
	app.Use(middleware.Logger())

	return &Server{
		db:        db,
		app:       app,
		validator: validate,
		logger:    logger,
	}
}

func (s *Server) Start() error {
	return s.app.Start(":8080")
}
