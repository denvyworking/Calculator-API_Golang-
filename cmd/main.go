package main

import (
	calculationservice "calc/iternal/CalculationService"
	"calc/iternal/db"
	"calc/iternal/handlers"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	database, err := db.InitDB()

	if err != nil {
		log.Fatalf("Could not connect to DB, %v", err)
	}

	e := echo.New()

	calcRepo := calculationservice.NewCalculationRepository(database)

	calcService := calculationservice.NewCalculationService(calcRepo)

	calcHandler := handlers.NewCalculationHandler(calcService)

	e.Use(middleware.CORS())

	e.Use(middleware.Logger())

	// все гет запросы и пост запросы по данным путям мы обрабатываем с помощью данных методов
	e.GET("/calculations", calcHandler.GetCalculations)
	e.POST("/calculations", calcHandler.PostCalculations)
	e.PATCH("/calculations/:id", calcHandler.PatchCalculations)
	e.DELETE("/calculations/:id", calcHandler.DeleteCalculations)

	e.Start("localhost:8080")
}
