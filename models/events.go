package models

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis"
)

// Event is base event type
type Event struct {
	AggregateID string
	Data        interface{}
	EventType   string
	CreatedAt   time.Time
}

const (
	EventTypeKey = "event_type"
	DataKey      = "data"
)

// SerializeRedisStream serializes to datatype that redis stream supports
func (e *Event) SerializeRedisStream() *redis.XAddArgs {
	dataJSON, _ := json.Marshal(e.Data)

	return &redis.XAddArgs{
		Stream: e.AggregateID,
		Values: map[string]interface{}{
			"event_type": e.EventType,
			"data":       dataJSON,
			"created_at": e.CreatedAt,
		},
	}
}

func DeserializeRedisStream(message redis.XMessage) *Event {
	return &Event{
		Data:      message.Values[DataKey],
		EventType: message.Values[EventTypeKey].(string),
	}
}

func BuildSessionEvents(events []*Event, sessionID string) (*Session, error) {
	if len(events) == 0 {
		return nil, errors.New("received no events")
	}

	if events[0].EventType != NewSessionEventType {
		return nil, errors.New("initial event is not valid")
	}

	session := &Session{
		ID:       sessionID,
		Computer: NewBoard(true),
		Player:   NewBoard(false),
	}

	for _, event := range events {
		session.Apply(event)
	}

	return session, nil
}

func (s *Session) Apply(event *Event) error {
	switch event.EventType {
	case NewSessionEventType:
		if err := s.ApplyCreateSessionEvent(event); err != nil {
			return err
		}
	case ShootEventType:
		if err := s.ApplyShootEvent(event); err != nil {
			return err
		}
	case DestroyShipEventType:
		if err := s.ApplyDestroyShipEvent(event); err != nil {
			return err
		}
	}

	return nil
}

// ApplyCreateSessionEvent add player and computer boards and their
// IDs to newly created session struct
func (s *Session) ApplyCreateSessionEvent(event *Event) error {
	body := event.Data.(string)
	var payload NewSessionEventData
	if err := json.Unmarshal([]byte(body), &payload); err != nil {
		return err
	}

	s.Computer.Battleships = payload.Computer.Battleships
	s.Player.Battleships = payload.Player.Battleships
	return nil
}

// ApplyShootEvent handles shooting cells
func (s *Session) ApplyShootEvent(event *Event) error {
	body := event.Data.(string)
	var payload ShootEventData
	if err := json.Unmarshal([]byte(body), &payload); err != nil {
		return err
	}

	if payload.IsComputer {
		s.Player.RegisterShot(payload.Cell)
	} else {
		s.Computer.RegisterShot(payload.Cell)
	}
	return nil
}

// ApplyDestroyShipEvent applies ship as dead if its all
// cells are destroyed
func (s *Session) ApplyDestroyShipEvent(event *Event) error {
	body := event.Data.(string)
	var payload DestroyShipEventData
	if err := json.Unmarshal([]byte(body), &payload); err != nil {
		return err
	}

	if payload.IsComputer {
		s.Computer.MarkShipIfDead(payload.ShipID)
	} else {
		s.Player.MarkShipIfDead(payload.ShipID)
	}
	return nil
}

const (
	NewSessionEventType  = "new_session"
	ShootEventType       = "shoot"
	DestroyShipEventType = "destroy_ship"
)

type NewSessionEventData struct {
	Session
}

func CreateNewSessionEvent(session *Session) *Event {
	return &Event{
		AggregateID: session.ID,
		Data:        session,
		EventType:   NewSessionEventType,
		CreatedAt:   time.Now(),
	}
}

type ShootEventData struct {
	Cell
	IsComputer bool `json:"is_computer"`
}

func CreateShootEvent(sessionID string, cell *Cell, isComputer bool) *Event {
	return &Event{
		AggregateID: sessionID,
		Data: ShootEventData{
			Cell:       *cell,
			IsComputer: isComputer,
		},
		EventType: ShootEventType,
		CreatedAt: time.Now(),
	}
}

type DestroyShipEventData struct {
	ShipID     string `json:"ship_id"`
	IsComputer bool   `json:"is_computer"`
}

func CreateDestroyShipEvent(sessionID string, shipID string, isComputer bool) *Event {
	return &Event{
		AggregateID: sessionID,
		Data: DestroyShipEventData{
			ShipID:     shipID,
			IsComputer: isComputer,
		},
		EventType: DestroyShipEventType,
		CreatedAt: time.Now(),
	}
}
