package ports

import (
	"context"

	"ApiSmart/internal/core/domain"
)

type AuthService interface {
	Register(ctx context.Context, req domain.RegisterRequest) (*domain.AuthResponse, error)
	Login(ctx context.Context, req domain.LoginRequest) (*domain.AuthResponse, error)
	ValidateToken(token string) (uint, error)
}

type SensorService interface {
	SaveSensorData(ctx context.Context, data *domain.SensorData) error
	GetAllSensorData(ctx context.Context) ([]domain.SensorData, error)
	GetLatestSensorData(ctx context.Context) (*domain.SensorData, error)
	GetAlerts(ctx context.Context, isRead *bool) ([]domain.Alert, error)
	MarkAlertAsRead(ctx context.Context, alertID uint) error
}

type AlertService interface {
	CheckAndCreateAlerts(data *domain.SensorData) []domain.Alert
}
