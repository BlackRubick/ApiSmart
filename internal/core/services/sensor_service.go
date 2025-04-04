package services

import (
	"context"

	"ApiSmart/internal/core/domain"
	"ApiSmart/internal/core/ports"
)

type sensorService struct {
	sensorRepo   ports.SensorRepository
	alertService ports.AlertService
}

func NewSensorService(sensorRepo ports.SensorRepository, alertService ports.AlertService) ports.SensorService {
	return &sensorService{
		sensorRepo:   sensorRepo,
		alertService: alertService,
	}
}

func (s *sensorService) SaveSensorData(ctx context.Context, data *domain.SensorData) error {
	// Guardar los datos del sensor
	if err := s.sensorRepo.SaveSensorData(ctx, data); err != nil {
		return err
	}

	// Verificar si se deben generar alertas
	alerts := s.alertService.CheckAndCreateAlerts(data)

	// Guardar las alertas generadas
	for _, alert := range alerts {
		if err := s.sensorRepo.SaveAlert(ctx, &alert); err != nil {
			return err
		}
	}

	return nil
}

func (s *sensorService) GetAllSensorData(ctx context.Context) ([]domain.SensorData, error) {
	return s.sensorRepo.GetAllSensorData(ctx)
}

func (s *sensorService) GetLatestSensorData(ctx context.Context) (*domain.SensorData, error) {
	return s.sensorRepo.GetLatestSensorData(ctx)
}

func (s *sensorService) GetAlerts(ctx context.Context, isRead *bool) ([]domain.Alert, error) {
	return s.sensorRepo.GetAlerts(ctx, isRead)
}

func (s *sensorService) MarkAlertAsRead(ctx context.Context, alertID uint) error {
	return s.sensorRepo.MarkAlertAsRead(ctx, alertID)
}
