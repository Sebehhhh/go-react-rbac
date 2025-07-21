package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"rbac-system/backend/internal/config"
	"rbac-system/backend/internal/database"
	"rbac-system/backend/internal/handlers"
	"rbac-system/backend/internal/middleware"
	"rbac-system/backend/internal/services"
	"rbac-system/backend/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	cfg := config.Load()

	if err := database.Connect(cfg); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := database.Migrate(); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	if err := database.SeedDatabase(cfg); err != nil {
		log.Fatal("Failed to seed database:", err)
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return utils.SendError(c, fiber.StatusInternalServerError, "internal_error", err.Error())
		},
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORS.AllowedOrigins[0],
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
	}))

	jwtService := utils.NewJWTService(cfg)
	authService := services.NewAuthService(jwtService)
	userService := services.NewUserService(services.NewRBACService())
	roleService := services.NewRoleService()
	dashboardService := services.NewDashboardService()
	rbacService := services.NewRBACService()

	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService, rbacService)
	roleHandler := handlers.NewRoleHandler(roleService)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)

	api := app.Group("/api")

	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.RefreshToken)
	auth.Post("/logout", middleware.AuthMiddleware(authService, jwtService), authHandler.Logout)

	profile := api.Group("/profile")
	profile.Use(middleware.AuthMiddleware(authService, jwtService))
	profile.Get("/", authHandler.Profile)
	profile.Put("/", authHandler.UpdateProfile)
	profile.Put("/password", authHandler.UpdatePassword)

	users := api.Group("/users")
	users.Use(middleware.AuthMiddleware(authService, jwtService))
	users.Get("/", middleware.RequirePermission(rbacService, "users", "read"), userHandler.GetUsers)
	users.Post("/", middleware.RequirePermission(rbacService, "users", "create"), userHandler.CreateUser)
	users.Get("/:id", middleware.SelfOrPermission(rbacService, "users", "read"), userHandler.GetUser)
	users.Put("/:id", middleware.SelfOrPermission(rbacService, "users", "update"), userHandler.UpdateUser)
	users.Delete("/:id", middleware.RequirePermission(rbacService, "users", "delete"), userHandler.DeleteUser)
	users.Put("/:id/activate", middleware.RequirePermission(rbacService, "users", "update"), userHandler.ActivateUser)
	users.Put("/:id/deactivate", middleware.RequirePermission(rbacService, "users", "update"), userHandler.DeactivateUser)
	users.Get("/:id/activity", middleware.SelfOrPermission(rbacService, "activity_logs", "read"), userHandler.GetUserActivity)
	users.Post("/bulk-actions", middleware.RequireRole(rbacService, "Super Admin"), userHandler.BulkActions)

	roles := api.Group("/roles")
	roles.Use(middleware.AuthMiddleware(authService, jwtService))
	roles.Get("/", middleware.RequirePermission(rbacService, "roles", "read"), roleHandler.GetRoles)
	roles.Post("/", middleware.RequireRole(rbacService, "Super Admin"), roleHandler.CreateRole)
	roles.Get("/:id", middleware.RequirePermission(rbacService, "roles", "read"), roleHandler.GetRole)
	roles.Put("/:id", middleware.RequireRole(rbacService, "Super Admin"), roleHandler.UpdateRole)
	roles.Delete("/:id", middleware.RequireRole(rbacService, "Super Admin"), roleHandler.DeleteRole)
	roles.Put("/:id/permissions", middleware.RequireRole(rbacService, "Super Admin"), roleHandler.AssignPermissions)
	roles.Get("/:id/permissions", middleware.RequirePermission(rbacService, "permissions", "read"), roleHandler.GetRolePermissions)

	permissions := api.Group("/permissions")
	permissions.Use(middleware.AuthMiddleware(authService, jwtService))
	permissions.Get("/", middleware.RequirePermission(rbacService, "permissions", "read"), roleHandler.GetPermissions)

	dashboard := api.Group("/dashboard")
	dashboard.Use(middleware.AuthMiddleware(authService, jwtService))
	dashboard.Get("/stats", middleware.RequirePermission(rbacService, "dashboard", "read"), dashboardHandler.GetStats)
	dashboard.Get("/role-distribution", middleware.RequirePermission(rbacService, "dashboard", "read"), dashboardHandler.GetRoleDistribution)
	dashboard.Get("/recent-activity", middleware.RequirePermission(rbacService, "activity_logs", "read"), dashboardHandler.GetRecentActivity)
	dashboard.Get("/user-analytics", middleware.RequirePermission(rbacService, "dashboard", "read"), dashboardHandler.GetUserAnalytics)
	dashboard.Get("/system-health", middleware.RequireRoleAny(rbacService, "Super Admin", "Admin"), dashboardHandler.GetSystemHealth)

	app.Use(middleware.ActivityLogger())

	app.Get("/health", func(c *fiber.Ctx) error {
		return utils.SendSuccess(c, fiber.StatusOK, "Server is healthy", map[string]string{
			"status": "ok",
			"env":    cfg.Server.Env,
		})
	})

	go func() {
		if err := app.Listen(":" + cfg.Server.Port); err != nil {
			log.Fatal("Failed to start server:", err)
		}
	}()

	log.Printf("Server starting on port %s", cfg.Server.Port)
	log.Printf("Environment: %s", cfg.Server.Env)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	if err := app.Shutdown(); err != nil {
		log.Fatal("Failed to shutdown server:", err)
	}
	log.Println("Server shutdown complete")
}