package service

import (
	"controller_server_golang/repository"
	"controller_server_golang/types/table"
)

type Service struct {
	repository *repository.Repository
}

func NewService(repository *repository.Repository) *Service {
	s := &Service{repository: repository}

	return s
}

func (s *Service) GetAvailableServerList() ([]*table.ServerInfo, error) {
	return s.repository.GetAvailableServerList()
}
