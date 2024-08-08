package postgres

import (
	"courier/config"
	"courier/storage"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage struct {
	Db            *sql.DB
	NotificationS storage.NotificationI
	TaskS         storage.TaskI
}

func NewPostgresStorage(config config.Config) (*Storage, error) {
	conn := fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%d sslmode=disable",
		config.DB_HOST, config.DB_USER, config.DB_NAME, config.DB_PASSWORD, config.DB_PORT)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	notification := NewNotificationRepo(db)
	task := NewTaskRepo(db)
	// comment := NewCommentsRepo(db)
	// shared := NewSharedMemoriesRepo(db)

	return &Storage{
		Db:            db,
		NotificationS: notification,
		TaskS:         task,
		// SharedMemoryS: shared,
		// CommentS:      comment,
	}, nil
}
