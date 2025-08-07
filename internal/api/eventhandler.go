package api

import (
	"encoding/json"
	api "event-processing-pipeline/internal/api/dtos"
	"event-processing-pipeline/internal/pipeline"
	"io"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type EventPipeline struct {
	ingestionChan chan api.EventDTO
	workerPool    []*Worker
	storage       Storage
	// metrics       *Metrics
	ctx *gin.Context
}

type Worker struct {
	Id       int
	jobChan  chan api.EventDTO
	pipeline *EventPipeline
}

type eventController struct {
	eventService pipeline.EventService
}

type EventController interface {
	HandleSingleEvent(ctx *gin.Context)
	HandleEventsBatch(ctx *gin.Context)
	GetMetrics(ctx *gin.Context)
}

func NewEventController(db *sqlx.DB) EventController {
	eventService := pipeline.NewEventService(db)

	return &eventController{
		eventService: eventService,
	}
}

func (c *eventController) HandleSingleEvent(ctx *gin.Context) {
	body, _ := io.ReadAll(ctx.Request.Body)
	var event api.EventDTO
	if err := json.Unmarshal(body, &event); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	c.eventService.Validate(*ctx, event)
}

func (c *eventController) HandleEventsBatch(ctx *gin.Context) {
	body, _ := io.ReadAll(ctx.Request.Body)
	var events []api.EventDTO
	if err := json.Unmarshal(body, &events); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{"status": "batch processing started"})

	workers := make([]Worker, len(events))
	wg := &sync.WaitGroup{}

	for i, event := range events {
		worker := &Worker{
			Id:      i,
			jobChan: make(chan api.EventDTO),
			pipeline: &EventPipeline{
				ingestionChan: make(chan api.EventDTO),
				ctx:           ctx,
			}}

		worker.Start(ctx)
	}

}

func (w *Worker) Start(ctx *gin.Context) {
	go func() {
		for {
			select {
			case job := <-w.jobChan:
				w.processJob(ctx, job)
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (c *eventController) GetMetrics(ctx *gin.Context) {
	// Assuming metrics are not implemented yet
	ctx.JSON(http.StatusOK, gin.H{"status": "metrics not implemented"})
}
