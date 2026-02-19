package service

import (
	"context"

	"coviar_backend/internal/domain"
	"coviar_backend/internal/repository"
)

type AdminService struct {
	repo repository.AdminRepository
}

func NewAdminService(repo repository.AdminRepository) *AdminService {
	return &AdminService{repo: repo}
}

func (s *AdminService) GetStats(ctx context.Context) (*domain.AdminStatsResponse, error) {
	return s.repo.GetStats(ctx)
}

func (s *AdminService) GetAllEvaluaciones(ctx context.Context, estado string, idBodega int) ([]domain.EvaluacionListItem, error) {
	return s.repo.GetAllEvaluaciones(ctx, estado, idBodega)
}
