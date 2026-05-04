package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	registry "github.com/retail/backend/cmd/app/register"
	"github.com/retail/backend/internal/platform/cache"
	"github.com/retail/backend/internal/platform/config"
	"github.com/retail/backend/internal/platform/database"
	"github.com/retail/backend/internal/platform/logger"
	"github.com/retail/backend/internal/platform/storage"
	"github.com/retail/backend/internal/platform/middleware"
	route "github.com/retail/backend/internal/routes"
	jwtpkg "github.com/retail/backend/pkg/jwt"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	logger.Init(cfg.App.Env)

	db, err := database.NewPostgres(cfg.Database.DSN)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	redisClient, err := cache.NewRedis(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to redis")
	}
	defer redisClient.Close()

	jwtMgr := jwtpkg.NewManager(cfg.JWT.Secret, cfg.JWT.ExpiryHours, cfg.JWT.RefreshDays, redisClient)

	minioClient, err := storage.NewMinIO(cfg.MinIO.Endpoint, cfg.MinIO.AccessKey, cfg.MinIO.SecretKey, cfg.MinIO.Bucket, cfg.MinIO.UseSSL)
	if err != nil {
		log.Warn().Err(err).Msg("failed to connect to MinIO, file uploads will not work")
	}

	repos := registry.NewRepositories(db)
	services := registry.NewServices(db, repos, jwtMgr, redisClient, cfg)
	handlers := registry.NewHandlers(services, minioClient)

	// Initialize Fiber App
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			log.Error().Err(err).Str("method", c.Method()).Str("path", c.Path()).Msg("Unhandled Error")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false, "error": "internal server error",
			})
		},
	})

	// Global Middleware
	app.Use(recover.New())
	app.Use(fiberLogger.New())
	app.Use(middleware.ErrorLogger())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders: "Content-Type,Authorization",
	}))

	// Register Routes
	route.RegisterRoutes(app, jwtMgr, handlers)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-quit
		log.Info().Msg("shutting down gracefully...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := app.ShutdownWithContext(ctx); err != nil {
			log.Error().Err(err).Msg("error during shutdown")
		}
	}()

	log.Info().Str("port", cfg.App.Port).Msg("retail API starting")
	if err := app.Listen(":" + cfg.App.Port); err != nil {
		log.Fatal().Err(err).Msg("server error")
	}
}
