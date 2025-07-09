package routes

import (
	"my-finance-hub-api/internal/infrastructure/container"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, container *container.Container) {
	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes (públicas)
		setupAuthRoutes(v1, container)

		// Protected routes
		protected := v1.Group("/")
		protected.Use(container.AuthMiddleware.RequireAuth())
		{
			setupProtectedRoutes(protected, container)
		}
	}
}

func setupAuthRoutes(group *gin.RouterGroup, container *container.Container) {
	auth := group.Group("/auth")
	{
		auth.POST("/register", container.AuthController.Register)
		auth.POST("/login", container.AuthController.Login)
		auth.POST("/logout", container.AuthController.Logout)
		auth.GET("/me", container.AuthMiddleware.RequireAuth(), container.AuthController.GetCurrentUser)
		auth.PUT("/profile", container.AuthMiddleware.RequireAuth(), container.AuthController.UpdateUser)
		auth.PUT("/password", container.AuthMiddleware.RequireAuth(), container.AuthController.ChangePassword)
	}
}

func setupProtectedRoutes(group *gin.RouterGroup, container *container.Container) {
	// Categories routes
	categories := group.Group("/categories")
	{
		categories.GET("/", container.CategoryController.GetCategories)
		categories.GET("", container.CategoryController.GetCategories)
		categories.POST("/", container.CategoryController.CreateCategory)
		categories.POST("", container.CategoryController.CreateCategory)
		categories.GET("/:id", container.CategoryController.GetCategory)
		categories.PUT("/:id", container.CategoryController.UpdateCategory)
		categories.PATCH("/:id", container.CategoryController.UpdateCategory)
		categories.DELETE("/:id", container.CategoryController.DeleteCategory)
	}

	// Goals routes
	goals := group.Group("/goals")
	{
		goals.GET("/", container.GoalController.GetGoals)
		goals.GET("", container.GoalController.GetGoals)
		goals.POST("/", container.GoalController.CreateGoal)
		goals.POST("", container.GoalController.CreateGoal)
		goals.GET("/:id", container.GoalController.GetGoal)
		goals.PUT("/:id", container.GoalController.UpdateGoal)
		goals.PATCH("/:id", container.GoalController.UpdateGoal)
		goals.DELETE("/:id", container.GoalController.DeleteGoal)
	}

	// Saving Goals (PiggyBanks) routes
	savingGoals := group.Group("/saving_goals")
	{
		savingGoals.GET("/", container.SavingGoalController.GetSavingGoals)
		savingGoals.GET("", container.SavingGoalController.GetSavingGoals)
		savingGoals.POST("/", container.SavingGoalController.CreateSavingGoal)
		savingGoals.POST("", container.SavingGoalController.CreateSavingGoal)
		savingGoals.GET("/:id", container.SavingGoalController.GetSavingGoal)
		savingGoals.PUT("/:id", container.SavingGoalController.UpdateSavingGoal)
		savingGoals.PATCH("/:id", container.SavingGoalController.UpdateSavingGoal)
		savingGoals.DELETE("/:id", container.SavingGoalController.DeleteSavingGoal)
		savingGoals.POST("/:id/deposit", container.SavingGoalController.Deposit)
	}

	// Transactions routes
	transactions := group.Group("/transactions")
	{
		transactions.GET("/", container.TransactionController.GetTransactions)
		transactions.GET("", container.TransactionController.GetTransactions)
		transactions.POST("/", container.TransactionController.CreateTransaction)
		transactions.POST("", container.TransactionController.CreateTransaction)
		transactions.GET("/:id", container.TransactionController.GetTransaction)
		transactions.PUT("/:id", container.TransactionController.UpdateTransaction)
		transactions.PATCH("/:id", container.TransactionController.UpdateTransaction)
		transactions.DELETE("/:id", container.TransactionController.DeleteTransaction)
		transactions.PATCH("/:id/paid", container.TransactionController.TogglePaid)
		// Endpoint específico para relatórios do dashboard
		transactions.GET("/reports", container.TransactionController.GetDashboardReports)
		transactions.GET("/reports/", container.TransactionController.GetDashboardReports)
	}

	// Reports routes
	reports := group.Group("/reports")
	{
		reports.GET("/", container.TransactionController.GetReports)
		reports.GET("", container.TransactionController.GetReports)
	}
}
