package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"sunchem-backend/internal/common/config"
	"sunchem-backend/internal/common/db"
	"sunchem-backend/internal/common/middleware"

	authHandler "sunchem-backend/internal/modules/auth/delivery/http"
	authRepo "sunchem-backend/internal/modules/auth/repository"
	authUC "sunchem-backend/internal/modules/auth/usecase"

	blogHandler "sunchem-backend/internal/modules/blog/delivery/http"
	blogRepo "sunchem-backend/internal/modules/blog/repository"
	blogUC "sunchem-backend/internal/modules/blog/usecase"

	mediaHandler "sunchem-backend/internal/modules/media/delivery/http"
	mediaRepo "sunchem-backend/internal/modules/media/repository"
	mediaUC "sunchem-backend/internal/modules/media/usecase"

	userHandler "sunchem-backend/internal/modules/users/delivery/http"

	settingHandler "sunchem-backend/internal/modules/settings/delivery/http"
	settingRepo "sunchem-backend/internal/modules/settings/repository"
	settingUC "sunchem-backend/internal/modules/settings/usecase"

	trafficHandler "sunchem-backend/internal/modules/analytics/delivery/http"
	trafficRepo "sunchem-backend/internal/modules/analytics/repository"
	trafficUC "sunchem-backend/internal/modules/analytics/usecase"

	productHandler "sunchem-backend/internal/modules/products/delivery/http"
	productRepo "sunchem-backend/internal/modules/products/repository"
	productUC "sunchem-backend/internal/modules/products/usecase"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Run() {
	_ = godotenv.Load()

	cfg := config.LoadConfig()

	database := db.InitDB(cfg)

	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		if err := db.AutoMigrate(database); err != nil {
			log.Fatal(err)
		}
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "seed" {
		if err := db.AutoMigrate(database); err != nil {
			log.Fatal(err)
		}
		SeedDB(database)
		return
	}

	if err := db.AutoMigrate(database); err != nil {
		log.Fatal(err)
	}
	if cfg.AppEnv == "dev" {
		seedAdmin(database)
	}
	SeedDB(database)

	userRepo := authRepo.NewUserRepository(database)
	blogRepository := blogRepo.NewBlogRepository(database)
	mediaRepository := mediaRepo.NewMediaRepository(database)
	settingRepository := settingRepo.NewSettingRepository(database)
	trafficRepository := trafficRepo.NewTrafficRepository(database)
	productRepository := productRepo.NewProductRepository(database)

	authUseCase := authUC.NewAuthUseCase(userRepo, cfg.JWTSecret, cfg.GenoractClientID, cfg.GenoractClientSecret)
	blogUseCase := blogUC.NewBlogUseCase(blogRepository)
	mediaUseCase := mediaUC.NewMediaUseCase(mediaRepository, cfg.UploadDir)
	settingUseCase := settingUC.NewSettingUseCase(settingRepository)
	trafficUseCase := trafficUC.NewTrafficUseCase(trafficRepository)
	productUseCase := productUC.NewProductUseCase(productRepository)

	authH := authHandler.NewAuthHandler(authUseCase)
	blogH := blogHandler.NewBlogHandler(blogUseCase)
	mediaH := mediaHandler.NewMediaHandler(mediaUseCase)
	userH := userHandler.NewUserHandler(userRepo)
	settingH := settingHandler.NewSettingHandler(settingUseCase)
	trafficH := trafficHandler.NewTrafficHandler(trafficUseCase)
	productH := productHandler.NewProductHandler(productUseCase)

	authMW := middleware.AuthMiddleware(cfg.JWTSecret)
	adminMW := middleware.AdminOnly(cfg.JWTSecret)

	r := gin.Default()
	r.Use(middleware.CORS())
	r.Use(middleware.ErrorHandler())

	api := r.Group("/api/v1")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})
		api.POST("/auth/login", authH.Login)
		api.POST("/auth/genoract/callback", authH.GenoractCallback)
		api.POST("/traffic/track", trafficH.Track)

		api.GET("/products", productH.List)
		api.GET("/products/:id", productH.GetByID)
		api.GET("/products/slug/:slug", productH.GetBySlug)

		api.GET("/posts", blogH.List)
		api.GET("/posts/slug/:slug", blogH.GetBySlug)

		protected := api.Group("")
		protected.Use(authMW)
		{
			protected.GET("/auth/me", authH.Me)

			protected.GET("/traffic/summary", trafficH.Summary)
			protected.GET("/traffic/detail", trafficH.Detail)
			protected.GET("/traffic/page-views", trafficH.PageViews)

			protected.GET("/posts/:id", blogH.GetByID)
			protected.POST("/posts", blogH.Create)
			protected.PUT("/posts/:id", blogH.Update)
			protected.DELETE("/posts/:id", blogH.Delete)

			protected.GET("/media", mediaH.List)
			protected.POST("/media/upload", mediaH.Upload)
			protected.DELETE("/media/:id", mediaH.Delete)

			protected.GET("/settings", settingH.GetAll)
			protected.PUT("/settings", settingH.Save)

			protected.POST("/products", productH.Create)
			protected.PUT("/products/:id", productH.Update)
			protected.DELETE("/products/:id", productH.Delete)

			admin := protected.Group("")
			admin.Use(adminMW)
			{
				admin.GET("/users", userH.List)
				admin.POST("/users", userH.Create)
				admin.PUT("/users/:id", userH.Update)
				admin.DELETE("/users/:id", userH.Delete)
			}
		}
	}

	r.Static("/uploads", cfg.UploadDir)

	srv := fmt.Sprintf(":%s", cfg.ServerPort)
	fmt.Printf("[Server] Running on http://localhost%s\n", srv)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := r.Run(srv); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	<-quit
	fmt.Println("\n[Server] Shutting down...")
}

func seedAdmin(database *gorm.DB) {
	hashed, _ := bcrypt.GenerateFromPassword([]byte("sunchem2024"), bcrypt.DefaultCost)
	users := []struct {
		ID       uint
		Username string
		Password string
		Name     string
		Role     string
	}{
		{1, "admin", string(hashed), "Quản trị viên", "admin"},
		{2, "nhanvien", string(hashed), "Nhân viên kinh doanh", "employee"},
		{3, "marketing", string(hashed), "Nhân viên marketing", "employee"},
	}
	for _, u := range users {
		database.Exec(`INSERT OR IGNORE INTO users (id, username, password, name, role) VALUES (?, ?, ?, ?, ?)`,
			u.ID, u.Username, u.Password, u.Name, u.Role)
	}
	_ = hashed
}
