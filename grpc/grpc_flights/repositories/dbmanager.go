package repositories

import (
	"assignment/grpc/grpc_flights/models"
	"context"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type FlightRepository interface {
	CreatFlight(ctx context.Context, model *models.Flight) (*models.Flight, error)
	UpdateFlight(ctx context.Context, model *models.Flight) (*models.Flight, error)
	SearchFlight(ctx context.Context, id uuid.UUID) (*models.Flight, error)
}

type dbManager struct {
	*gorm.DB
}

func NewDBManager() (FlightRepository, error) {
	db, err := gorm.Open(postgres.Open("host=localhost user=postgres password=Hieu200599 dbname=flights port=5432 sslmode=disable TimeZone=Asia/Ho_Chi_Minh"))
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.Flight{})

	return &dbManager{db}, nil
}

//tạo flight mới
func (m *dbManager) CreatFlight(ctx context.Context, model *models.Flight) (*models.Flight, error) {
	if err := m.Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

//tìm flight bằng id
func (m *dbManager) SearchFlight(ctx context.Context, id uuid.UUID) (*models.Flight, error) {
	var result models.Flight
	err := m.Where(&models.Flight{ID: id}).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

//update flight
func (m *dbManager) UpdateFlight(ctx context.Context, model *models.Flight) (*models.Flight, error) {
	err := m.Where(&models.Flight{ID: model.ID}).Updates(model).Error
	if err != nil {
		return nil, err
	}
	return model, nil
}
