package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/renlin-code/mock-shop-api/pkg/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.userSignUp)
		auth.POST("/confirm-email", h.userConfirmEmail)
		auth.POST("/sign-in", h.userSignIn)
		auth.POST("/password-recovery", h.recoveryUserPassword)
		auth.POST("/password-update", h.updateUserPassword)
	}
	profile := router.Group("/profile", h.userIdentity)
	{
		profile.GET("/", h.getUserProfile)
		profile.PUT("/", h.updateUserProfile)

		orders := profile.Group("/orders")
		{
			orders.POST("/", h.userCreateOrder)
			orders.GET("/", h.userGetAllOrder)
			orders.GET("/:id", h.userGetOrderById)
		}
	}

	api := router.Group("/api")
	{
		categories := api.Group("/categories")
		{
			categories.GET("/", h.getAllCategories)
			categories.GET("/:id", h.getCategoryById)

			products := categories.Group(":id/products")
			{
				products.GET("/", h.getCategoryProducts)
			}
		}

		products := api.Group("/products")
		{
			products.GET("/", h.getAllProducts)
		}

		{
			products.GET("/:id", h.getProductById)
		}
	}

	return router
}
