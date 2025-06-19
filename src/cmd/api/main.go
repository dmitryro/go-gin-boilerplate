package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go_api/docs"
	"go_api/internal/handlers"
	"go_api/internal/middleware"
	"go_api/internal/models"
	"go_api/internal/services"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title GO API
// @version 1.0
// @description This is the API documentation for GO API.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@medrobotix.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8081
// @BasePath /
// @schemes http

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter the token with the 'Bearer ' prefix, e.g., 'Bearer eyJhbGci...'

func main() {
	var err error

	// Resolve root path
	rootPath, err := filepath.Abs(filepath.Join("./"))
	if err != nil {
		log.Fatal("Unable to resolve root path:", err)
	}

	envPath := filepath.Join(rootPath, ".env")

	log.Println("Loaded .env from:", envPath)
	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf("Error loading .env file from %s: %v", envPath, err)
	}

	// Swagger info config
	docs.SwaggerInfo.Host = os.Getenv("API_HOST")
	if docs.SwaggerInfo.Host == "" {
		docs.SwaggerInfo.Host = "localhost:8081"
	}
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http"}

	// DB connection setup
	dsn := "host=" + os.Getenv("PG_HOST") +
		" user=" + os.Getenv("PG_USER") +
		" password=" + os.Getenv("PG_PASSWORD") +
		" dbname=" + os.Getenv("PG_DATABASE") +
		" port=" + os.Getenv("PG_PORT") +
		" sslmode=disable client_encoding=UTF8"

	var db *gorm.DB
	maxRetries := 10
	retryInterval := 5 * time.Second

	log.Printf("Connection string is: %s", dsn)

	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
		})
		if err == nil {
			sqlDB, _ := db.DB()
			err = sqlDB.Ping()
		}
		if err == nil {
			log.Println("Successfully connected to the database")
			break
		}
		log.Printf("Failed to connect to DB (attempt %d/%d): %v. Retrying in %v...", i+1, maxRetries, err, retryInterval)
		time.Sleep(retryInterval)
	}

	if err != nil {
		log.Fatalf("Failed to connect to DB after %d attempts: %v", maxRetries, err)
	}

	// Ensure unique constraint on roles.name
	sql := `DO $$
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'uni_roles_name') THEN
		ALTER TABLE roles ADD CONSTRAINT uni_roles_name UNIQUE (name);
	END IF;
END $$;`
	if err := db.Exec(sql).Error; err != nil {
		log.Fatalf("Failed to ensure unique constraint: %v", err)
	}

	// Migrate DB schema
	if err := db.AutoMigrate(&models.User{}, &models.Role{}); err != nil {
		log.Fatalf("Failed to migrate DB: %v", err)
	}

	// Create Gin router
	r := gin.Default()

	// Swagger UI route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Instantiate services
	loginService := services.NewLoginService(db)
	registerService := services.NewRegisterService(db)
	userService := services.NewUserService(db)
	roleService := services.NewRoleService(db)

	// Instantiate handlers
	loginHandler := handlers.NewLoginHandler(loginService)
	registerHandler := handlers.NewRegisterHandler(registerService)
	userHandler := handlers.NewUserHandler(userService)
	roleHandler := handlers.NewRoleHandler(roleService)

	// Public (no auth) routes

	// @Summary User login
	// @Tags auth
	// @Accept json
	// @Produce json
	// @Param credentials body handlers.LoginRequest true "User credentials"
	// @Success 200 {object} handlers.LoginResponse
	// @Failure 401 {object} handlers.ErrorResponse
	// @Router /api/login [post]
	r.POST("/api/login", loginHandler.Login)

	// @Summary User registration
	// @Tags auth
	// @Accept json
	// @Produce json
	// @Param registration body handlers.RegisterRequest true "User registration data"
	// @Success 201 {object} handlers.RegisterResponse
	// @Failure 400 {object} handlers.ErrorResponse
	// @Router /api/register [post]
	r.POST("/api/register", registerHandler.Register)

	// Protected API group with JWT auth applied
	api := r.Group("/api")
	api.Use(middleware.JwtAuthMiddleware())

	// User routes - any authenticated user
	userRoutes := api.Group("/users")
	userRoutes.Use(middleware.RoleAuthMiddleware("any"))
	{
		// @Summary Get all users
		// @Description Get all users
		// @Tags users
		// @Produce json
		// @Success 200 {array} models.User
		// @Security BearerAuth
		// @Router /api/users [get]
		userRoutes.GET("", userHandler.GetAllUsers)

		// @Summary Get a user by ID
		// @Description Get a user by ID
		// @Tags users
		// @Produce json
		// @Param id path int true "User ID"
		// @Success 200 {object} models.User
		// @Failure 400 {object} handlers.ErrorResponse
		// @Failure 404 {object} handlers.ErrorResponse
		// @Security BearerAuth
		// @Router /api/users/{id} [get]
		userRoutes.GET("/:id", userHandler.GetUserByID)

		// @Summary Create a new user
		// @Description Create a new user with the provided details
		// @Tags users
		// @Accept json
		// @Produce json
		// @Param user body models.UserCreateRequest true "User information"
		// @Success 201 {object} models.User
		// @Failure 400 {object} handlers.ErrorResponse
		// @Security BearerAuth
		// @Router /api/users [post]
		userRoutes.POST("", userHandler.CreateUser)

		// @Summary Update an existing user
		// @Description Update an existing user
		// @Tags users
		// @Accept json
		// @Produce json
		// @Param id path int true "User ID"
		// @Param user body models.User true "User object to be updated"
		// @Success 200 {object} models.User
		// @Failure 400 {object} handlers.ErrorResponse
		// @Security BearerAuth
		// @Router /api/users/{id} [put]
		userRoutes.PUT("/:id", userHandler.UpdateUser)

		// @Summary Delete a user
		// @Description Delete a user
		// @Tags users
		// @Produce json
		// @Param id path int true "User ID"
		// @Success 204 "No Content"
		// @Failure 400 {object} handlers.ErrorResponse
		// @Security BearerAuth
		// @Router /api/users/{id} [delete]
		userRoutes.DELETE("/:id", userHandler.DeleteUser)

		// @Summary Get a user by email
		// @Description Get a user by email
		// @Tags users
		// @Produce json
		// @Param email path string true "User Email"
		// @Success 200 {object} models.User
		// @Failure 404 {object} handlers.ErrorResponse
		// @Security BearerAuth
		// @Router /api/users/email/{email} [get]
		userRoutes.GET("/email/:email", userHandler.GetUserByEmail)

		// @Summary Get a user by username
		// @Description Get a user by username
		// @Tags users
		// @Produce json
		// @Param username path string true "User Username"
		// @Success 200 {object} models.User
		// @Failure 404 {object} handlers.ErrorResponse
		// @Security BearerAuth
		// @Router /api/users/username/{username} [get]
		userRoutes.GET("/username/:username", userHandler.GetUserByUsername)

		// @Summary Get users by role ID
		// @Description Get all users with the given role ID
		// @Tags users
		// @Produce json
		// @Param role_id path int true "Role ID"
		// @Success 200 {array} models.User
		// @Failure 400 {object} handlers.ErrorResponse
		// @Security BearerAuth
		// @Router /api/users/role/{role_id} [get]
		userRoutes.GET("/role/:role_id", userHandler.GetUsersByRoleID)

		// @Summary Change current user's password
		// @Description Change the current user's password
		// @Tags users
		// @Accept json
		// @Produce json
		// @Param body body handlers.ChangePasswordRequest true "Old and new passwords"
		// @Success 200 {object} map[string]string
		// @Failure 400 {object} handlers.ErrorResponse
		// @Failure 401 {object} handlers.ErrorResponse
		// @Security BearerAuth
		// @Router /api/users/password [post]
		userRoutes.POST("/password", userHandler.ChangePassword)
	}

	// Role routes - only admin role
	roleRoutes := api.Group("/roles")
	roleRoutes.Use(middleware.RoleAuthMiddleware("admin"))
	{
		// @Summary Get all roles
		// @Description Get all roles
		// @Tags roles
		// @Produce json
		// @Success 200 {array} models.Role
		// @Security BearerAuth
		// @Router /api/roles [get]
		roleRoutes.GET("", roleHandler.GetRoles)

		// @Summary Get a role by ID
		// @Description Get a role by ID
		// @Tags roles
		// @Produce json
		// @Param id path int true "Role ID"
		// @Success 200 {object} models.Role
		// @Failure 404 {object} handlers.ErrorResponse
		// @Security BearerAuth
		// @Router /api/roles/{id} [get]
		roleRoutes.GET("/:id", roleHandler.GetRoleByID)

		// @Summary Get a role by name
		// @Description Get a role by name
		// @Tags roles
		// @Produce json
		// @Param name path string true "Role Name"
		// @Success 200 {object} models.Role
		// @Failure 404 {object} handlers.ErrorResponse
		// @Security BearerAuth
		// @Router /api/roles/name/{name} [get]
		roleRoutes.GET("/name/:name", roleHandler.GetRoleByName)

		// @Summary Create a new role
		// @Description Create a new role
		// @Tags roles
		// @Accept json
		// @Produce json
		// @Param role body models.RoleCreateRequest true "Role information"
		// @Success 201 {object} models.Role
		// @Failure 400 {object} handlers.ErrorResponse
		// @Security BearerAuth
		// @Router /api/roles [post]
		roleRoutes.POST("", roleHandler.CreateRole)
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on :%s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

