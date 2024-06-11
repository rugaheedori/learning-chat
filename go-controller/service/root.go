package service

import (
	"controller_server_golang/repository"
	"controller_server_golang/types/table"
	"encoding/json"
	"fmt"
	"log"

	. "github.com/confluentinc/confluent-kafka-go/kafka"
)

type Service struct {
	repository *repository.Repository

	AvgServerList map[string]bool
}

func NewService(repository *repository.Repository) *Service {
	s := &Service{repository: repository, AvgServerList: make(map[string]bool)}

	s.setServerInfo()

	if err := s.repository.Kafka.RegisterSubTopic("chat"); err != nil {
		panic(err)
	} else {
		go s.loopSubKafka()
	}

	return s
}

// 서브 스레드로 들어오는 이벤트를 감지해야 함
func (s *Service) loopSubKafka() {
	for {
		ev := s.repository.Kafka.Pool(100)

		switch event := ev.(type) {
		case *Message:

			type ServerInfoEvent struct {
				IP     string
				Status bool
			}

			var decoder ServerInfoEvent

			if err := json.Unmarshal(event.Value, &decoder); err != nil {
				log.Println("Failed To Decode Event", event.Value)
			} else {
				fmt.Println(decoder)
				s.AvgServerList[decoder.IP] = decoder.Status
			}

			fmt.Println(event)
		case *Error:
			log.Print("Failed To Pooling Event", event.Error())
		}
	}
}

func (s *Service) setServerInfo() {
	if serverList, err := s.GetAvailableServerList(); err != nil {
		panic(err)
	} else {
		for _, server := range serverList {
			s.AvgServerList[server.IP] = true
		}
	}
}

func (s *Service) GetAvgServrList() []string {
	var res []string

	for ip, available := range s.AvgServerList {
		if available {
			res = append(res, ip)
		}
	}

	return res
}

func (s *Service) GetAvailableServerList() ([]*table.ServerInfo, error) {
	return s.repository.GetAvailableServerList()
}
