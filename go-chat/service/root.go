package service

import (
	"chat_server_golang/repository"
	"chat_server_golang/types/schema"
	"encoding/json"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Service struct {
	repository *repository.Repository
}

func NewService(repository *repository.Repository) *Service {
	s := &Service{repository: repository}

	return s
}

func (s *Service) PublishServerStatusEvent(ip string, status bool) {
	// Kafka에 이벤트 전송
	type ServerInfoEvent struct {
		IP     string
		Status bool
	}

	e := &ServerInfoEvent{IP: ip, Status: status}
	ch := make(chan kafka.Event)

	// 값을 전송할 때는 배열 바이트 값으로 전송해야 함
	if v, err := json.Marshal(e); err != nil {
		log.Println("Failed To Marshal")
	} else if result, err := s.PublishEvent("chat", v, ch); err != nil {
		// TODO Send Event To Kafka
		log.Println("Failed To Send Event to Kafka", "err", err)
	} else {
		log.Println("Success To Send Event", result)
	}
}

func (s *Service) PublishEvent(topic string, value []byte, ch chan kafka.Event) (kafka.Event, error) {
	return s.repository.Kafka.PublishEvent(topic, value, ch)
}

func (s *Service) ServerSet(ip string, available bool) error {
	if err := s.repository.ServerSet(ip, available); err != nil {
		log.Println("Failed To ServerSet", "ip", ip, "available", available)
		return err
	} else {
		return nil
	}
}

func (s *Service) InsertChatting(user, message, roomName string) {
	if err := s.repository.InsertChatting(user, message, roomName); err != nil {
		log.Println("Failed To Chat", "err", err)
	}
}

func (s *Service) EnterRoom(roomName string) ([]*schema.Chat, error) {
	if res, err := s.repository.GetChatList(roomName); err != nil {
		log.Println("Failed To Get All Chat LIist", "err", err.Error())
		return nil, err
	} else {
		return res, nil
	}
}

func (s *Service) RoomList() ([]*schema.Room, error) {
	if res, err := s.repository.RoomList(); err != nil {
		log.Println("Failed To Get All Room LIist", "err", err.Error())
		return nil, err
	} else {
		return res, nil
	}
}

func (s *Service) MakeRoom(name string) error {
	if err := s.repository.MakeRoom(name); err != nil {
		log.Println("Failed To Make New Room", "err", err.Error())
		return err
	} else {
		return nil
	}
}

func (s *Service) Room(name string) (*schema.Room, error) {
	if res, err := s.repository.Room(name); err != nil {
		log.Println("Failed To Get Room", "err", err.Error())
		return nil, err
	} else {
		return res, nil
	}
}
