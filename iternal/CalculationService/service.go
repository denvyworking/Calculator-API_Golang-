package calculationservice

import (
	"fmt"

	"github.com/Knetic/govaluate"
	"github.com/google/uuid"
)

// бизнес логика нашего бекенда

type CalculationService interface {
	CreateCalculation(expression string) (Calculation, error)
	GetAllCalculations() ([]Calculation, error)
	GetCalculationsByID(id string) (Calculation, error)
	UpdateCalculation(id, expression string) (Calculation, error)
	DeleteCalculation(id string) error
}

type calcService struct {
	repo CalculationRepository
}

func (s *calcService) calculateExpression(expression string) (string, error) {
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

func NewCalculationService(r CalculationRepository) CalculationService {
	return &calcService{repo: r}
}

func (s *calcService) CreateCalculation(expression string) (Calculation, error) {
	// считаем то что передали нам
	result, err := s.calculateExpression(expression)
	if err != nil {
		return Calculation{}, err
	}

	calc := Calculation{
		ID:         uuid.NewString(),
		Expression: expression,
		Result:     result,
	}

	if err := s.repo.CreateCalculation(calc); err != nil {
		return Calculation{}, err
	}

	return calc, err

}

func (s *calcService) GetAllCalculations() ([]Calculation, error) {
	return s.repo.GetAllCalculations()
}

func (s *calcService) GetCalculationsByID(id string) (Calculation, error) {
	return s.repo.GetCalculationsByID(id)
}

func (s *calcService) UpdateCalculation(id string, expression string) (Calculation, error) {
	calc, err := s.repo.GetCalculationsByID(id)

	if err != nil {
		return Calculation{}, err
	}

	result, err := s.calculateExpression(expression)
	if err != nil {
		return Calculation{}, err
	}
	calc.Expression = expression
	calc.Result = result

	if err := s.repo.UpdateCalculation(calc); err != nil {
		return Calculation{}, err
	}

	return calc, nil
}

func (s *calcService) DeleteCalculation(id string) error {
	return s.repo.DeleteCalculation(id)
}
