package repository

import (
	"chat_server_golang/config"
	"chat_server_golang/types/schema"
	"database/sql"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Repository struct {
	cfg *config.Config

	db *sql.DB
}

const (
	room       = "chatting.room"
	chat       = "chatting.chat"
	serverInfo = "chatting.serverInfo"
)

func NewRepository(cfg *config.Config) (*Repository, error) {
	r := &Repository{cfg: cfg}
	var err error

	if r.db, err = sql.Open(cfg.DB.Database, cfg.DB.URL); err != nil {
		return nil, err
	} else {
		return r, nil
	}
}

func (s *Repository) RoomList() ([]*schema.Room, error) {
	// TODO 페이징 추가하기
	qs := query([]string{"SELECT * FROM", room})

	// cursor: 쿼리에 대해서 데이터를 가지고 있으며, 메모리에 할당됨
	if cursor, err := s.db.Query(qs); err != nil {
		return nil, err
	} else {
		defer cursor.Close()

		var result []*schema.Room

		for cursor.Next() {
			d := new(schema.Room)

			if err = cursor.Scan(
				&d.ID,
				&d.Name,
				&d.CreateAt,
				&d, d.UpdateAt,
			); err != nil {
				return nil, err
			} else {
				result = append(result, d)
			}
		}

		if len(result) == 0 {
			return []*schema.Room{}, nil
		}
		return result, nil
	}
}

func (s *Repository) MakeRoom(name string) error {
	_, err := s.db.Exec("INSERT INTO chatting.room(name) VALUES(?)", name)
	return err
}

func (s *Repository) Room(name string) (*schema.Room, error) {
	d := new(schema.Room)
	qs := query([]string{"SELECT * FROM", room, "WHERE name = ?"})

	err := s.db.QueryRow(qs, name).Scan(
		&d.ID,
		&d.Name,
		&d.CreateAt,
		&d, d.UpdateAt,
	)

	return d, err
}

func query(qs []string) string {
	return strings.Join(qs, " ") + ";"
}
