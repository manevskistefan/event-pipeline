package pipeline

import (
	"errors"
	api "event-processing-pipeline/internal/api/dtos"
	"event-processing-pipeline/internal/storage"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type eventService struct {
	eventRepository storage.EventRepository
}

type Validator interface {
	Validate(ctx gin.Context, event api.EventDTO) error
}

type Processor interface {
	Process(ctx gin.Context, event api.EventDTO) (*storage.ProcessedEvent, error)
}

type Storage interface {
	Store(ctx gin.Context, events []storage.ProcessedEvent) error
}

type EventService interface {
	Validator
	Processor
	Storage
}

func NewEventService(db *sqlx.DB) EventService {
	eventRepository := storage.NewEventRepository(db)

	return &eventService{
		eventRepository: eventRepository,
	}
}

func (s *eventService) Validate(ctx gin.Context, event api.EventDTO) error {
	if event.Type == "" {
		return errors.New("event type is required")
	}

	if event.Source == "" {
		return errors.New("event source is required")
	}

	return nil
}

func (s *eventService) Process(ctx gin.Context, event api.EventDTO) (*storage.ProcessedEvent, error) {
	time.Sleep(10)

	return &storage.ProcessedEvent{
		ID:        *event.ID,
		Type:      storage.EventType(event.Type),
		Source:    storage.Source(event.Source),
		Timestamp: event.Timestamp,
		UserID:    event.UserID,
		Data: storage.Data{
			Action:   event.Data.Action,
			Value:    event.Data.Value,
			Metadata: event.Data.Metadata,
		},
	}, nil
}

func (s *eventService) Store(ctx gin.Context, events []storage.ProcessedEvent) error {
	for _, event := range events {
		savedEvent, err := s.eventRepository.InsertEvent(
			event.ID,
			storage.EventType(event.Type),
			storage.Source(event.Source),
			event.Timestamp,
			event.UserID,
			storage.Data{
				Action:   event.Data.Action,
				Value:    event.Data.Value,
				Metadata: event.Data.Metadata,
			})

		if err != nil {
			log.Println("Event saved:", savedEvent)
			continue
		}

		return err
	}

	return nil
}

func (w *Worker) processJob(ctx *gin.Context, job api.EventDTO) {

}
