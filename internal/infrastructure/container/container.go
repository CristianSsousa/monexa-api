package container

import (
	"my-finance-hub-api/internal/application/interfaces"
	"my-finance-hub-api/internal/application/services"
	"my-finance-hub-api/internal/domain/repositories"
	dbRepos "my-finance-hub-api/internal/infrastructure/database/repositories"
	"my-finance-hub-api/internal/infrastructure/http/controllers"
	"my-finance-hub-api/internal/infrastructure/http/middleware"

	"gorm.io/gorm"
)

type Container struct {
	// Database
	DB *gorm.DB

	// Repositories
	UserRepository        repositories.UserRepository
	CategoryRepository    repositories.CategoryRepository
	GoalRepository        repositories.GoalRepository
	SavingGoalRepository  repositories.SavingGoalRepository
	TransactionRepository repositories.TransactionRepository

	// Services
	AuthService        interfaces.AuthService
	CategoryService    interfaces.CategoryService
	GoalService        interfaces.GoalService
	SavingGoalService  interfaces.SavingGoalService
	TransactionService interfaces.TransactionService

	// Controllers
	AuthController        *controllers.AuthController
	CategoryController    *controllers.CategoryController
	GoalController        *controllers.GoalController
	SavingGoalController  *controllers.SavingGoalController
	TransactionController *controllers.TransactionController

	// Middleware
	AuthMiddleware *middleware.AuthMiddleware
}

func NewContainer(db *gorm.DB) *Container {
	container := &Container{
		DB: db,
	}

	container.initRepositories()
	container.initServices()
	container.initControllers()
	container.initMiddleware()

	return container
}

func (c *Container) initRepositories() {
	c.UserRepository = dbRepos.NewUserRepository(c.DB)
	c.CategoryRepository = dbRepos.NewCategoryRepository(c.DB)
	c.GoalRepository = dbRepos.NewGoalRepository(c.DB)
	c.SavingGoalRepository = dbRepos.NewSavingGoalRepository(c.DB)
	c.TransactionRepository = dbRepos.NewTransactionRepository(c.DB)
}

func (c *Container) initServices() {
	c.AuthService = services.NewAuthService(c.UserRepository)
	c.CategoryService = services.NewCategoryService(c.CategoryRepository)
	c.GoalService = services.NewGoalService(c.GoalRepository)
	c.SavingGoalService = services.NewSavingGoalService(c.SavingGoalRepository)
	c.TransactionService = services.NewTransactionService(c.TransactionRepository)
}

func (c *Container) initControllers() {
	c.AuthController = controllers.NewAuthController(c.AuthService)
	c.CategoryController = controllers.NewCategoryController(c.CategoryService)
	c.GoalController = controllers.NewGoalController(c.GoalService)
	c.SavingGoalController = controllers.NewSavingGoalController(c.SavingGoalService)
	c.TransactionController = controllers.NewTransactionController(c.TransactionService)
}

func (c *Container) initMiddleware() {
	c.AuthMiddleware = middleware.NewAuthMiddleware(c.AuthService)
}
