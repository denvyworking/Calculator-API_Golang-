package handlers

import (
	calculationservice "calc/iternal/CalculationService"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CalculationHandler struct {
	// из пакета calculationservice достаем интерфейс CalculationService
	service calculationservice.CalculationService
}

func NewCalculationHandler(s calculationservice.CalculationService) *CalculationHandler {
	return &CalculationHandler{service: s}
}

func (h *CalculationHandler) GetCalculations(c echo.Context) error {
	calculations, err := h.service.GetAllCalculations()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not get calculations"})
	}

	return c.JSON(http.StatusOK, calculations)
}

// post ручка

func (h *CalculationHandler) PostCalculations(c echo.Context) error {
	var req calculationservice.CalculationRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	calc, err := h.service.CreateCalculation(req.Expression)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Could not create this"})
	}
	return c.JSON(http.StatusCreated, calc)
}

func (h *CalculationHandler) PatchCalculations(c echo.Context) error {
	id := c.Param("id")
	// новый запрос
	var req calculationservice.CalculationRequest

	// мы декодируем новый запрос
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	updateCalc, err := h.service.UpdateCalculation(id, req.Expression)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Could not update calculation"})
	}

	return c.JSON(http.StatusOK, updateCalc)
}

func (h *CalculationHandler) DeleteCalculations(c echo.Context) error {
	id := c.Param("id")

	err := h.service.DeleteCalculation(id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not delete calculation"})
	}

	return c.NoContent(http.StatusNoContent)
}
