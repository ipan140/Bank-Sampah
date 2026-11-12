package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bank-sampah/core-go/app/modules/auth"
	"bank-sampah/core-go/app/modules/csr"
	"bank-sampah/core-go/app/modules/pickup"
	"bank-sampah/core-go/app/modules/reward"
	"bank-sampah/core-go/app/modules/tracking"
	"bank-sampah/core-go/app/modules/user"
	"bank-sampah/core-go/app/modules/wallet"
	"bank-sampah/core-go/app/modules/weighbridge"
	"bank-sampah/core-go/common/middleware"
	"bank-sampah/core-go/config"
	"bank-sampah/core-go/database/migrations"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	log.Println("==================================================")
	log.Println("   BANK SAMPAH MODULAR MONOLITH SERVER (CORE-GO)   ")
	log.Println("==================================================")

	// 1. Load Configuration
	cfg := config.LoadConfig()

	// 2. Initialize Database (Postgres)
	db, err := config.InitDB(cfg)
	if err != nil {
		log.Printf("⚠ Warning: DB connection failed: %v. Running in offline/mock mode if needed.", err)
	} else {
		// Run migrations
		if err := migrations.AutoMigrate(db); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
	}

	// 3. Initialize Redis
	rdb, err := config.InitRedis(cfg)
	if err != nil {
		log.Printf("⚠ Warning: Redis connection failed: %v", err)
	}

	// 4. Initialize Echo Router
	e := echo.New()
	e.HideBanner = true

	// Standard Global Middlewares
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))
	e.Use(middleware.AuditLogger())

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{
			"status":    "healthy",
			"timestamp": time.Now().Format(time.RFC3339),
			"version":   "1.0.0",
		})
	})

	apiV1 := e.Group("/api/v1")

	// 5. Dependency Injection & Module Wiring
	if db != nil {
		// Module: Auth
		authRepo := auth.NewRepository(db)
		authCacheRepo := auth.NewCacheRepository(rdb)
		authService := auth.NewService(authRepo, authCacheRepo, cfg)
		authHandler := auth.NewHandler(authService)
		auth.RegisterRoutes(apiV1, authHandler, cfg)

		// Module: User
		userRepo := user.NewRepository(db)
		userService := user.NewService(userRepo)
		userHandler := user.NewHandler(userService)
		user.RegisterRoutes(apiV1, userHandler, cfg)

		// Module: Wallet
		walletRepo := wallet.NewRepository(db)
		walletService := wallet.NewService(walletRepo)
		walletHandler := wallet.NewHandler(walletService)
		wallet.RegisterRoutes(apiV1, walletHandler, cfg)

		// Module: Pickup
		pickupRepo := pickup.NewRepository(db)
		pickupService := pickup.NewService(pickupRepo)
		pickupHandler := pickup.NewHandler(pickupService)
		pickup.RegisterRoutes(apiV1, pickupHandler, cfg)

		// Module: Weighbridge
		weighRepo := weighbridge.NewRepository(db)
		weighService := weighbridge.NewService(weighRepo)
		// Wire cross-module event: Weighbridge ticket settled -> Credit wallet
		weighService.SetSettlementCallback(func(ctx context.Context, nasabahID uuid.UUID, amount float64, desc string, refID uuid.UUID) error {
			_, err := walletService.CreditBalance(ctx, nasabahID, amount, desc, refID)
			return err
		})
		weighHandler := weighbridge.NewHandler(weighService)
		weighbridge.RegisterRoutes(apiV1, weighHandler, cfg)

		// Module: Tracking
		trackRepo := tracking.NewRepository(db)
		trackRedis := tracking.NewRedisRepository(rdb)
		trackService := tracking.NewService(trackRepo, trackRedis)
		trackHandler := tracking.NewHandler(trackService)
		tracking.RegisterRoutes(apiV1, trackHandler, cfg)

		// Module: Reward
		rewardRepo := reward.NewRepository(db)
		rewardService := reward.NewService(rewardRepo)
		rewardHandler := reward.NewHandler(rewardService)
		reward.RegisterRoutes(apiV1, rewardHandler, cfg)

		// Module: CSR
		csrRepo := csr.NewRepository(db)
		csrService := csr.NewService(csrRepo)
		csrHandler := csr.NewHandler(csrService)
		csr.RegisterRoutes(apiV1, csrHandler, cfg)
	}

	// 6. Graceful Server Start
	address := fmt.Sprintf(":%s", cfg.AppPort)
	go func() {
		log.Printf("🚀 Server running on port %s", cfg.AppPort)
		if err := e.Start(address); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server stopped abruptly: %v", err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("✓ Server exited properly")
}
