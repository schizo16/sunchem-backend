package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

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

	categoryHandler "sunchem-backend/internal/modules/categories/delivery/http"
	categoryRepo "sunchem-backend/internal/modules/categories/repository"
	categoryUC "sunchem-backend/internal/modules/categories/usecase"

	tagHandler "sunchem-backend/internal/modules/tags/delivery/http"
	tagRepo "sunchem-backend/internal/modules/tags/repository"
	tagUC "sunchem-backend/internal/modules/tags/usecase"

	seoHandler "sunchem-backend/internal/modules/seo/delivery/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// in-memory org storage (single-tenant, survives until restart)
var orgsStore = []gin.H{
	{"id": "genoract", "name": "Genoract", "slug": "genoract"},
}
var orgsMu sync.RWMutex

// staticTopPaths lists all first-path-segments after /api/v1/ that are "flat" routes.
// Requests starting with any of these are handled by the main router directly.
// All other /api/v1/<x>/... requests are treated as site-scoped (siteId is stripped).
var staticTopPaths = map[string]struct{}{
	"auth":     {},
	"posts":    {},
	"products": {},
	"media":    {},
	"settings": {},
	"users":    {},
	"traffic":  {},
	"sites":    {},
	"orgs":     {},
	"organizations": {},
	"uploads":  {},
}

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
	// All users are admins in single-tenant admin panel
	database.Exec("UPDATE users SET role = 'admin' WHERE role = 'employee'")

	userRepo := authRepo.NewUserRepository(database)
	blogRepository := blogRepo.NewBlogRepository(database)
	mediaRepository := mediaRepo.NewMediaRepository(database)
	settingRepository := settingRepo.NewSettingRepository(database)
	trafficRepository := trafficRepo.NewTrafficRepository(database)
	productRepository := productRepo.NewProductRepository(database)
	categoryRepository := categoryRepo.NewCategoryRepository(database)
	tagRepository := tagRepo.NewTagRepository(database)

	authUseCase := authUC.NewAuthUseCase(userRepo, cfg.JWTSecret, cfg.GenoractClientID, cfg.GenoractClientSecret)
	authUseCase.SetOIDCConfig(cfg.OIDCAuthority, cfg.OIDCClientID, cfg.OIDCRedirectURI)
	blogUseCase := blogUC.NewBlogUseCase(blogRepository)
	mediaUseCase := mediaUC.NewMediaUseCase(mediaRepository, cfg.UploadDir)
	settingUseCase := settingUC.NewSettingUseCase(settingRepository)
	trafficUseCase := trafficUC.NewTrafficUseCase(trafficRepository)
	productUseCase := productUC.NewProductUseCase(productRepository)
	categoryUseCase := categoryUC.NewCategoryUseCase(categoryRepository)
	tagUseCase := tagUC.NewTagUseCase(tagRepository)

	authH := authHandler.NewAuthHandler(authUseCase)
	blogH := blogHandler.NewBlogHandler(blogUseCase)
	mediaH := mediaHandler.NewMediaHandler(mediaUseCase)
	userH := userHandler.NewUserHandler(userRepo)
	settingH := settingHandler.NewSettingHandler(settingUseCase)
	trafficH := trafficHandler.NewTrafficHandler(trafficUseCase)
	productH := productHandler.NewProductHandler(productUseCase)
	categoryH := categoryHandler.NewCategoryHandler(categoryUseCase)
	tagH := tagHandler.NewTagHandler(tagUseCase)
	seoH := seoHandler.NewSEOHandler(settingUseCase)

	authMW := middleware.AuthMiddleware(cfg.JWTSecret)
	adminMW := middleware.AdminOnly(cfg.JWTSecret)

	r := gin.Default()
	r.Use(middleware.CORS())
	r.Use(middleware.ErrorHandler())

	// ── Main API router (flat routes) ──────────────────────────────────────
	api := r.Group("/api/v1")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		// Setup status — always installed (single-tenant, no setup flow)
		api.GET("/setup", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"data": gin.H{"is_installed": true}})
		})
		api.POST("/setup", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"data": gin.H{"is_installed": true}})
		})

		// Auth
		api.POST("/auth/login", authH.Login)
		api.POST("/auth/login-bundle", authH.LoginBundle)
		api.POST("/auth/genoract/callback", authH.GenoractCallback)
		api.GET("/auth/config", authH.GetConfig)
		api.POST("/auth/token", authH.Token)
		api.POST("/auth/refresh", authH.RefreshToken)

		// Traffic
		api.POST("/traffic/track", trafficH.Track)

		// Products (public)
		api.GET("/products", productH.List)
		api.GET("/products/slug/:slug", productH.GetBySlug)
		api.GET("/products/:id", productH.GetByID)

		// Blog (public)
		api.GET("/posts", blogH.List)
		api.GET("/posts/slug/:slug", blogH.GetBySlug)

		// Settings general (public)
		api.GET("/settings/general", settingH.GetGeneral)

		// Sites & Orgs (hardcoded single-tenant, public)
		api.GET("/sites/my", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"data": []gin.H{
				{"id": "sunchem", "name": "Sunchem", "domain": "sunchem.vn", "status": "active"},
			}})
		})
		api.GET("/orgs/my", func(c *gin.Context) {
			orgsMu.RLock()
			defer orgsMu.RUnlock()
			if len(orgsStore) > 0 {
				c.JSON(http.StatusOK, gin.H{"data": orgsStore[0]})
			} else {
				c.JSON(http.StatusOK, gin.H{"data": gin.H{
					"id": "genoract", "name": "Genoract", "slug": "genoract",
				}})
			}
		})
		api.GET("/organizations", func(c *gin.Context) {
			orgsMu.RLock()
			defer orgsMu.RUnlock()
			c.JSON(http.StatusOK, gin.H{"data": orgsStore})
		})

		// Protected routes
		protected := api.Group("")
		protected.Use(authMW)
		{
			protected.GET("/auth/me", authH.Me)
			protected.GET("/users/me", authH.UsersMe)

			// Traffic
			protected.GET("/traffic/summary", trafficH.Summary)
			protected.GET("/traffic/detail", trafficH.Detail)
			protected.GET("/traffic/page-views", trafficH.PageViews)

			// Blog
			protected.GET("/posts/:id", blogH.GetByID)
			protected.POST("/posts", blogH.Create)
			protected.PUT("/posts/:id", blogH.Update)
			protected.DELETE("/posts/:id", blogH.Delete)

			// Media
			protected.GET("/media", mediaH.List)
			protected.POST("/media/upload", mediaH.Upload)
			protected.DELETE("/media/:id", mediaH.Delete)

			// Settings
			protected.GET("/settings", settingH.GetAll)
			protected.PUT("/settings", settingH.Save)
			protected.GET("/settings/oidc", settingH.GetOIDC)
			protected.PUT("/settings/oidc", settingH.SaveOIDC)
			protected.GET("/settings/storage", settingH.GetStorage)
			protected.PUT("/settings/storage", settingH.SaveStorage)

			// Products
			protected.POST("/products", productH.Create)
			protected.PUT("/products/:id", productH.Update)
			protected.DELETE("/products/:id", productH.Delete)

			// Organizations (protected, not admin-only)
			protected.POST("/organizations", func(c *gin.Context) {
				var req struct {
					Name string `json:"name"`
					Slug string `json:"slug"`
				}
				c.ShouldBindJSON(&req)
				if req.Name == "" {
					req.Name = "Genoract"
				}
				if req.Slug == "" {
					req.Slug = "genoract"
				}
				org := gin.H{
					"id":   fmt.Sprintf("org_%x", time.Now().UnixMilli()),
					"name": req.Name,
					"slug": req.Slug,
				}
				orgsMu.Lock()
				orgsStore = append(orgsStore, org)
				orgsMu.Unlock()
				c.JSON(http.StatusOK, gin.H{"data": org})
			})

			// Users roles (hardcoded for frontend compatibility)
			protected.GET("/users/roles", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"data": []gin.H{
					{"id": 1, "name": "admin", "permissions": []string{"*"}},
					{"id": 2, "name": "editor", "permissions": []string{"posts:write", "posts:read"}},
					{"id": 3, "name": "author", "permissions": []string{"posts:write"}},
					{"id": 4, "name": "subscriber", "permissions": []string{"posts:read"}},
				}})
			})

			// Users listing (auth required, not admin-only)
			protected.GET("/users", userH.List)

			// Admin-only: user mutations
			admin := protected.Group("")
			admin.Use(adminMW)
			{
				admin.POST("/users", userH.Create)
				admin.PUT("/users/:id", userH.Update)
				admin.DELETE("/users/:id", userH.Delete)
			}
		}
	}

	// ── Site-scoped router (handles /:siteId/... requests) ─────────────────
	// Built as a separate Gin engine to avoid wildcard/static conflict in the main tree.
	// The main router's NoRoute handler rewrites paths and delegates here.
	siteEngine := buildSiteRouter(authMW, blogH, mediaH, settingH, categoryH, tagH, seoH)

	// NoRoute catches /api/v1/<siteId>/... paths not matched above
	r.NoRoute(func(c *gin.Context) {
		const apiPrefix = "/api/v1/"
		reqPath := c.Request.URL.Path
		if !strings.HasPrefix(reqPath, apiPrefix) {
			c.Status(http.StatusNotFound)
			return
		}
		tail := reqPath[len(apiPrefix):]
		slashIdx := strings.Index(tail, "/")
		if slashIdx < 0 {
			c.Status(http.StatusNotFound)
			return
		}
		firstSeg := tail[:slashIdx]
		// If this first segment is a known static top-level path, it's a genuine 404
		if _, isStatic := staticTopPaths[firstSeg]; isStatic {
			c.Status(http.StatusNotFound)
			return
		}
		// Otherwise treat first segment as siteId and strip it
		newPath := "/" + tail[slashIdx+1:]
		c.Request.URL.Path = newPath
		siteEngine.HandleContext(c)
	})

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

// buildSiteRouter creates a dedicated Gin engine for site-scoped routes.
// The path has already had the /:siteId prefix stripped, so routes here
// are registered without the siteId segment.
func buildSiteRouter(
	authMW gin.HandlerFunc,
	blogH *blogHandler.BlogHandler,
	mediaH *mediaHandler.MediaHandler,
	settingH *settingHandler.SettingHandler,
	categoryH *categoryHandler.CategoryHandler,
	tagH *tagHandler.TagHandler,
	seoH *seoHandler.SEOHandler,
) *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Recovery()) // recover from panics in site-scoped handlers
	// Note: CORS and ErrorHandler are already in the main engine's middleware chain.

	// Public site routes
	engine.GET("/posts", blogH.List)
	engine.GET("/categories", categoryH.List)
	engine.GET("/tags", tagH.List)
	engine.GET("/seo", seoH.GetSEO)

	// Protected site routes
	prot := engine.Group("")
	prot.Use(authMW)
	{
		// Posts
		prot.GET("/posts/:id", blogH.GetByID)
		prot.POST("/posts", blogH.Create)
		prot.PUT("/posts/:id", blogH.Update)
		prot.DELETE("/posts/:id", blogH.Delete)
		prot.PATCH("/posts/:id/moderate", blogH.Update)

		// Media
		prot.GET("/media", mediaH.List)
		prot.POST("/media/upload", mediaH.Upload)
		prot.DELETE("/media/:id", mediaH.Delete)

		// Settings
		prot.GET("/settings", settingH.GetAll)
		prot.PUT("/settings", settingH.Save)

		// Categories
		prot.POST("/categories", categoryH.Create)
		prot.PUT("/categories/:id", categoryH.Update)
		prot.DELETE("/categories/:id", categoryH.Delete)

		// Tags
		prot.POST("/tags", tagH.Create)
		prot.PUT("/tags/:id", tagH.Update)
		prot.DELETE("/tags/:id", tagH.Delete)

		// SEO
		prot.PUT("/seo", seoH.SaveSEO)
	}

	return engine
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
