package routes

import (
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"go_api/internal/handlers"
	"go_api/internal/middleware"
	"go_api/internal/services"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter sets up all the routes for the application
func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// Enable CORS middleware with custom configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	swaggerYamlDir := os.Getenv("SWAGGER_YAML_DIR")
	swaggerJsonDir := os.Getenv("SWAGGER_JSON_DIR")

	// Serve swagger docs
	r.GET("/docs/swagger.json", func(c *gin.Context) {
		c.File(swaggerYamlDir)
	})
	r.GET("/docs/swagger.yaml", func(c *gin.Context) {
		c.File(swaggerJsonDir)
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome to the API!"})
	})

	// Public routes: login and register
	loginService := services.NewLoginService(db)
	loginHandler := handlers.NewLoginHandler(loginService)
	r.POST("/api/login", loginHandler.Login)

	registerService := services.NewRegisterService(db)
	registerHandler := handlers.NewRegisterHandler(registerService)
	r.POST("/api/register", registerHandler.Register)

	// Protected routes with JWT middleware and permission-based access control
	protected := r.Group("/api")
	protected.Use(middleware.JwtAuthMiddleware()) // Apply JWT middleware to the entire group

	// User routes with RBAC permissions
	userService := services.NewUserService(db)
	userHandler := handlers.NewUserHandler(userService)
	protected.GET("/users", middleware.PermissionAuthMiddleware("read"), userHandler.GetAllUsers)
	protected.POST("/users", middleware.PermissionAuthMiddleware("create"), userHandler.CreateUser)
	protected.GET("/users/:id", middleware.PermissionAuthMiddleware("read"), userHandler.GetUserByID)
	protected.GET("/users/email/:email", middleware.PermissionAuthMiddleware("read"), userHandler.GetUserByEmail)
	protected.GET("/users/username/:username", middleware.PermissionAuthMiddleware("read"), userHandler.GetUserByUsername)

	// New route to get users by role ID
	protected.GET("/users/role/:role_id", middleware.PermissionAuthMiddleware("read"), userHandler.GetUsersByRoleID)

	// Role routes with RBAC permissions
	roleService := services.NewRoleService(db)
	roleHandler := handlers.NewRoleHandler(roleService)
	protected.GET("/roles", middleware.PermissionAuthMiddleware("read"), roleHandler.GetRoles)
	protected.POST("/roles", middleware.PermissionAuthMiddleware("create"), roleHandler.CreateRole)
	protected.GET("/roles/:id", middleware.PermissionAuthMiddleware("read"), roleHandler.GetRoleByID)
	protected.GET("/roles/name/:name", middleware.PermissionAuthMiddleware("read"), roleHandler.GetRoleByName)

	return r
}

