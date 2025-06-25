package main

import (
	"github.com/MrWhok/IMK-FP-BACKEND/client/restclient"
	"github.com/MrWhok/IMK-FP-BACKEND/configuration"
	"github.com/MrWhok/IMK-FP-BACKEND/controller"
	_ "github.com/MrWhok/IMK-FP-BACKEND/docs"
	"github.com/MrWhok/IMK-FP-BACKEND/entity"
	"github.com/MrWhok/IMK-FP-BACKEND/exception"
	repository "github.com/MrWhok/IMK-FP-BACKEND/repository/impl"
	service "github.com/MrWhok/IMK-FP-BACKEND/service/impl"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/robfig/cron"
)

// @title Go Fiber Clean Architecture
// @version 1.0.0
// @description Baseline project using Go Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:9999
// @BasePath /
// @schemes http https
// @securityDefinitions.apikey JWT
// @in header
// @name Authorization
// @description Authorization For JWT
func main() {
	//setup configuration
	config := configuration.New()
	database := configuration.NewDatabase(config)

	err := database.AutoMigrate(
		&entity.User{},
		&entity.UserRole{},
		&entity.Product{},
		&entity.Cart{},
		&entity.CartItem{},
		&entity.Transaction{},
		&entity.TransactionDetail{},
		&entity.Gift{},
	)

	// database.AutoMigrate(&entity.Cart{})
	// database.AutoMigrate(&entity.User{})
	// database.AutoMigrate(&entity.UserRole{})
	// database.AutoMigrate(&entity.Product{})
	// err := database.AutoMigrate(&entity.Cart{})
	if err != nil {
		panic("AutoMigrate failed: " + err.Error())
	}

	redis := configuration.NewRedis(config)

	//repository
	productRepository := repository.NewProductRepositoryImpl(database)
	transactionRepository := repository.NewTransactionRepositoryImpl(database)
	transactionDetailRepository := repository.NewTransactionDetailRepositoryImpl(database)
	userRepository := repository.NewUserRepositoryImpl(database)
	cartRepository := repository.NewCartRepositoryImpl(database)
	newsRepository := repository.NewFileNewsRepo("data/cache.json")
	giftRepository := repository.NewGiftRepositoryImpl(database)

	//rest client
	httpBinRestClient := restclient.NewHttpBinRestClient()

	//service
	productService := service.NewProductServiceImpl(&productRepository, redis)
	transactionService := service.NewTransactionServiceImpl(&transactionRepository, cartRepository, productRepository, userRepository)
	transactionDetailService := service.NewTransactionDetailServiceImpl(&transactionDetailRepository)
	userService := service.NewUserServiceImpl(&userRepository)
	httpBinService := service.NewHttpBinServiceImpl(&httpBinRestClient)
	cartService := service.NewCartServiceImpl(cartRepository, productRepository)
	newsService := service.NewNewsServiceImpl(newsRepository)
	giftService := service.NewGiftServiceImpl(&giftRepository, redis)

	//controller
	productController := controller.NewProductController(&productService, config)
	transactionController := controller.NewTransactionController(&transactionService, config)
	transactionDetailController := controller.NewTransactionDetailController(&transactionDetailService, config)
	userController := controller.NewUserController(&userService, config)
	httpBinController := controller.NewHttpBinController(&httpBinService)
	cartController := controller.NewCartController(&cartService, config)
	newsController := controller.NewNewsController(&newsService)
	giftController := controller.NewGiftController(&giftService, config)

	//setup fiber
	app := fiber.New(configuration.NewFiberConfiguration())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000, https://hoppscotch.io, http://34.101.249.2:3000",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	// cronjob
	cronJob := cron.New()
	cronJob.AddFunc("0 0 * * *", newsService.FetchAndUpdate)
	cronJob.Start()
	go newsService.FetchAndUpdate()

	//routing
	productController.Route(app)
	transactionController.Route(app)
	transactionDetailController.Route(app)
	userController.Route(app)
	httpBinController.Route(app)
	cartController.Route(app)
	newsController.Route(app)
	giftController.Route(app)

	//swagger
	app.Get("/swagger/*", swagger.HandlerDefault)

	// media
	app.Static("/media", "./media")

	//start app
	err = app.Listen(config.Get("SERVER.PORT"))
	exception.PanicLogging(err)
}
