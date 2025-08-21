package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Knetic/govaluate"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func initDB() {
	dsn := "host=localhost user=postgres password=yourpassword dbname=postgres port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	if err := db.AutoMigrate(&Calculation{}); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}
}

type Calculation struct {
	ID         string `gorm:"primaryKey" json:"id"`
	Expression string `json:"expression"`
	Result     string `json:"result"`
}

type CalculationRequest struct {
	Expression string `json:"expression"`
}

//var calculations = []Calculation{}

func CalculateExpression(expression string) (string, error) {
	// передаем стоку с выражением в виде стоки : 52+52
	expr, err := govaluate.NewEvaluableExpression(expression)
	// отслеживание невалидной строки например : 34++34
	if err != nil {
		return "", err
	}
	result, err := expr.Evaluate(nil)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", result), nil
}

// get ручка

func getCalculations(c echo.Context) error {
	var calculations []Calculation

	if err := db.Find(&calculations).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not get calculations"})
	}
	return c.JSON(http.StatusOK, calculations)
}

// post ручка

func postCalculations(c echo.Context) error {
	var req CalculationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	result, err := CalculateExpression(req.Expression)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid expression"})
	}

	calc := Calculation{
		ID:         uuid.NewString(),
		Expression: req.Expression,
		Result:     result,
	}

	if err := db.Create(&calc).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not add calculations"})
	}

	return c.JSON(http.StatusCreated, calc)
}

func patchCalculations(c echo.Context) error {
	id := c.Param("id")
	// новый запрос
	var req CalculationRequest
	// мы декодируем новый запрос
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	// считаем новый результат
	result, err := CalculateExpression(req.Expression)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid expression"})
	}

	var calc Calculation
	if err := db.First(&calc, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Could not find expression"})
	}
	calc.Expression = req.Expression
	calc.Result = result

	if err := db.Save(&calc).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not update calculations"})
	}
	return c.JSON(http.StatusOK, calc)
}

func deleteCalculations(c echo.Context) error {
	id := c.Param("id")

	result := db.Delete(&Calculation{}, "id = ?", id)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not delete calculation"})
	}

	if result.RowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Calculation not found"})
	}

	return c.NoContent(http.StatusNoContent)
}

func main() {
	initDB()

	e := echo.New()

	e.Use(middleware.CORS())

	e.Use(middleware.Logger())

	// все гет запросы и пост запросы по данным путям мы обрабатываем с помощью данных методов
	e.GET("/calculations", getCalculations)
	e.POST("/calculations", postCalculations)
	e.PATCH("/calculations/:id", patchCalculations)
	e.DELETE("/calculations/:id", deleteCalculations)

	e.Start("localhost:8080")
}
