package service

import (
	"controller_server_golang/repository"
	"controller_server_golang/types/table"
)

type Service struct {
	repository *repository.Repository

	avgServerList map[string]bool
}

func NewService(repository *repository.Repository) *Service {
	s := &Service{repository: repository, avgServerList: make(map[string]bool)}

	s.setServerInfo()

	return s
}

func (s *Service) setServerInfo() {
	if serverList, err := s.GetAvailableServerList(); err != nil {
		panic(err)
	} else {
		for _, server := range serverList {
			s.avgServerList[server.IP] = true
		}
	}
}

func (s *Service) GetAvailableServerList() ([]*table.ServerInfo, error) {
	return s.repository.GetAvailableServerList()
}
