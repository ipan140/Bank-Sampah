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
	"bank-sampah/core-go/app/modules/compliance"
	"bank-sampah/core-go/app/modules/csr"
	"bank-sampah/core-go/app/modules/drop_point"
	"bank-sampah/core-go/app/modules/finance_report"
	"bank-sampah/core-go/app/modules/gamification"
	"bank-sampah/core-go/app/modules/inventory"
	"bank-sampah/core-go/app/modules/marketplace"
	"bank-sampah/core-go/app/modules/member_card"
	"bank-sampah/core-go/app/modules/notification"
	"bank-sampah/core-go/app/modules/pickup"
	"bank-sampah/core-go/app/modules/ppob"
	"bank-sampah/core-go/app/modules/reward"
	"bank-sampah/core-go/app/modules/tracking"
	"bank-sampah/core-go/app/modules/user"
	"bank-sampah/core-go/app/modules/wallet"
	"bank-sampah/core-go/app/modules/weighbridge"
	"bank-sampah/core-go/common/middleware"
	"bank-sampah/core-go/config"
	"bank-sampah/core-go/database/migrations"
	"bank-sampah/core-go/database/seeds"
	_ "bank-sampah/core-go/docs"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Bank Sampah Core ERP API
// @version 1.0.0
// @description Modular Monolith Backend API with Clean Layered Architecture & Swagger Documentation
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
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
		// Run auto migrations for all modules
		if err := migrations.AutoMigrate(db); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}

		// Run enterprise database seeders
		if err := seeds.SeedAll(db); err != nil {
			log.Printf("⚠ Warning running seeders: %v", err)
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

	// Swagger UI documentation
	e.GET("/swagger/*", echoSwagger.WrapHandler)

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

		// Module: Inventory & B2B Offtaker Sales
		invRepo := inventory.NewRepository(db)
		invService := inventory.NewService(invRepo)
		invHandler := inventory.NewHandler(invService)
		inventory.RegisterRoutes(apiV1, invHandler, cfg)

		// Module: Pickup
		pickupRepo := pickup.NewRepository(db)
		pickupService := pickup.NewService(pickupRepo)
		pickupHandler := pickup.NewHandler(pickupService)
		pickup.RegisterRoutes(apiV1, pickupHandler, cfg)

		// Module: Notification
		notifRepo := notification.NewRepository(db)
		notifService := notification.NewService(notifRepo)
		notifHandler := notification.NewHandler(notifService)
		notification.RegisterRoutes(apiV1, notifHandler, cfg)

		// Module: Weighbridge (Timbangan Digital)
		weighRepo := weighbridge.NewRepository(db)
		weighService := weighbridge.NewService(weighRepo)
		// Wire cross-module event:
		// 1) Ticket Settled -> Auto-credit Wallet
		// 2) Ticket Settled -> Stock In Inventory Gudang Sampah
		// 3) Ticket Settled -> Send In-App Notification
		weighService.SetSettlementCallback(func(ctx context.Context, nasabahID uuid.UUID, amount float64, desc string, refID uuid.UUID) error {
			_, err := walletService.CreditBalance(ctx, nasabahID, amount, desc, refID)
			if err == nil {
				// Record stock in
				ticket, _ := weighRepo.FindByID(ctx, refID)
				if ticket != nil {
					_ = invService.RecordStockIn(ctx, ticket.Category, ticket.NetWeight, ticket.PricePerKg, ticket.ID)
				}
				// Notify nasabah
				_, _ = notifService.SendNotification(ctx, &notification.SendNotificationDTO{
					UserID:        nasabahID,
					Title:         "Saldo Sampah Masuk!",
					Message:       fmt.Sprintf("Setoran sampah Anda berhasil diverifikasi. Saldo bertambah Rp %.0f (%s)", amount, desc),
					Channel:       notification.ChannelInApp,
					ReferenceType: "WALLET_CREDIT",
				})
			}
			return err
		})
		weighHandler := weighbridge.NewHandler(weighService)
		weighbridge.RegisterRoutes(apiV1, weighHandler, cfg)

		// Module: Member Card & Digital QR Code
		cardRepo := member_card.NewRepository(db)
		cardService := member_card.NewService(cardRepo)
		cardHandler := member_card.NewHandler(cardService)
		member_card.RegisterRoutes(apiV1, cardHandler, cfg)

		// Module: Tracking (Redis Geospatial)
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

		// Module: CSR & Environmental Impact
		csrRepo := csr.NewRepository(db)
		csrService := csr.NewService(csrRepo)
		csrHandler := csr.NewHandler(csrService)
		csr.RegisterRoutes(apiV1, csrHandler, cfg)

		// Module: Drop Point & Mobile Pickup Schedules
		dpRepo := drop_point.NewRepository(db)
		dpService := drop_point.NewService(dpRepo)
		dpHandler := drop_point.NewHandler(dpService)
		drop_point.RegisterRoutes(apiV1, dpHandler, cfg)

		// Module: Finance Report (P&L & Operating Expenses)
		finRepo := finance_report.NewRepository(db)
		finService := finance_report.NewService(finRepo)
		finHandler := finance_report.NewHandler(finService)
		finance_report.RegisterRoutes(apiV1, finHandler, cfg)

		// Module: PPOB, Tabungan Emas & Sedekah Sampah
		ppobRepo := ppob.NewRepository(db)
		ppobService := ppob.NewService(ppobRepo)
		ppobHandler := ppob.NewHandler(ppobService)
		ppob.RegisterRoutes(apiV1, ppobHandler, cfg)

		// Module: Marketplace Eco-Products
		marketRepo := marketplace.NewRepository(db)
		marketService := marketplace.NewService(marketRepo)
		marketHandler := marketplace.NewHandler(marketService)
		marketplace.RegisterRoutes(apiV1, marketHandler, cfg)

		// Module: Gamification & Leaderboard
		gameRepo := gamification.NewRepository(db)
		gameService := gamification.NewService(gameRepo)
		gameHandler := gamification.NewHandler(gameService)
		gamification.RegisterRoutes(apiV1, gameHandler, cfg)

		// Module: EPR Compliance & B2B Contracts
		compRepo := compliance.NewRepository(db)
		compService := compliance.NewService(compRepo)
		compHandler := compliance.NewHandler(compService)
		compliance.RegisterRoutes(apiV1, compHandler, cfg)
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
