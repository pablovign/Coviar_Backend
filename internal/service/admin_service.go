package service

import (
	"context"

	"coviar_backend/internal/domain"
)

type AdminRepository interface {
	GetStats(ctx context.Context) (*domain.AdminStatsResponse, error)
	GetAllEvaluaciones(ctx context.Context, estado string, idBodega int) ([]domain.EvaluacionListItem, error)
}

type AdminService struct {
	repo AdminRepository
}

func NewAdminService(repo AdminRepository) *AdminService {
	return &AdminService{repo: repo}
}

func (s *AdminService) GetStats(ctx context.Context) (*domain.AdminStatsResponse, error) {
	return s.repo.GetStats(ctx)
}

func (s *AdminService) GetAllEvaluaciones(ctx context.Context, estado string, idBodega int) ([]domain.EvaluacionListItem, error) {
	return s.repo.GetAllEvaluaciones(ctx, estado, idBodega)
}
