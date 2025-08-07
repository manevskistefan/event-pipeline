package storage

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type EventType string

type Source string

type Data struct {
	Action   string                 `db:"action"`
	Value    float32                `db:"value"`
	Metadata map[string]interface{} `db:"metadata"`
}

type ProcessedEvent struct {
	ID        string    `db:"id"`
	Type      EventType `db:"type"`
	Source    Source    `db:"source"`
	Timestamp time.Time `db:"timestamp"`
	UserID    *string   `db:"user_id"`
	Data      Data      `db:"data"`
}

type eventRepository struct {
	db *sqlx.DB
}

type EventRepository interface {
	InsertEvent(id string, eventType EventType, source Source, timestamp time.Time, userId *string, data Data) (*ProcessedEvent, error)
}

func NewEventRepository(db *sqlx.DB) EventRepository {
	return &eventRepository{
		db: db,
	}
}

func (r *eventRepository) InsertEvent(id string, eventType EventType, source Source, timestamp time.Time, userId *string, data Data) (*ProcessedEvent, error) {

	event := &ProcessedEvent{
		ID:        id,
		Type:      eventType,
		Source:    source,
		Timestamp: timestamp,
		UserID:    userId,
		Data:      data,
	}

	query := `INSERT INTO events (id, type, source, timestamp, user_id, action, value, metadata) 
			  VALUES (:id, :type, :source, :timestamp, :user_id, :action, :value, :metadata)`

	_, err := r.db.NamedExec(query, event)
	if err != nil {
		return nil, err
	}

	return event, nil
}
