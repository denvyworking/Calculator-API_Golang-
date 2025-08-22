package calculationservice

// только работа с базой данных

import "gorm.io/gorm"

type CalculationRepository interface {
	CreateCalculation(calc Calculation) error
	GetAllCalculations() ([]Calculation, error)
	GetCalculationsByID(id string) (Calculation, error)
	UpdateCalculation(calc Calculation) error
	DeleteCalculation(id string) error
}

type calcRepository struct {
	db *gorm.DB
}

func NewCalculationRepository(db *gorm.DB) CalculationRepository {
	return &calcRepository{db: db}
}

func (cr *calcRepository) CreateCalculation(calc Calculation) error {
	return cr.db.Create(&calc).Error
}

func (cr *calcRepository) GetAllCalculations() ([]Calculation, error) {
	var calculations []Calculation
	err := cr.db.Find(&calculations).Error
	return calculations, err
}

func (cr *calcRepository) GetCalculationsByID(id string) (Calculation, error) {
	var calc Calculation
	err := cr.db.First(&calc, "id = ?", id).Error
	return calc, err
}

func (cr *calcRepository) UpdateCalculation(calc Calculation) error {
	return cr.db.Save(&calc).Error
}

func (cr *calcRepository) DeleteCalculation(id string) error {
	return cr.db.Delete(&Calculation{}, "id = ?", id).Error
}
