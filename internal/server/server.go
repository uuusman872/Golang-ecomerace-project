package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/uuusman872/learning-go-shop/internal/config"
	service "github.com/uuusman872/learning-go-shop/internal/services"
	"gorm.io/gorm"
)

type Server struct {
	config         *config.Config
	db             *gorm.DB
	logger         *zerolog.Logger
	authServices   *service.AuthService
	productService *service.ProductService
	userService    *service.UserService
	uploadService  *service.UploadService
	cartService    *service.CartService
	orderService   *service.OrderService
}

func New(
	cfg *config.Config,
	db *gorm.DB,
	logger *zerolog.Logger,
	authService *service.AuthService,
	productService *service.ProductService,
	userService *service.UserService,
	uploadService *service.UploadService,
	cartService *service.CartService,
	orderService *service.OrderService,
) *Server {
	return &Server{
		config:         cfg,
		db:             db,
		logger:         logger,
		authServices:   authService,
		productService: productService,
		userService:    userService,
		uploadService:  uploadService,
		cartService:    cartService,
		orderService:   orderService,
	}
}

func (s *Server) SetupRoutes() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(s.corsMiddleware())
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/health", s.healthCheck)
	router.Static("/uploads", "./uploads")
	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", s.register)
			auth.POST("/login", s.login)
			auth.POST("/refresh", s.refreshToken)
			auth.POST("/logout", s.logout)
		}
		protected := api.Group("/")
		protected.Use(s.authMiddleware())
		{
			users := protected.Group("/users")
			{
				userRoutes := users
				userRoutes.GET("/profile", s.getProfile)
				userRoutes.PUT("/profile", s.updateProfile)
			}

			categories := protected.Group("/categories")
			{
				categoryRoute := categories
				categoryRoute.POST("/", s.adminMiddleware(), s.createCategory)
				categoryRoute.PUT("/:id", s.adminMiddleware(), s.updateCategory)
				categoryRoute.DELETE("/:id", s.adminMiddleware(), s.deleteCategory)
			}

			products := protected.Group("/products")
			{
				productRoutes := products
				productRoutes.POST("/", s.adminMiddleware(), s.createProduct)
				productRoutes.PUT("/:id", s.adminMiddleware(), s.updateProduct)
				productRoutes.DELETE("/:id", s.adminMiddleware(), s.deleteProduct)
				productRoutes.POST("/:id/images", s.adminMiddleware(), s.uploadProductImage)
			}
			cart := protected.Group("/cart")
			{
				cartRoute := cart
				cartRoute.GET("/", s.getCart)
				cartRoute.POST("/items", s.addToCart)
				cartRoute.PUT("/items/:id", s.updateCartItem)
				cartRoute.DELETE("/items/:id", s.removeFromCart)
			}
			order := protected.Group("/orders")
			{
				orderRoutes := order
				orderRoutes.POST("/", s.createOrder)
				orderRoutes.GET("/", s.getOrders)
				orderRoutes.GET("/:id", s.getOrder)
			}
		}

	}
	api.GET("/categories", s.getCategory)
	api.GET("/products", s.getProducts)
	api.GET("/products/:id", s.getProduct)
	return router
}

// healthCheck godoc
// @Summary      Health check
// @Tags         Health
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /health [get]
func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s *Server) corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
